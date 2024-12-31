// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import (
	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/agent"
	"github.com/lunarianss/Luna/internal/infrastructure/core"
)

func (tc *ToolController) GetIcon(c *gin.Context) {

	params := &dto.ListIconUri{}

	if err := c.ShouldBindUri(params); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	icon_path, err := tc.toolService.GetIconPath(c, params.Provider)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	c.File(icon_path)
}
