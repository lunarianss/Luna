// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package route

import (
	"github.com/gin-gonic/gin"

	"github.com/lunarianss/Luna/internal/api-server/config"
	controller "github.com/lunarianss/Luna/internal/api-server/controller/gin/v1/model-provider/provider"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/middleware"
	repo_impl "github.com/lunarianss/Luna/internal/api-server/repository"
	"github.com/lunarianss/Luna/internal/api-server/service"
	"github.com/lunarianss/Luna/internal/pkg/mysql"
)

type ModelProviderRoutes struct{}

func (r *ModelProviderRoutes) Register(g *gin.Engine) error {
	gormIns, err := mysql.GetMySQLIns(nil)

	if err != nil {
		return err
	}

	// repos
	providerRepo := repo_impl.NewProviderRepoImpl(gormIns)
	modelProviderRepo := repo_impl.NewModelProviderRepoImpl(gormIns)
	accountRepo := repo_impl.NewAccountRepoImpl(gormIns)
	tenantRepo := repo_impl.NewTenantRepoImpl(gormIns)
	providerConfigurationsManager := domain_service.NewProviderConfigurationsManager(providerRepo, modelProviderRepo, "", nil)

	// domain
	modelProviderDomain := domain_service.NewProviderDomain(providerRepo, modelProviderRepo, providerConfigurationsManager)
	accountDomain := accountDomain.NewAccountDomain(accountRepo, nil, nil, nil, tenantRepo)

	// config
	config, err := config.GetLunaRuntimeConfig()

	if err != nil {
		return err
	}

	// service
	modelProviderService := service.NewModelProviderService(modelProviderDomain, accountDomain, config)

	modelProviderController := controller.NewModelProviderController(modelProviderService)

	v1 := g.Group("/v1")
	modelProviderAuthV1 := v1.Group("/console/api/workspaces/current")
	modelProviderNoAuthV1 := v1.Group("/console/api/workspaces/current")

	modelProviderAuthV1.Use(middleware.TokenAuthMiddleware())

	modelProviderAuthV1.GET("/model-providers", modelProviderController.List)
	modelProviderNoAuthV1.GET("/model-providers/:provider/:iconType/:lang", modelProviderController.ListIcons)
	modelProviderAuthV1.POST("/model-providers/:provider", modelProviderController.SaveProviderCredential)

	return nil
}

func (r *ModelProviderRoutes) GetModule() string {
	return "providers"
}
