package controller

import (
	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
	"github.com/lunarianss/Luna/internal/infrastructure/core"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
)

func (cc *ChatController) ConsoleConversationDetail(c *gin.Context) {
	params := &dto.DetailConversationUri{}

	if err := c.ShouldBindUri(params); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}
	userID, err := util.GetUserIDFromGin(c)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	conversationDetail, err := cc.chatService.DetailConversation(c, userID, params.ConversationID, params.AppID)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, conversationDetail)

}
