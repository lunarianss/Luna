// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import (
	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/provider"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (mc *ModelController) SaveModelCredential(c *gin.Context) {
	paramsUri := &dto.CreateModelCredentialUri{}
	paramsBody := &dto.CreateModelCredentialBody{}

	if err := c.ShouldBindUri(paramsUri); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	if err := c.ShouldBindJSON(paramsBody); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	if err := mc.modelProviderService.SaveModelCredentials(c, "9ecdc361-cbc1-4c9b-8fb9-827dff4c145a", paramsBody.Model, paramsBody.ModelType, paramsUri.Provider, paramsBody.Credentials); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, core.GetSuccessResponse())
}
