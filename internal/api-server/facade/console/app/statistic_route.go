// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package route

import (
	"github.com/gin-gonic/gin"
	service "github.com/lunarianss/Luna/internal/api-server/application"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	controller "github.com/lunarianss/Luna/internal/api-server/interface/gin/v1/statistic"
	"github.com/lunarianss/Luna/internal/api-server/middleware"
	repo_impl "github.com/lunarianss/Luna/internal/api-server/repository"
	"github.com/lunarianss/Luna/internal/infrastructure/mysql"
)

type StatisticRoutes struct{}

func (a *StatisticRoutes) Register(g *gin.Engine) error {
	gormIns, err := mysql.GetMySQLIns(nil)

	if err != nil {
		return err
	}

	// repos
	accountRepo := repo_impl.NewAccountRepoImpl(gormIns)
	tenantRepo := repo_impl.NewTenantRepoImpl(gormIns)
	// appRepo := repo_impl.NewAppRepoImpl(gormIns)
	messageRepo := repo_impl.NewMessageRepoImpl(gormIns)
	// providerRepo := repo_impl.NewProviderRepoImpl(gormIns)
	// webAppRepo := repo_impl.NewWebAppRepoImpl(gormIns)
	// modelProviderRepo := repo_impl.NewModelProviderRepoImpl(gormIns)
	// providerConfigurationsManager := domain_service.NewProviderConfigurationsManager(providerRepo, modelProviderRepo, "", nil)

	// domain
	// providerDomain := domain_service.NewProviderDomain(providerRepo, modelProviderRepo, providerConfigurationsManager)
	// appDomain := appDomain.NewAppDomain(appRepo, webAppRepo, gormIns)
	accountDomain := accountDomain.NewAccountDomain(accountRepo, nil, nil, nil, tenantRepo)
	chatDomain := chatDomain.NewChatDomain(messageRepo)

	// service
	chatService := service.NewStatisticService(chatDomain, accountDomain)

	appController := controller.NewSetupController(chatService)

	v1 := g.Group("/v1")
	modelProviderV1 := v1.Group("/console/api")
	statisticsGroup := modelProviderV1.Group("/apps/:appID/statistics")
	statisticsGroup.Use(middleware.TokenAuthMiddleware())
	statisticsGroup.GET("/daily-conversations", appController.DailyConversations)
	// statisticsGroup.GET("/daily-messages", appController.ChatMessageList)
	// statisticsGroup.GET("/daily-end-users", appController.ChatConversationList)
	// statisticsGroup.GET("/token-costs", appController.GetAnnotationCount)
	// statisticsGroup.GET("/average-session-interactions", appController.ConsoleConversationDetail)
	// statisticsGroup.GET("/user-satisfaction-rate", appController.ConsoleConversationDetail)
	// statisticsGroup.GET("/average-response-time", appController.ConsoleConversationDetail)
	// statisticsGroup.GET("/tokens-per-second", appController.ConsoleConversationDetail)
	return nil
}

func (r *StatisticRoutes) GetModule() string {
	return "statistic"
}
