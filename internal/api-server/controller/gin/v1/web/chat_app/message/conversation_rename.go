package controller

import (
	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/web_app"
	"github.com/lunarianss/Luna/internal/api-server/pkg/util"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (mc *MessageController) RenameConversion(c *gin.Context) {
	params := &dto.ConversationIDUrl{}
	paramsBody := dto.NewRenameConversationRequest()

	if err := c.ShouldBindUri(params); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	if err := c.ShouldBind(paramsBody); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	appID, _, endUserID, err := util.GetWebAppFromGin(c)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	if err := mc.webAppService.RenameConversation(c, appID, endUserID, params.ConversationID, paramsBody); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, core.GetSuccessResponse())
}
