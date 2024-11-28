// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package route

import (
	"github.com/gin-gonic/gin"
	service "github.com/lunarianss/Luna/internal/api-server/application"
	"github.com/lunarianss/Luna/internal/api-server/config"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	controller "github.com/lunarianss/Luna/internal/api-server/interface/gin/v1/model-provider/model"
	"github.com/lunarianss/Luna/internal/api-server/middleware"
	repo_impl "github.com/lunarianss/Luna/internal/api-server/repository"
	"github.com/lunarianss/Luna/internal/pkg/mysql"
)

type ModelRoutes struct{}

func (r *ModelRoutes) Register(g *gin.Engine) error {
	gormIns, err := mysql.GetMySQLIns(nil)

	if err != nil {
		return err
	}

	// dao
	modelRepo := repo_impl.NewModelProviderRepoImpl(gormIns)
	providerRepo := repo_impl.NewProviderRepoImpl(gormIns)
	accountRepo := repo_impl.NewAccountRepoImpl(gormIns)
	tenantRepo := repo_impl.NewTenantRepoImpl(gormIns)

	// config
	config, err := config.GetLunaRuntimeConfig()

	if err != nil {
		return err
	}

	providerConfigurationsManager := domain_service.NewProviderConfigurationsManager(providerRepo, modelRepo, "", nil)
	// domain
	providerDomain := domain_service.NewProviderDomain(providerRepo, modelRepo, providerConfigurationsManager)
	accountDomain := accountDomain.NewAccountDomain(accountRepo, nil, config, nil, tenantRepo)

	// service
	modelService := service.NewModelService(providerDomain, accountDomain, config)

	modelController := controller.NewModelController(modelService)

	v1 := g.Group("/v1")
	modelProviderV1 := v1.Group("/console/api/workspaces/current")
	modelProviderV1.Use(middleware.TokenAuthMiddleware())

	modelProviderV1.POST("/model-providers/:provider/models", modelController.SaveModelCredential)
	modelProviderV1.GET("/model-providers/:provider/models/parameter-rules", modelController.ParameterRules)
	modelProviderV1.GET("/models/model-types/:modelType", modelController.GetAccountAvailableModels)

	modelProviderV1.GET("/default-model", modelController.GetDefaultModelByType)

	return nil
}

func (r *ModelRoutes) GetModule() string {
	return "model"
}
