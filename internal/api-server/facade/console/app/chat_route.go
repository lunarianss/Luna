// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package route

import (
	"github.com/gin-gonic/gin"
	service "github.com/lunarianss/Luna/internal/api-server/application"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	controller "github.com/lunarianss/Luna/internal/api-server/interface/gin/v1/chat"
	"github.com/lunarianss/Luna/internal/api-server/middleware"
	repo_impl "github.com/lunarianss/Luna/internal/api-server/repository"
	"github.com/lunarianss/Luna/internal/infrastructure/mysql"
)

type ChatRoutes struct{}

func (a *ChatRoutes) Register(g *gin.Engine) error {
	gormIns, err := mysql.GetMySQLIns(nil)

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
	annotationRepo := repo_impl.NewAnnotationRepoImpl(gormIns)
	providerConfigurationsManager := domain_service.NewProviderConfigurationsManager(providerRepo, modelProviderRepo, "", nil)

	// domain
	providerDomain := domain_service.NewProviderDomain(providerRepo, modelProviderRepo, tenantRepo, providerConfigurationsManager)
	appDomain := appDomain.NewAppDomain(appRepo, webAppRepo, gormIns)
	accountDomain := accountDomain.NewAccountDomain(accountRepo, nil, nil, nil, tenantRepo)
	chatDomain := chatDomain.NewChatDomain(messageRepo, annotationRepo)

	// service
	chatService := service.NewChatService(appDomain, providerDomain, accountDomain, chatDomain)
	annotationService := service.NewAnnotationService(appDomain, providerDomain, accountDomain, chatDomain)
	chatController := controller.NewChatController(chatService, annotationService)

	v1 := g.Group("/v1")
	modelProviderV1 := v1.Group("/console/api")
	modelProviderV1.Use(middleware.TokenAuthMiddleware())
	modelProviderV1.POST("/apps/:appID/chat-messages", chatController.ChatMessage)
	modelProviderV1.POST("/apps/:appID/audio-to-text", chatController.AudioToChatMessage)
	modelProviderV1.POST("/apps/:appID/text-to-audio", chatController.TextToAudio)
	modelProviderV1.POST("/apps/:appID/annotations", chatController.InsertAnnotationFormMessage)

	modelProviderV1.GET("/apps/:appID/chat-messages", chatController.ChatMessageList)
	modelProviderV1.GET("/apps/:appID/chat-conversations", chatController.ChatConversationList)
	modelProviderV1.GET("/apps/:appID/annotations/count", chatController.GetAnnotationCount)
	modelProviderV1.GET("/apps/:appID/chat-conversations/:conversationID", chatController.ConsoleConversationDetail)
	return nil
}

func (r *ChatRoutes) GetModule() string {
	return "chat"
}
