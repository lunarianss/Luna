// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package controller

import (
	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
	"github.com/lunarianss/Luna/internal/infrastructure/core"
)

func (ac *ChatController) ChatMessageList(c *gin.Context) {

	paramsUrl := &dto.CreateChatMessageUri{}
	paramsQuery := dto.NewListChatMessageQuery()

	if err := c.ShouldBindUri(paramsUrl); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	if err := c.ShouldBind(paramsQuery); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	messages, err := ac.chatService.ListConsoleMessagesOfConversation(c, paramsUrl.AppID, paramsQuery)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, messages)
}
