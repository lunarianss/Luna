// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package route

import (
	"github.com/gin-gonic/gin"
	controller "github.com/lunarianss/Luna/internal/api-server/interface/gin/v1/setup"
)

type SetupRoutes struct{}

func (a *SetupRoutes) Register(g *gin.Engine) error {

	accountController := controller.NewSetupController()
	v1 := g.Group("/v1")
	authV1 := v1.Group("/console/api")
	authV1.GET("/setup", accountController.ValidateSetup)
	return nil
}

func (r *SetupRoutes) GetModule() string {
	return "setup"
}
