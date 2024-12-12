// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package route

import (
	"github.com/gin-gonic/gin"
	controller "github.com/lunarianss/Luna/internal/api-server/interface/gin/v1/tag"
	"github.com/lunarianss/Luna/internal/api-server/middleware"
)

type TagRoutes struct{}

func (a *TagRoutes) Register(g *gin.Engine) error {

	accountController := controller.NewTagController()
	v1 := g.Group("/v1")
	routerPrefix := v1.Group("/console/api")
	routerPrefix.Use(middleware.TokenAuthMiddleware())
	authV1 := routerPrefix.Group("/workspaces/current")
	authV1.GET("/tool-providers", accountController.List)
	routerPrefix.GET("/tags", accountController.List)
	return nil
}

func (r *TagRoutes) GetModule() string {
	return "tag"
}
