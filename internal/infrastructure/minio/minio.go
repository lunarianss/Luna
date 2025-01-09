package minio

import (
	"fmt"
	"sync"

	"github.com/lunarianss/Luna/infrastructure/log"
	s3_minio "github.com/lunarianss/Luna/infrastructure/minio"
	"github.com/lunarianss/Luna/internal/infrastructure/options"
	"github.com/minio/minio-go/v7"
)

var (
	once        sync.Once
	MinioClient *minio.Client
)

func GetMinioClient(opt *options.MinioOptions) (*minio.Client, error) {

	var (
		err         error
		minioClient *minio.Client
	)
	once.Do(func() {
		minioClient, err = s3_minio.New(opt)

		if err != nil {
			log.Error(err)
		}

		MinioClient = minioClient
	})

	if MinioClient == nil || err != nil {
		return nil, fmt.Errorf("failed to get minio client, minioFactory: %+v, error: %w", MinioClient, err)
	}
	return MinioClient, nil
}
