package minio

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/kerilOvs/profile_sevice/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

/*
accessKey: "ErKo2T8FdXxvR4phdFop"
  secretKey: "uV8044JzvasUXq6u62RD6DD3JVSsJq3x4w225AAl"
*/

type Client struct {
	Client       *minio.Client
	Bucket       string
	PublicHost   string // Добавляем поле для публичного хоста
	PublicPrefix string // Публичный префикс
}

func New(ctx context.Context, cfg config.MinioConfig) (*Client, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio init error: %w", err)
	}

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("bucket check error: %w", err)
	}

	if !exists {
		if err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("bucket creation error: %w", err)
		}
	}

	return &Client{
		Client:       client,
		Bucket:       cfg.Bucket,
		PublicHost:   "localhost",
		PublicPrefix: "pub/",
	}, nil
}

func (c *Client) GenerateUploadURL1(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	presignedURL, err := c.Client.PresignedPutObject(ctx, c.Bucket, objectName, expiry)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

// Добавляем методы для работы с объектами -- не используется блять....
func (c *Client) PutObject1(ctx context.Context, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	return c.Client.PutObject(ctx, c.Bucket, objectName, reader, objectSize, opts)
}

func (c *Client) PresignedGetObject4(ctx context.Context, objectName string, expiry time.Duration, reqParams url.Values) (string, error) {
	// return c.client.PresignedGetObject(ctx, c.bucket, objectName, expiry, reqParams)
	presignedURL, err := c.Client.PresignedGetObject(ctx, c.Bucket, objectName, expiry, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}
