// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package route

import (
	"github.com/gin-gonic/gin"
	service "github.com/lunarianss/Luna/internal/api-server/application"
	"github.com/lunarianss/Luna/internal/api-server/config"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	datasetDomain "github.com/lunarianss/Luna/internal/api-server/domain/dataset/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	controller "github.com/lunarianss/Luna/internal/api-server/interface/gin/v1/app"
	"github.com/lunarianss/Luna/internal/api-server/middleware"
	repo_impl "github.com/lunarianss/Luna/internal/api-server/repository"
	"github.com/lunarianss/Luna/internal/infrastructure/mysql"
	"github.com/lunarianss/Luna/internal/infrastructure/redis"
)

type AppRoutes struct{}

func (a *AppRoutes) Register(g *gin.Engine) error {
	redisIns, err := redis.GetRedisIns(nil)

	if err != nil {
		return err
	}

	gormIns, err := mysql.GetMySQLIns(nil)

	if err != nil {
		return err
	}

	// config
	config, err := config.GetLunaRuntimeConfig()

	if err != nil {
		return err
	}

	// repos
	accountRepo := repo_impl.NewAccountRepoImpl(gormIns)
	tenantRepo := repo_impl.NewTenantRepoImpl(gormIns)
	appRepo := repo_impl.NewAppRepoImpl(gormIns)
	messageRepo := repo_impl.NewMessageRepoImpl(gormIns)
	providerRepo := repo_impl.NewProviderRepoImpl(gormIns)
	webAppRepo := repo_impl.NewWebAppRepoImpl(gormIns)
	modelProviderRepo := repo_impl.NewModelProviderRepoImpl(gormIns)
	providerConfigurationsManager := domain_service.NewProviderConfigurationsManager(providerRepo, modelProviderRepo, "", nil)
	annotationRepo := repo_impl.NewAnnotationRepoImpl(gormIns)
	datasetRepo := repo_impl.NewDatasetRepoImpl(gormIns)

	// domain
	datasetDomain := datasetDomain.NewDatasetDomain(datasetRepo)
	providerDomain := domain_service.NewProviderDomain(providerRepo, modelProviderRepo, tenantRepo, providerConfigurationsManager)
	appDomain := appDomain.NewAppDomain(appRepo, webAppRepo, gormIns)
	accountDomain := accountDomain.NewAccountDomain(accountRepo, nil, nil, nil, tenantRepo)
	chatDomain := chatDomain.NewChatDomain(messageRepo, annotationRepo)

	// service
	appService := service.NewAppService(appDomain, providerDomain, accountDomain, chatDomain, gormIns, config)
	chatService := service.NewChatService(appDomain, providerDomain, accountDomain, chatDomain, datasetDomain, redisIns)

	appController := controller.NewAppController(appService, chatService)

	v1 := g.Group("/v1")
	modelProviderV1 := v1.Group("/console/api")
	modelProviderV1.Use(middleware.TokenAuthMiddleware())
	modelProviderV1.POST("/apps", appController.Create)
	modelProviderV1.GET("/apps", appController.List)
	modelProviderV1.GET("/apps/:appID", appController.Detail)
	modelProviderV1.POST("/apps/:appID/model-config", appController.UpdateModelConfig)
	return nil
}

func (r *AppRoutes) GetModule() string {
	return "app"
}
