package storage

import (
	"context"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ObjectStorage interface {
	Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error
	GenerateReadURL(ctx context.Context, objectName string, ttl time.Duration) (string, error)
}

type MinioService struct {
	client *minio.Client
	bucket string
}

func NewMinioServiceFromEnv() (*MinioService, error) {
	endpoint := strings.TrimSpace(os.Getenv("MINIO_ENDPOINT"))
	accessKey := strings.TrimSpace(os.Getenv("MINIO_ACCESS_KEY"))
	secretKey := strings.TrimSpace(os.Getenv("MINIO_SECRET_KEY"))
	bucket := strings.TrimSpace(os.Getenv("MINIO_BUCKET"))
	useSSLRaw := strings.TrimSpace(os.Getenv("MINIO_USE_SSL"))

	if endpoint == "" || accessKey == "" || secretKey == "" || bucket == "" {
		return nil, errors.New("configuração do MinIO incompleta")
	}

	useSSL := false
	if useSSLRaw != "" {
		parsed, err := strconv.ParseBool(useSSLRaw)
		if err != nil {
			return nil, errors.New("MINIO_USE_SSL inválido")
		}
		useSSL = parsed
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &MinioService{client: client, bucket: bucket}, nil
}

func (s *MinioService) ensureBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{})
}

func (s *MinioService) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error {
	if objectName == "" {
		return errors.New("nome do objeto ausente")
	}
	if reader == nil {
		return errors.New("arquivo ausente")
	}
	if err := s.ensureBucket(ctx); err != nil {
		return err
	}

	opts := minio.PutObjectOptions{}
	if contentType != "" {
		opts.ContentType = contentType
	}

	_, err := s.client.PutObject(ctx, s.bucket, objectName, reader, size, opts)
	return err
}

func (s *MinioService) GenerateReadURL(ctx context.Context, objectName string, ttl time.Duration) (string, error) {
	if objectName == "" {
		return "", errors.New("nome do objeto ausente")
	}
	url, err := s.client.PresignedGetObject(ctx, s.bucket, objectName, ttl, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
