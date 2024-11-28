// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import (
	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/auth"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (ac *AuthController) SendEmailCode(c *gin.Context) {

	params := &dto.SendEmailCodeRequest{}

	if err := c.ShouldBind(params); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	sendResp, err := ac.authService.SendEmailCode(c, params)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, sendResp)
}
