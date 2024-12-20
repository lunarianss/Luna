// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package route

import (
	"github.com/gin-gonic/gin"
	service "github.com/lunarianss/Luna/internal/api-server/application"
	"github.com/lunarianss/Luna/internal/api-server/config"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	controller "github.com/lunarianss/Luna/internal/api-server/interface/gin/v1/workspace"
	"github.com/lunarianss/Luna/internal/api-server/middleware"
	repo_impl "github.com/lunarianss/Luna/internal/api-server/repository"
	"github.com/lunarianss/Luna/internal/infrastructure/email"
	"github.com/lunarianss/Luna/internal/infrastructure/mysql"
	"github.com/lunarianss/Luna/internal/infrastructure/redis"
)

type WorkspaceRoutes struct{}

func (a *WorkspaceRoutes) Register(g *gin.Engine) error {
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

	// domain

	accountDomain := accountDomain.NewAccountDomain(accountRepo, redisIns, config, email, tenantRepo)
	// service
	tenantService := service.NewTenantService(accountDomain)

	workspaceController := controller.NewWorkspaceController(tenantService)
	v1 := g.Group("/v1")
	authV1 := v1.Group("/console/api")
	authV1.Use(middleware.TokenAuthMiddleware())
	authV1.GET("/workspaces", workspaceController.List)
	authV1.GET("/workspaces/current", workspaceController.GetTenantCurrentWorkspace)

	return nil
}

func (r *WorkspaceRoutes) GetModule() string {
	return "setup"
}
