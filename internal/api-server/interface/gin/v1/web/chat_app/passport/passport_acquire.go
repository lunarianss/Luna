// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/lunarianss/Luna/internal/infrastructure/core"
	"github.com/lunarianss/Luna/infrastructure/errors"
)

func (pc *PassportController) Acquire(c *gin.Context) {

	appCode := c.GetHeader("X-App-Code")

	if appCode == "" {
		core.WriteResponse(c, errors.WithCode(code.ErrAppCodeNotFound, "app code not exist in X-App-Code http header"), nil)
		return
	}

	passport, err := pc.passportService.AcquirePassport(c, appCode)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, passport)
}
