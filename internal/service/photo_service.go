package service

import (
	"context"
	"fmt"
	"io"

	//"mime"
	//"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type PhotoService struct {
	minioClient *minio.Client
	bucket      string
}

func NewPhotoService(client *minio.Client, bucket string) *PhotoService {
	return &PhotoService{
		minioClient: client,
		bucket:      bucket,
	}
}

func (s *PhotoService) UploadPhoto(ctx context.Context, file io.Reader, size int64) (string, error) {
	// Генерируем уникальное имя файла с правильным расширением
	objectName := uuid.New().String() + ".jpg"

	_, err := s.minioClient.PutObject(
		ctx,
		s.bucket,
		objectName,
		file,
		size,
		minio.PutObjectOptions{
			ContentType: "image/jpeg",
		},
	)

	if err != nil {
		return "", fmt.Errorf("upload failed: %w", err)
	}

	return objectName, nil
}

func (s *PhotoService) GetPhotoURL(objectName string, expiry time.Duration) (string, error) {
	url, err := s.minioClient.PresignedGetObject(
		context.Background(),
		s.bucket,
		objectName,
		expiry,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}
