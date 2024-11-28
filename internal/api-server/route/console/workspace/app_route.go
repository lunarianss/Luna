// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package route

import (
	"github.com/gin-gonic/gin"
	service "github.com/lunarianss/Luna/internal/api-server/application"
	"github.com/lunarianss/Luna/internal/api-server/config"
	controller "github.com/lunarianss/Luna/internal/api-server/controller/gin/v1/app"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/middleware"
	repo_impl "github.com/lunarianss/Luna/internal/api-server/repository"
	"github.com/lunarianss/Luna/internal/pkg/mysql"
)

type AppRoutes struct{}

func (a *AppRoutes) Register(g *gin.Engine) error {
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

	// domain
	providerDomain := domain_service.NewProviderDomain(providerRepo, modelProviderRepo, providerConfigurationsManager)
	appDomain := appDomain.NewAppDomain(appRepo, webAppRepo, gormIns)
	accountDomain := accountDomain.NewAccountDomain(accountRepo, nil, nil, nil, tenantRepo)
	chatDomain := chatDomain.NewChatDomain(messageRepo)

	// service
	appService := service.NewAppService(appDomain, providerDomain, accountDomain, gormIns, config)
	chatService := service.NewChatService(appDomain, providerDomain, accountDomain, chatDomain)

	appController := controller.NewAppController(appService, chatService)

	v1 := g.Group("/v1")
	modelProviderV1 := v1.Group("/console/api")
	modelProviderV1.Use(middleware.TokenAuthMiddleware())
	modelProviderV1.POST("/apps", appController.Create)
	modelProviderV1.GET("/apps", appController.List)
	modelProviderV1.GET("/apps/:appID", appController.Detail)
	modelProviderV1.POST("/apps/:appID/model-config", appController.UpdateModelConfig)
	modelProviderV1.POST("/apps/:appID/chat-messages", appController.ChatMessage)
	return nil
}

func (r *AppRoutes) GetModule() string {
	return "app"
}
