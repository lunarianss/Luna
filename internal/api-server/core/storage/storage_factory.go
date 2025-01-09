package storage

import (
	"context"

	storage_interface "github.com/lunarianss/Luna/internal/api-server/core/storage/interface"
	minio_storage "github.com/lunarianss/Luna/internal/api-server/core/storage/minio"
	"github.com/lunarianss/Luna/internal/infrastructure/minio"
)

type ILunaStorage interface {
	Save(ctx context.Context, filename string, data []byte) error
	LoadOnce(ctx context.Context, filename string) ([]byte, error)
	LoadStream(ctx context.Context, filename string) (chan []byte, error)
	Download(ctx context.Context, filename string, targetPath string) error
	Exists(ctx context.Context, filename string) (bool, error)
	Delete(ctx context.Context, filename string) error
}

type LunaStorage struct {
	Processor storage_interface.ObjectStorage
}

func NewStorage(ctx context.Context, bucketName string, storageName storage_interface.StorageProvider) (ILunaStorage, error) {

	lunaStorage := &LunaStorage{}

	if err := lunaStorage.InitStorageProcessor(ctx, bucketName, storageName); err != nil {
		return nil, err
	}
	return lunaStorage, nil
}

func (ls *LunaStorage) InitStorageProcessor(ctx context.Context, bucketName string, storageName storage_interface.StorageProvider) error {

	if storageName == storage_interface.MINIO {
		client, err := minio.GetMinioClient(nil)

		if err != nil {
			return err
		}
		ls.Processor = minio_storage.NewMinioStorage(client, bucketName)
	}
	return nil

}

func (ls *LunaStorage) Save(ctx context.Context, filename string, data []byte) error {
	return ls.Processor.Save(ctx, filename, data)
}

func (ls *LunaStorage) LoadOnce(ctx context.Context, filename string) ([]byte, error) {
	return ls.Processor.LoadOnce(ctx, filename)
}
func (ls *LunaStorage) LoadStream(ctx context.Context, filename string) (chan []byte, error) {
	return ls.Processor.LoadStream(ctx, filename)
}
func (ls *LunaStorage) Download(ctx context.Context, filename string, targetPath string) error {
	return ls.Processor.Download(ctx, filename, targetPath)
}
func (ls *LunaStorage) Exists(ctx context.Context, filename string) (bool, error) {
	return ls.Processor.Exists(ctx, filename)
}
func (ls *LunaStorage) Delete(ctx context.Context, filename string) error {
	return ls.Processor.Delete(ctx, filename)
}
