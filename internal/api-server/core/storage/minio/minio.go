package minio_storage

import (
	"bytes"
	"context"
	"io"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/infrastructure/log"
	storage_interface "github.com/lunarianss/Luna/internal/api-server/core/storage/interface"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/minio/minio-go/v7"
)

type MinioStorage struct {
	client     *minio.Client
	bucketName string
}

func NewMinioStorage(client *minio.Client, bucketName string) storage_interface.ObjectStorage {
	return &MinioStorage{
		client:     client,
		bucketName: bucketName,
	}
}

var _ storage_interface.ObjectStorage = (*MinioStorage)(nil)

func (mc *MinioStorage) Save(ctx context.Context, filename string, data []byte) error {
	if _, err := mc.client.PutObject(ctx, mc.bucketName, filename, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{}); err != nil {
		return errors.WithSCode(code.ErrMinio, err.Error())
	}
	return nil
}

func (mc *MinioStorage) LoadOnce(ctx context.Context, filename string) ([]byte, error) {
	object, err := mc.client.GetObject(ctx, mc.bucketName, filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.WithSCode(code.ErrMinio, err.Error())
	}
	defer object.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, object)

	if err != nil {
		return nil, errors.WithSCode(code.ErrMinio, err.Error())
	}

	return buf.Bytes(), nil
}

func (mc *MinioStorage) LoadStream(ctx context.Context, filename string) (chan []byte, error) {
	object, err := mc.client.GetObject(ctx, mc.bucketName, filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.WithSCode(code.ErrMinio, err.Error())
	}

	chunks := make(chan []byte, 1)

	go func() {
		for {
			buf := make([]byte, 1024*1025*5)
			n, err := object.Read(buf)

			if n > 0 {
				chunks <- buf[:n]
			}

			if err != nil {
				close(chunks)
				if !errors.Is(err, io.EOF) {
					log.Errorf("minio load stream error: %s", err.Error())
				}
				return
			}
		}
	}()

	return chunks, nil
}

func (mc *MinioStorage) Download(ctx context.Context, filename string, targetPath string) error {
	err := mc.client.FGetObject(ctx, mc.bucketName, filename, targetPath, minio.GetObjectOptions{})
	return err
}

func (mc *MinioStorage) Exists(ctx context.Context, filename string) (bool, error) {
	_, err := mc.client.StatObject(ctx, mc.bucketName, filename, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (mc *MinioStorage) Delete(ctx context.Context, filename string) error {
	return nil
}
