package miniostorage

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

//func (s *StorageMinio) DeleteFile() error {
//	return nil
//}

//func (s *MinioStorage) Stat(ctx context.Context, bucket, objectName string) (*minio.ObjectInfo, error) {
//	childCtx, cancel := context.WithTimeout(ctx, s.config.Timeout)
//	defer cancel()
//
//	objectInfo, err := s.client.StatObject(childCtx, bucket, objectName, minio.StatObjectOptions{})
//	if err != nil {
//		return nil, err
//	}
//
//	return &objectInfo, nil
//}
//
//func (s *MinioStorage) Get(ctx context.Context, bucket, objectName string) (*minio.Object, error) {
//	s.log.Debugf("getting file %s", objectName)
//
//	childCtx, cancel := context.WithTimeout(ctx, s.config.Timeout)
//	defer cancel()
//
//	object, err := s.client.GetObject(childCtx, bucket, objectName, minio.GetObjectOptions{})
//	if err != nil {
//		return nil, err
//	}
//
//	return object, nil
//}
//
//func (s *StorageMinio) Put(
//	ctx context.Context, bucket, objectName string, reader io.Reader, objectSize int64, options Options,
//) error {
//
//	childCtx, cancel := context.WithTimeout(ctx, s.config.Timeout)
//	defer cancel()
//
//	if _, err := s.client.PutObject(childCtx, bucket, objectName, reader, objectSize, minio.PutObjectOptions{
//		ContentType: options.ContentType,
//	}); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (s *MinioStorage) FPut(ctx context.Context, bucket, objectName, filePath string, options Options) error {
//	s.log.Debugf("storing file %s", objectName)
//
//	childCtx, cancel := context.WithTimeout(ctx, s.config.Timeout)
//	defer cancel()
//
//	if _, err := s.client.FPutObject(childCtx, bucket, objectName, filePath, minio.PutObjectOptions{
//		ContentType: options.ContentType,
//	}); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (s *StorageMinio) PresignedGet(
//	ctx context.Context, bucket, objectName string, expires time.Duration, reqParams url.Values,
//) (*url.URL, error) {
//	s.log.Debugf("getting URL to file %s", objectName)
//
//	childCtx, cancel := context.WithTimeout(ctx, s.config.Timeout)
//	defer cancel()
//
//	u, err := s.client.PresignedGetObject(childCtx, bucket, objectName, expires, reqParams)
//	if err != nil {
//		return nil, err
//	}
//
//	return u, nil
//}
//
//func (s *MinioStorage) Remove(ctx context.Context, bucket, objectName string) error {
//	s.log.Debugf("removing file %s", objectName)
//
//	childCtx, cancel := context.WithTimeout(ctx, s.config.Timeout)
//	defer cancel()
//
//	if err := s.client.RemoveObject(childCtx, bucket, objectName, minio.RemoveObjectOptions{}); err != nil {
//		return err
//	}
//
//	return nil
//}
