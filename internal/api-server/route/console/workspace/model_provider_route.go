// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package route

import (
	"github.com/gin-gonic/gin"
	controller "github.com/lunarianss/Hurricane/internal/api-server/controller/gin/v1/blog"
	"github.com/lunarianss/Hurricane/internal/api-server/dao"
	domain "github.com/lunarianss/Hurricane/internal/api-server/domain/blog"
	"github.com/lunarianss/Hurricane/internal/api-server/service"
	"github.com/lunarianss/Hurricane/internal/pkg/mysql"
)

type ModelProviderRoutes struct{}

func (r *ModelProviderRoutes) Register(g *gin.Engine) error {
	gormIns, err := mysql.GetMySQLIns(nil)

	if err != nil {
		return err
	}

	// dao
	blogDao := dao.NewBlogDao(gormIns)

	// domain
	blogDomain := domain.NewBlogDomain(blogDao)

	// service
	blogService := service.NewBlogService(blogDomain)
	blogController := controller.NewBlogController(blogService)

	v1 := g.Group("/v1")
	blogV1 := v1.Group("/console/workspace/current")
	blogV1.GET("/model-providers", blogController.List)
	return nil
}

func (r *ModelProviderRoutes) GetModule() string {
	return "providers"
}
