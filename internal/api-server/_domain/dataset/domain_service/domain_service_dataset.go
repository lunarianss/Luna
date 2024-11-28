package domain_service

import (
	"context"

	dto "github.com/lunarianss/Luna/internal/api-server/dto/dataset"
)

type DatasetDomain struct {
}

func NewDatasetDomain() *DatasetDomain {
	return &DatasetDomain{}
}

func (d *DatasetDomain) GetFileUploadConfiguration(ctx context.Context) (*dto.FileUploadConfigurationResponse, error) {
	return dto.NewFileUploadConfigurationResponse(), nil
}
