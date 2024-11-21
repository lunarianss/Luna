package service

import (
	"context"

	domain "github.com/lunarianss/Luna/internal/api-server/domain/dataset"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/dataset"
)

type DatasetService struct {
	DatasetDomain *domain.DatasetDomain
}

func NewDatasetService(domain *domain.DatasetDomain) *DatasetService {
	return &DatasetService{
		DatasetDomain: domain,
	}
}

func (s *DatasetService) GetFileUploadConfiguration(ctx context.Context) (*dto.FileUploadConfigurationResponse, error) {
	return s.DatasetDomain.GetFileUploadConfiguration(ctx)

}
