// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/api-server/config"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	controller "github.com/lunarianss/Luna/internal/api-server/interface/gin/v1/account"
	"github.com/lunarianss/Luna/internal/api-server/middleware"
	repo_impl "github.com/lunarianss/Luna/internal/api-server/repository"

	service "github.com/lunarianss/Luna/internal/api-server/application"
	"github.com/lunarianss/Luna/internal/infrastructure/email"
	"github.com/lunarianss/Luna/internal/infrastructure/mysql"
	"github.com/lunarianss/Luna/internal/infrastructure/redis"
)

type AccountRoute struct {
}

func (a *AccountRoute) Register(g *gin.Engine) error {
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

	// repo
	accountRepo := repo_impl.NewAccountRepoImpl(gormIns)
	tenantRepo := repo_impl.NewTenantRepoImpl(gormIns)

	// domain
	accountDomain := domain.NewAccountDomain(accountRepo, redisIns, config, email, tenantRepo)
	tenantDomain := domain.NewTenantDomain(tenantRepo)

	// service
	accountService := service.NewAccountService(accountDomain, tenantDomain, gormIns, nil)

	accountController := controller.NewAccountController(accountService)

	v1 := g.Group("/v1")

	authV1 := v1.Group("/console/api")
	accountV1 := authV1.Group("/account").Use(middleware.TokenAuthMiddleware())
	accountV1.GET("/profile", accountController.GetAccountProfile)
	return nil
}

func (r *AccountRoute) GetModule() string {
	return "account"
}
