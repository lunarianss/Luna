package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/infrastructure/log"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
	"github.com/lunarianss/Luna/internal/infrastructure/core"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
)

func (ac *ChatController) TextToAudio(c *gin.Context) {
	params := &dto.TextToAudioRequest{}
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

	if err := ac.chatService.TextToAudio(c, paramsUrl.AppID, params.Text, params.MessageID, "", userID); err != nil {
		log.Errorf("%#+v", err)
	}
}
