package controller

import (
	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/provider"
	"github.com/lunarianss/Luna/internal/api-server/pkg/util"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (mc *ModelController) GetDefaultModelByType(c *gin.Context) {

	paramsQuery := &dto.DefaultModelByTypeQuery{}

	if err := c.ShouldBind(paramsQuery); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	userID, err := util.GetUserIDFromGin(c)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	defaultModel, err := mc.ModelProviderService.GetDefaultModelByType(c, userID, paramsQuery.ModelType)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, defaultModel)
}
