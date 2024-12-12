// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import (
	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/provider"
	"github.com/lunarianss/Luna/internal/infrastructure/core"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
)

func (mc *ModelController) GetProviderModels(c *gin.Context) {
	paramsUri := &dto.CreateModelCredentialUri{}

	if err := c.ShouldBindUri(paramsUri); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	userID, err := util.GetUserIDFromGin(c)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	providerModels, err := mc.modelProviderService.GetProviderModels(c, userID, paramsUri.Provider)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, providerModels)
}
