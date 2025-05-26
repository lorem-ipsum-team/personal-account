package service

import (
	"context"
	"fmt"
	"io"
	"strings"

	//"mime"
	//"path/filepath"
	"time"

	"github.com/kerilOvs/profile_sevice/internal/config"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type PhotoService struct {
	minioClient *minio.Client
	bucket      string
	pubPrefix   string
	host        string
}

func NewPhotoService(client *minio.Client, cfg config.MinioConfig) *PhotoService {
	return &PhotoService{
		minioClient: client,
		bucket:      cfg.Bucket,
		pubPrefix:   cfg.PubPrefix,
		host:        cfg.Host,
	}
}

func (s *PhotoService) UploadPhoto(ctx context.Context, file io.Reader, size int64) (string, error) {
	// Генерируем уникальное имя файла с правильным расширением
	objectName := s.pubPrefix + "/" + uuid.New().String() + ".jpg"

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

	PublicPrefix := s.pubPrefix
	PublicHost := s.host

	if strings.HasPrefix(objectName, PublicPrefix) {
		publicURL := fmt.Sprintf("%s/%s/%s", PublicHost, s.bucket, objectName)

		return publicURL, nil
	}

	return "hui", nil
}
