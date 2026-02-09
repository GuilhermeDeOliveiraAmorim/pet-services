package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ObjectStorage interface {
	UploadImage(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error)
}

type MinioService struct {
	client *minio.Client
	bucket string
	publicRead bool
}

func NewMinioServiceFromEnv() (*MinioService, error) {
	endpoint := strings.TrimSpace(os.Getenv("MINIO_ENDPOINT"))
	accessKey := strings.TrimSpace(os.Getenv("MINIO_ACCESS_KEY"))
	secretKey := strings.TrimSpace(os.Getenv("MINIO_SECRET_KEY"))
	bucket := strings.TrimSpace(os.Getenv("MINIO_BUCKET"))
	useSSLRaw := strings.TrimSpace(os.Getenv("MINIO_USE_SSL"))
	publicReadRaw := strings.TrimSpace(os.Getenv("MINIO_PUBLIC_READ"))

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
	publicRead := false
	if publicReadRaw != "" {
		parsed, err := strconv.ParseBool(publicReadRaw)
		if err != nil {
			return nil, errors.New("MINIO_PUBLIC_READ inválido")
		}
		publicRead = parsed
	}

	return &MinioService{client: client, bucket: bucket, publicRead: publicRead}, nil
}

func (s *MinioService) ensureBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return err
	}
	if exists {
		if s.publicRead {
			return s.ensurePublicReadPolicy(ctx)
		}
		return nil
	}
	if err := s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{}); err != nil {
		return err
	}
	if s.publicRead {
		return s.ensurePublicReadPolicy(ctx)
	}
	return nil
}

func (s *MinioService) ensurePublicReadPolicy(ctx context.Context) error {
	policy := fmt.Sprintf(`{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {"AWS": ["*"]},
      "Action": ["s3:GetObject"],
      "Resource": ["arn:aws:s3:::%s/*"]
    }
  ]
}`, s.bucket)

	return s.client.SetBucketPolicy(ctx, s.bucket, policy)
}

func (s *MinioService) UploadImage(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	if objectName == "" {
		return "", errors.New("nome do objeto ausente")
	}
	if reader == nil {
		return "", errors.New("arquivo ausente")
	}
	if err := s.ensureBucket(ctx); err != nil {
		return "", err
	}

	opts := minio.PutObjectOptions{}
	if contentType != "" {
		opts.ContentType = contentType
	}

	_, err := s.client.PutObject(ctx, s.bucket, objectName, reader, size, opts)
	if err != nil {
		return "", err
	}

	endpointURL := strings.TrimRight(s.client.EndpointURL().String(), "/")
	return fmt.Sprintf("%s/%s/%s", endpointURL, s.bucket, objectName), nil
}
