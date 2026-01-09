package minio

import (
	"bytes"
	"context"
	"fmt"

	"pet-services-api/internal/application/logging"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioService struct {
	Client *minio.Client
	Bucket string
	Logger logging.LoggerService
}

func NewMinioService(endpoint, accessKey, secretKey, bucket string, useSSL bool, logger logging.LoggerService) (*MinioService, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		if logger != nil {
			logger.Log(logging.Logger{
				TypeLog: logging.LoggerTypes.ERROR,
				Layer:   logging.LoggerLayers.SERVICES,
				From:    "NewMinioService",
				Message: fmt.Sprintf("Erro ao criar cliente Minio: %v", err),
				Error:   err,
			})
		}

		return nil, err
	}

	return &MinioService{Client: client, Bucket: bucket, Logger: logger}, nil
}

func (s *MinioService) Upload(ctx context.Context, objectName string, fileData []byte, contentType string) (string, error) {
	reader := bytes.NewReader(fileData)

	_, err := s.Client.PutObject(ctx, s.Bucket, objectName, reader, int64(len(fileData)), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		if s.Logger != nil {
			s.Logger.Log(logging.Logger{
				Context: ctx,
				TypeLog: logging.LoggerTypes.ERROR,
				Layer:   logging.LoggerLayers.SERVICES,
				From:    "MinioService.Upload",
				Message: fmt.Sprintf("Erro ao fazer upload para Minio: %v", err),
				Error:   err,
			})
		}

		return "", err
	}

	url := fmt.Sprintf("https://%s/%s/%s", s.Client.EndpointURL().Host, s.Bucket, objectName)
	if s.Logger != nil {
		s.Logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.INFO,
			Layer:   logging.LoggerLayers.SERVICES,
			From:    "MinioService.Upload",
			Message: fmt.Sprintf("Upload realizado com sucesso: %s", url),
		})
	}

	return url, nil
}
