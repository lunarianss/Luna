package storage_interface

import "context"

type StorageProvider string

var (
	MINIO StorageProvider = "minio"
)

type ObjectStorage interface {
	Save(ctx context.Context, filename string, data []byte) error
	LoadOnce(ctx context.Context, filename string) ([]byte, error)
	LoadStream(ctx context.Context, filename string) (chan []byte, error)
	Download(ctx context.Context, filename string, targetPath string) error
	Exists(ctx context.Context, filename string) (bool, error)
	Delete(ctx context.Context, filename string) error
}
