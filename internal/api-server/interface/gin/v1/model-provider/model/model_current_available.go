// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import (
	"github.com/gin-gonic/gin"
	common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/provider"
	"github.com/lunarianss/Luna/internal/pkg/core"
	"github.com/lunarianss/Luna/internal/pkg/util"
)

func (mc *ModelController) GetAccountAvailableModels(c *gin.Context) {
	params := &dto.GetAccountAvailableModelsRequest{}

	if err := c.ShouldBindUri(params); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	userID, err := util.GetUserIDFromGin(c)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	availableModels, err := mc.modelProviderService.GetAccountAvailableModels(c, userID, common.ModelType(params.ModelType))

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, availableModels)
}
