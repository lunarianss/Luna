// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/domain_service"
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
