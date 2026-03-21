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

func NewGCSServiceFromEnv() (*GCSService, error) {
	ctx := context.Background()

	bucketName := strings.TrimSpace(os.Getenv("GCS_BUCKET_NAME"))
	credentialsPath := strings.TrimSpace(os.Getenv("GCS_CREDENTIALS_PATH"))

	if bucketName == "" {
		bucketName = strings.TrimSpace(os.Getenv("IMAGE_BUCKET_NAME"))
	}
	if credentialsPath == "" {
		credentialsPath = strings.TrimSpace(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	}

	if bucketName == "" {
		return nil, errors.New("GCS_BUCKET_NAME ou IMAGE_BUCKET_NAME são obrigatórios")
	}

	if credentialsPath != "" {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credentialsPath)
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar cliente GCS: %w", err)
	}

	bucket := client.Bucket(bucketName)

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

	if _, err := io.Copy(w, reader); err != nil {
		w.Close()
		return fmt.Errorf("falha ao fazer upload do objeto: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("falha ao finalizar upload: %w", err)
	}

	return nil
}

func (s *GCSService) GenerateReadURL(ctx context.Context, objectName string, ttl time.Duration) (string, error) {
	if objectName == "" {
		return "", errors.New("nome do objeto ausente")
	}

	url, err := s.bucket.SignedURL(objectName, &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(ttl),
	})

	if err != nil {
		return "", fmt.Errorf("falha ao gerar URL assinada: %w", err)
	}

	return url, nil
}

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

func (s *GCSService) Close() error {
	return s.client.Close()
}
