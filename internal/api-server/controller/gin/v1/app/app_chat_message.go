package controller

import (
	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/app"
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

	core.WriteResponse(c, nil, nil)
}
