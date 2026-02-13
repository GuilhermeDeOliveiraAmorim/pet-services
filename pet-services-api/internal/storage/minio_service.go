package storage

import (
	"context"
	"errors"
	"io"
	"net/url"
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
	Delete(ctx context.Context, objectName string) error
}

type MinioService struct {
	client         *minio.Client
	bucket         string
	publicEndpoint string
	creds          *credentials.Credentials
	region         string
}

func NewMinioServiceFromEnv() (*MinioService, error) {
	endpoint := strings.TrimSpace(os.Getenv("MINIO_ENDPOINT"))
	accessKey := strings.TrimSpace(os.Getenv("MINIO_ACCESS_KEY"))
	secretKey := strings.TrimSpace(os.Getenv("MINIO_SECRET_KEY"))
	bucket := strings.TrimSpace(os.Getenv("MINIO_BUCKET"))
	useSSLRaw := strings.TrimSpace(os.Getenv("MINIO_USE_SSL"))
	publicEndpoint := strings.TrimSpace(os.Getenv("MINIO_PUBLIC_ENDPOINT"))
	region := strings.TrimSpace(os.Getenv("MINIO_REGION"))
	if region == "" {
		region = "us-east-1"
	}

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

	creds := credentials.NewStaticV4(accessKey, secretKey, "")
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  creds,
		Secure: useSSL,
		Region: region,
	})
	if err != nil {
		return nil, err
	}
	return &MinioService{client: client, bucket: bucket, publicEndpoint: publicEndpoint, creds: creds, region: region}, nil
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
	client := s.client
	if s.publicEndpoint != "" {
		publicURL, err := url.Parse(s.publicEndpoint)
		if err != nil {
			return "", err
		}
		endpoint := publicURL.Host
		if endpoint == "" {
			endpoint = strings.TrimSpace(s.publicEndpoint)
		}
		useSSL := strings.EqualFold(publicURL.Scheme, "https")
		presignClient, err := minio.New(endpoint, &minio.Options{
			Creds:        s.creds,
			Secure:       useSSL,
			BucketLookup: minio.BucketLookupPath,
			Region:       s.region,
		})
		if err != nil {
			return "", err
		}
		client = presignClient
	}

	presignedURL, err := client.PresignedGetObject(ctx, s.bucket, objectName, ttl, nil)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

func (s *MinioService) Delete(ctx context.Context, objectName string) error {
	if objectName == "" {
		return errors.New("nome do objeto ausente")
	}
	return s.client.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{})
}
