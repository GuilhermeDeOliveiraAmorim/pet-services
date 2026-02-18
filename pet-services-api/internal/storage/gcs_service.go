package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
)

type GCSService struct {
	client     *storage.Client
	bucket     *storage.BucketHandle
	bucketName string
}

// NewGCSServiceFromEnv initializes a GCS service from environment variables
func NewGCSServiceFromEnv() (*GCSService, error) {
	ctx := context.Background()

	// Try new environment variables first (GCS_*)
	bucketName := strings.TrimSpace(os.Getenv("GCS_BUCKET_NAME"))
	credentialsPath := strings.TrimSpace(os.Getenv("GCS_CREDENTIALS_PATH"))

	// Fallback to existing variables
	if bucketName == "" {
		bucketName = strings.TrimSpace(os.Getenv("IMAGE_BUCKET_NAME"))
	}
	if credentialsPath == "" {
		credentialsPath = strings.TrimSpace(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	}

	if bucketName == "" {
		return nil, errors.New("GCS_BUCKET_NAME ou IMAGE_BUCKET_NAME são obrigatórios")
	}

	// Set credentials if provided
	if credentialsPath != "" {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credentialsPath)
	}

	// Create client
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar cliente GCS: %w", err)
	}

	bucket := client.Bucket(bucketName)

	// Verify bucket exists
	_, err = bucket.Attrs(ctx)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("falha ao acessar bucket GCS: %w", err)
	}

	return &GCSService{
		client:     client,
		bucket:     bucket,
		bucketName: bucketName,
	}, nil
}

// Upload uploads an object to GCS
func (s *GCSService) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error {
	if objectName == "" {
		return errors.New("nome do objeto ausente")
	}
	if reader == nil {
		return errors.New("arquivo ausente")
	}

	wctx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()

	obj := s.bucket.Object(objectName)
	w := obj.NewWriter(wctx)

	if contentType != "" {
		w.ContentType = contentType
	}

	// Copy the data
	if _, err := io.Copy(w, reader); err != nil {
		w.Close()
		return fmt.Errorf("falha ao fazer upload do objeto: %w", err)
	}

	// Close writer (this completes the upload)
	if err := w.Close(); err != nil {
		return fmt.Errorf("falha ao finalizar upload: %w", err)
	}

	return nil
}

// GenerateReadURL generates a signed URL for reading the object
func (s *GCSService) GenerateReadURL(ctx context.Context, objectName string, ttl time.Duration) (string, error) {
	if objectName == "" {
		return "", errors.New("nome do objeto ausente")
	}

	// Generate signed URL
	url, err := s.bucket.SignedURL(objectName, &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(ttl),
	})

	if err != nil {
		return "", fmt.Errorf("falha ao gerar URL assinada: %w", err)
	}

	return url, nil
}

// Delete deletes an object from GCS
func (s *GCSService) Delete(ctx context.Context, objectName string) error {
	if objectName == "" {
		return errors.New("nome do objeto ausente")
	}

	wctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	obj := s.bucket.Object(objectName)

	if err := obj.Delete(wctx); err != nil && err != storage.ErrObjectNotExist {
		return fmt.Errorf("falha ao deletar objeto: %w", err)
	}

	return nil
}

// Close closes the GCS client
func (s *GCSService) Close() error {
	return s.client.Close()
}
