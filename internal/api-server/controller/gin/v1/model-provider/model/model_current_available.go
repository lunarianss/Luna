package controller

import (
	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/provider"
	"github.com/lunarianss/Luna/internal/api-server/entities/base"
	"github.com/lunarianss/Luna/internal/api-server/pkg/util"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (mc *ModelController) GetAccountAvailableModels(c *gin.Context) {
	params := &dto.GetAccountAvailableModelsRequest{}

	if err := c.ShouldBindUri(params); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	userID, err := util.GetUserIDFromGin(c)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	availableModels, err := mc.modelProviderService.GetAccountAvailableModels(c, userID, base.ModelType(params.ModelType))

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, availableModels)
}
