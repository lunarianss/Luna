package service

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/_domain/dataset/domain_service"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/dataset"
)

type DatasetService struct {
	datasetDomain *domain_service.DatasetDomain
}

func NewDatasetService(domain *domain_service.DatasetDomain) *DatasetService {
	return &DatasetService{
		datasetDomain: domain,
	}
}

func (s *DatasetService) GetFileUploadConfiguration(ctx context.Context) (*dto.FileUploadConfigurationResponse, error) {
	return s.datasetDomain.GetFileUploadConfiguration(ctx)

}
