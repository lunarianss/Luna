package service

import (
	"context"

	domain "github.com/lunarianss/Luna/internal/api-server/domain/dataset"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/dataset"
)

type DatasetService struct {
	datasetDomain *domain.DatasetDomain
}

func NewDatasetService(domain *domain.DatasetDomain) *DatasetService {
	return &DatasetService{
		datasetDomain: domain,
	}
}

func (s *DatasetService) GetFileUploadConfiguration(ctx context.Context) (*dto.FileUploadConfigurationResponse, error) {
	return s.datasetDomain.GetFileUploadConfiguration(ctx)

}
