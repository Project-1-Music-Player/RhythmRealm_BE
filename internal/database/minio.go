package database

import (
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOService interface {
	Health() map[string]string
	StreamMusic(ctx echo.Context, bucketName, objectName string) (*minio.Object, error)
}

type minIOService struct {
	client *minio.Client
}

func NewMinIO() MinIOService {
	endpoint := "localhost:9000"
	accessKeyID := "minio"
	secretAccessKey := "minio123"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatalln(err)
	}

	s := &minIOService{
		client: minioClient,
	}
	return s
}

func (s *minIOService) Health() map[string]string {
	exists, err := s.client.BucketExists(context.Background(), "music")
	if err != nil || !exists {
		log.Fatalf("storage down: %v", err)
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *minIOService) StreamMusic(ctx echo.Context, bucketName, objectName string) (*minio.Object, error) {
	obj, err := s.client.GetObject(ctx.Request().Context(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return obj, nil
}
