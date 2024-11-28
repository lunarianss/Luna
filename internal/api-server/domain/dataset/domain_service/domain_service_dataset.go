// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

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
