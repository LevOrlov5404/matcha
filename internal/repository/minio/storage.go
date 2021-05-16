package minio

import (
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type StorageMinio struct {
	client *minio.Client
	log    *logrus.Entry
	config Config
}

type Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Timeout   time.Duration
}

func New(config Config, log *logrus.Entry) (*StorageMinio, error) {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	return &StorageMinio{
		config: config,
		client: client,
		log:    log,
	}, nil
}

func (s *StorageMinio) PutFile(
	ctx context.Context, bucketName, objectName, contentType string, reader io.Reader,
) error {
	childCtx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()

	if _, err := s.client.PutObject(childCtx, bucketName, objectName, reader, -1, minio.PutObjectOptions{
		ContentType: contentType,
	}); err != nil {
		return errors.Wrap(err, "failed to put to minio")
	}

	return nil
}

func (s *StorageMinio) GetFileURL(
	ctx context.Context, bucket, objectName string, expires time.Duration,
) (url string, err error) {
	childCtx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()

	u, err := s.client.PresignedGetObject(childCtx, bucket, objectName, expires, nil)
	if err != nil {
		return "", err
	}

	return u.String(), nil
}

func (s *StorageMinio) DeleteFile(ctx context.Context, bucket, objectName string) error {
	childCtx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()

	if err := s.client.RemoveObject(childCtx, bucket, objectName, minio.RemoveObjectOptions{}); err != nil {
		return err
	}

	return nil
}
