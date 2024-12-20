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
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	webAppDomain "github.com/lunarianss/Luna/internal/api-server/domain/web_app/domain_service"
	controller "github.com/lunarianss/Luna/internal/api-server/interface/gin/v1/web/chat_app/chat"
	"github.com/lunarianss/Luna/internal/api-server/middleware"
	repo_impl "github.com/lunarianss/Luna/internal/api-server/repository"
	"github.com/lunarianss/Luna/internal/infrastructure/email"
	"github.com/lunarianss/Luna/internal/infrastructure/mysql"
	"github.com/lunarianss/Luna/internal/infrastructure/redis"
)

type WebChatRoutes struct{}

func (a *WebChatRoutes) Register(g *gin.Engine) error {

	gormIns, err := mysql.GetMySQLIns(nil)

	if err != nil {
		return err
	}

	redisIns, err := redis.GetRedisIns(nil)

	if err != nil {
		return err
	}

	email, err := email.GetEmailSMTPIns(nil)

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
	webAppRepo := repo_impl.NewWebAppRepoImpl(gormIns)
	messageRepo := repo_impl.NewMessageRepoImpl(gormIns)
	annotationRepo := repo_impl.NewAnnotationRepoImpl(gormIns)
	providerRepo := repo_impl.NewProviderRepoImpl(gormIns)
	modelProviderRepo := repo_impl.NewModelProviderRepoImpl(gormIns)
	providerConfigurationsManager := domain_service.NewProviderConfigurationsManager(providerRepo, modelProviderRepo, "", nil)

	// domain
	appDomain := appDomain.NewAppDomain(appRepo, webAppRepo, gormIns)
	webAppDomain := webAppDomain.NewWebAppDomain(webAppRepo)
	accountDomain := accountDomain.NewAccountDomain(accountRepo, redisIns, config, email, tenantRepo)
	chatDomain := chatDomain.NewChatDomain(messageRepo, annotationRepo)

	// domain
	providerDomain := domain_service.NewProviderDomain(providerRepo, modelProviderRepo, tenantRepo, providerConfigurationsManager)
	webChatService := service.NewWebChatService(webAppDomain, accountDomain, appDomain, config, providerDomain, chatDomain)

	webSiteController := controller.NewWebChatController(webChatService)
	v1 := g.Group("/v1")
	authV1 := v1.Group("/api")
	authV1.Use(middleware.WebTokenAuthMiddleware())
	authV1.POST("/chat-messages", webSiteController.Chat)
	authV1.POST("/audio-to-text", webSiteController.AudioToChatMessage)
	authV1.POST("/text-to-audio", webSiteController.TextToAudio)
	return nil
}

func (r *WebChatRoutes) GetModule() string {
	return "web_chat"
}
