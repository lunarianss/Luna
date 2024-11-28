// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import (
	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/app"
	"github.com/lunarianss/Luna/internal/pkg/core"
	"github.com/lunarianss/Luna/internal/pkg/util"
)

func (ac *AppController) List(c *gin.Context) {

	params := &dto.ListAppRequest{}

	if err := c.ShouldBind(params); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	userID, err := util.GetUserIDFromGin(c)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	appList, err := ac.appService.ListTenantApps(c, params, userID)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, appList)

}
