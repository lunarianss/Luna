// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package route

import (
	"github.com/gin-gonic/gin"
	service "github.com/lunarianss/Luna/internal/api-server/application"
	"github.com/lunarianss/Luna/internal/api-server/config"
	"github.com/lunarianss/Luna/internal/api-server/core/tools"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	agentDomain "github.com/lunarianss/Luna/internal/api-server/domain/agent/domain_service"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	controller "github.com/lunarianss/Luna/internal/api-server/interface/gin/v1/tool"
	"github.com/lunarianss/Luna/internal/api-server/middleware"
	repo_impl "github.com/lunarianss/Luna/internal/api-server/repository"
	"github.com/lunarianss/Luna/internal/infrastructure/mysql"
)

type ToolRoutes struct{}

func (a *ToolRoutes) Register(g *gin.Engine) error {
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
	webAppRepo := repo_impl.NewWebAppRepoImpl(gormIns)
	agentRepo := repo_impl.NewAgentRepoImpl(gormIns)

	appDomain := appDomain.NewAppDomain(appRepo, webAppRepo, gormIns)
	accountDomain := accountDomain.NewAccountDomain(accountRepo, nil, nil, nil, tenantRepo)

	agentDomain := agentDomain.NewAgentDomain(agentDomain.NewToolTransformService(config), tools.NewToolManager(), agentRepo)
	// service
	toolService := service.NewToolService(accountDomain, appDomain, agentDomain)

	toolController := controller.NewToolController(toolService)

	v1 := g.Group("/v1")
	toolV1 := v1.Group("/console/api/workspaces/current")
	unAuthV1 := v1.Group("/console/api/workspaces/current")
	toolV1.Use(middleware.TokenAuthMiddleware())
	toolV1.GET("/tools/builtin", toolController.List)
	toolV1.GET("/tools/api", toolController.ListAPI)
	toolV1.GET("/tools/workflow", toolController.ListAPI)
	toolV1.GET("/tool-labels", toolController.ListLabels)
	unAuthV1.GET("/tool-provider/builtin/:provider/icon", toolController.GetIcon)
	return nil
}

func (r *ToolRoutes) GetModule() string {
	return "tool"
}
