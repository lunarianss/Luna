// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package route

import (
	"github.com/gin-gonic/gin"
	service "github.com/lunarianss/Luna/internal/api-server/application"
	"github.com/lunarianss/Luna/internal/api-server/config"
	"github.com/lunarianss/Luna/internal/api-server/core/tools"
	agentDomain "github.com/lunarianss/Luna/internal/api-server/domain/agent/domain_service"
	controller "github.com/lunarianss/Luna/internal/api-server/interface/gin/v1/files"
	repo_impl "github.com/lunarianss/Luna/internal/api-server/repository"
	"github.com/lunarianss/Luna/internal/infrastructure/mysql"
)

type ToolFilesRoutes struct{}

func (a *ToolFilesRoutes) Register(g *gin.Engine) error {
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

	agentRepo := repo_impl.NewAgentRepoImpl(gormIns)
	// domain
	agentDomain := agentDomain.NewAgentDomain(agentDomain.NewToolTransformService(config), tools.NewToolManager(), agentRepo)

	// service

	fileService := service.NewFileService(agentDomain, config)

	fileController := controller.NewFileController(fileService)
	v1 := g.Group("/v1")
	v1.GET("/files/tools/:filename", fileController.PreviewFile)
	return nil
}

func (r *ToolFilesRoutes) GetModule() string {
	return "files/tool"
}
