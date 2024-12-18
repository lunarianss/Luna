// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import (
	"github.com/gin-gonic/gin"
	biz_entity_app_generate "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_generate"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
	"github.com/lunarianss/Luna/internal/infrastructure/core"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
)

func (cc *ServiceChatController) Chat(c *gin.Context) {
	params := &dto.ServiceCreateChatMessageBody{}

	if err := c.ShouldBind(params); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	app, tenant, err := util.GetServiceTokenFromGin(c)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	if params.ResponseMode == "streaming" {
		if err := cc.serviceChatService.Chat(c, app, tenant, params, biz_entity_app_generate.ServiceAPI); err != nil {
			core.WriteResponse(c, err, nil)
		}
		return
	}

	if params.ResponseMode == "blocking" {
		if llmResult, err := cc.serviceChatService.ChatNonStream(c, app, tenant, params, biz_entity_app_generate.ServiceAPI); err != nil {
			core.WriteResponse(c, err, nil)
		} else {
			core.WriteResponse(c, nil, llmResult)
		}
		return
	}
}
