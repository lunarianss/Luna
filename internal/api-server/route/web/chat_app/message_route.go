// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/api-server/config"
	controller "github.com/lunarianss/Luna/internal/api-server/controller/gin/v1/web/chat_app/message"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	webAppDomain "github.com/lunarianss/Luna/internal/api-server/domain/web_app/domain_service"
	repo_impl "github.com/lunarianss/Luna/internal/api-server/repository"

	service "github.com/lunarianss/Luna/internal/api-server/application"
	"github.com/lunarianss/Luna/internal/api-server/middleware"
	"github.com/lunarianss/Luna/internal/pkg/email"
	"github.com/lunarianss/Luna/internal/pkg/mysql"
	"github.com/lunarianss/Luna/internal/pkg/redis"
)

type WebMessageRoutes struct{}

func (a *WebMessageRoutes) Register(g *gin.Engine) error {

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
	messageRepo := repo_impl.NewMessageRepoImpl(gormIns)
	providerRepo := repo_impl.NewProviderRepoImpl(gormIns)
	webAppRepo := repo_impl.NewWebAppRepoImpl(gormIns)
	modelProviderRepo := repo_impl.NewModelProviderRepoImpl(gormIns)
	providerConfigurationsManager := domain_service.NewProviderConfigurationsManager(providerRepo, modelProviderRepo, "", nil)

	// domain
	providerDomain := domain_service.NewProviderDomain(providerRepo, modelProviderRepo, providerConfigurationsManager)
	appDomain := appDomain.NewAppDomain(appRepo, webAppRepo, gormIns)
	accountDomain := accountDomain.NewAccountDomain(accountRepo, redisIns, config, email, tenantRepo)
	chatDomain := chatDomain.NewChatDomain(messageRepo)
	webAppDomain := webAppDomain.NewWebAppDomain(webAppRepo)

	webMessageService := service.NewWebMessageService(webAppDomain, accountDomain, appDomain, config, providerDomain, chatDomain)

	webSiteController := controller.NewMessageController(webMessageService)
	v1 := g.Group("/v1")
	authV1 := v1.Group("/api")
	authV1.Use(middleware.WebTokenAuthMiddleware())

	authV1.GET("/messages", webSiteController.ListMessages)
	authV1.GET("/conversations", webSiteController.ListConversation)
	authV1.PATCH("/conversations/:conversationID/pin", webSiteController.PinnedConversion)
	authV1.PATCH("/conversations/:conversationID/unpin", webSiteController.UnPinnedConversation)
	authV1.DELETE("/conversations/:conversationID", webSiteController.DeleteConversion)
	authV1.POST("/conversations/:conversationID/name", webSiteController.RenameConversion)
	return nil
}

func (r *WebMessageRoutes) GetModule() string {
	return "message"
}
