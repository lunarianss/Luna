// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import "github.com/lunarianss/Luna/internal/api-server/service"

type DatasetController struct {
	datasetService *service.DatasetService
}

func NewDatasetController(datasetSrv *service.DatasetService) *DatasetController {
	return &DatasetController{
		datasetService: datasetSrv,
	}
}
