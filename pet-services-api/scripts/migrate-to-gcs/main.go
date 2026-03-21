package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	minioEndpoint := flag.String("minio-endpoint", "localhost:9000", "MinIO endpoint")
	minioBucket := flag.String("minio-bucket", "pet-services", "MinIO bucket name")
	minioAccessKey := flag.String("minio-access-key", "petimages", "MinIO access key")
	minioSecretKey := flag.String("minio-secret-key", "petimages123", "MinIO secret key")

	gcsBucket := flag.String("gcs-bucket", "pet-services-bucket", "GCS bucket name")
	gcsProject := flag.String("gcs-project", "", "GCS project ID")

	flag.Parse()

	if *gcsProject == "" {
		log.Fatal("--gcs-project é obrigatório")
	}

	ctx := context.Background()

	minioClient, err := minio.New(*minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(*minioAccessKey, *minioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("Erro ao conectar MinIO: %v", err)
	}

	gcsClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Erro ao criar cliente GCS: %v", err)
	}
	defer gcsClient.Close()

	gcsBucketHandle := gcsClient.Bucket(*gcsBucket)

	objectsCh := minioClient.ListObjects(ctx, *minioBucket, minio.ListObjectsOptions{
		Recursive: true,
	})

	migratedCount := 0
	errorCount := 0

	for object := range objectsCh {
		if object.Err != nil {
			log.Printf("Erro ao listar objetos: %v", object.Err)
			errorCount++
			continue
		}

		objectName := object.Key
		fmt.Printf("Migrando: %s... ", objectName)

		reader, err := minioClient.GetObject(ctx, *minioBucket, objectName, minio.GetObjectOptions{})
		if err != nil {
			fmt.Printf("ERRO (download): %v\n", err)
			errorCount++
			continue
		}
		defer reader.Close()

		wctx, cancel := context.WithCancel(ctx)
		w := gcsBucketHandle.Object(objectName).NewWriter(wctx)

		if object.ContentType != "" {
			w.ContentType = object.ContentType
		}

		_, err = io.Copy(w, reader)
		if err != nil {
			fmt.Printf("ERRO (upload): %v\n", err)
			w.Close()
			cancel()
			errorCount++
			continue
		}

		err = w.Close()
		if err != nil {
			fmt.Printf("ERRO (close): %v\n", err)
			cancel()
			errorCount++
			continue
		}

		cancel()
		fmt.Printf("OK\n")
		migratedCount++
	}

	fmt.Printf("\n=== Resumo da Migração ===\n")
	fmt.Printf("Migrados com sucesso: %d\n", migratedCount)
	fmt.Printf("Erros: %d\n", errorCount)

	if errorCount > 0 {
		os.Exit(1)
	}
}
