package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/api-server/core/app/apps/entities"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
	"github.com/lunarianss/Luna/internal/api-server/pkg/util"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (ac *AppController) ChatMessage(c *gin.Context) {

	params := &dto.CreateChatMessageBody{}
	paramsUrl := &dto.CreateChatMessageUri{}

	if err := c.ShouldBindUri(paramsUrl); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}
	if err := c.ShouldBind(params); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	userID, err := util.GetUserIDFromGin(c)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	if err := ac.chatService.Generate(c, paramsUrl.AppID, userID, params, entities.DEBUGGER, true); err != nil {
		core.WriteResponse(c, err, nil)
	}
}
