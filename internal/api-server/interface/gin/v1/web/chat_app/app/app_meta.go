// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ac *WebAppController) AppMeta(c *gin.Context) {

	c.JSON(http.StatusOK, map[string]map[string]any{
		"tool_icons": {},
	})
}
