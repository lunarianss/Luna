package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/infrastructure/log"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
	"github.com/lunarianss/Luna/internal/infrastructure/core"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
)

func (ac *WebChatController) TextToAudio(c *gin.Context) {
	params := &dto.TextToAudioRequest{}

	if err := c.ShouldBind(params); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	appID, _, endUserID, err := util.GetWebAppFromGin(c)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	if err := ac.webAppService.TextToAudio(c, appID, params.Text, params.MessageID, "", endUserID); err != nil {
		log.Errorf("%#+v", err)
	}
}
