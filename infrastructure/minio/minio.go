package s3_minio

import (
	"context"

	"github.com/fatih/color"
	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/internal/infrastructure/options"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioS3Storage struct {
}

func New(opt *options.MinioOptions) (*minio.Client, error) {
	minioClient, err := minio.New(opt.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(opt.AccessKey, opt.SecretKey, ""),
		Secure: opt.UseSSL,
	})

	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, opt.Bucket)

	if err != nil {
		return nil, err
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, opt.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	log.Info(color.GreenString("minio is ready!"))

	return minioClient, nil
}
