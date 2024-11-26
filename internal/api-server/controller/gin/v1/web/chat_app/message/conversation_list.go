package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/api-server/core/app/apps/entities"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/web_app"
	"github.com/lunarianss/Luna/internal/api-server/pkg/util"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (mc *MessageController) ListConversation(c *gin.Context) {
	params := &dto.ListConversationQuery{}

	if err := c.ShouldBind(params); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	appID, _, endUserID, err := util.GetWebAppFromGin(c)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	conversations, err := mc.webAppService.ListConversations(c, appID, endUserID, params, entities.WEB_APP)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, conversations)
}
