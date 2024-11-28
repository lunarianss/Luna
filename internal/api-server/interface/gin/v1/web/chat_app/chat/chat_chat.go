// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import (
	"github.com/gin-gonic/gin"
	biz_entity_app_generate "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_generate"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
	"github.com/lunarianss/Luna/internal/pkg/core"
	"github.com/lunarianss/Luna/internal/pkg/util"
)

func (cc *WebChatController) Chat(c *gin.Context) {
	params := &dto.CreateChatMessageBody{}

	if err := c.ShouldBind(params); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	appID, _, endUserID, err := util.GetWebAppFromGin(c)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	if err := cc.webAppService.Chat(c, appID, endUserID, params, biz_entity_app_generate.WebApp, true); err != nil {
		core.WriteResponse(c, err, nil)
	}
}
