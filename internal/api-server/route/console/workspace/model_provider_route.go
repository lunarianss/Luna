// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package route

import (
	"github.com/gin-gonic/gin"

	controller "github.com/lunarianss/Luna/internal/api-server/controller/gin/v1/model-provider"
	"github.com/lunarianss/Luna/internal/api-server/dao"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/model-provider"
	"github.com/lunarianss/Luna/internal/api-server/service"
	"github.com/lunarianss/Luna/internal/pkg/mysql"
)

type ModelProviderRoutes struct{}

func (r *ModelProviderRoutes) Register(g *gin.Engine) error {
	gormIns, err := mysql.GetMySQLIns(nil)

	if err != nil {
		return err
	}

	// dao
	modelProviderDao := dao.NewModelProvider(gormIns)
	// domain
	modelProviderDomain := domain.NewModelProviderDomain(modelProviderDao)

	// service
	modelProviderService := service.NewModelProviderService(modelProviderDomain)
	modelProviderController := controller.NewModelProviderController(modelProviderService)

	v1 := g.Group("/v1")
	blogV1 := v1.Group("/console/workspace/current")
	blogV1.GET("/model-providers", modelProviderController.List)
	blogV1.GET("/model-providers/:provider/:iconType/:lang", modelProviderController.ListIcons)
	return nil
}

func (r *ModelProviderRoutes) GetModule() string {
	return "providers"
}
