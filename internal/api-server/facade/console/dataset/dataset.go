// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package route

import (
	"github.com/gin-gonic/gin"
	service "github.com/lunarianss/Luna/internal/api-server/application"
	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/domain_service"
	controller "github.com/lunarianss/Luna/internal/api-server/interface/gin/v1/dataset"
	"github.com/lunarianss/Luna/internal/api-server/middleware"
	repo_impl "github.com/lunarianss/Luna/internal/api-server/repository"
	"github.com/lunarianss/Luna/internal/infrastructure/mysql"
)

type DatasetRoutes struct {
}

func (a *DatasetRoutes) Register(g *gin.Engine) error {
	gormIns, err := mysql.GetMySQLIns(nil)

	if err != nil {
		return err
	}

	// redisIns, err := redis.GetRedisIns(nil)

	// if err != nil {
	// 	return err
	// }

	// email, err := email.GetEmailSMTPIns(nil)

	// if err != nil {
	// 	return err
	// }

	// // config
	// config, err := config.GetLunaRuntimeConfig()

	// if err != nil {
	// 	return err
	// }

	datasetRepo := repo_impl.NewDatasetRepoImpl(gormIns)

	// domain

	datasetDomain := domain_service.NewDatasetDomain(datasetRepo)
	// service
	datasetService := service.NewDatasetService(datasetDomain)
	datasetController := controller.NewDatasetController(datasetService)

	v1 := g.Group("/v1")
	authV1 := v1.Group("/console/api")
	authV1.Use(middleware.TokenAuthMiddleware())
	authV1.GET("/files/upload", datasetController.GetFileUploadConfiguration)
	return nil
}

func (r *DatasetRoutes) GetModule() string {
	return "dataset"
}
