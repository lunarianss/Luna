package controller

import (
	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/provider"
	"github.com/lunarianss/Luna/internal/api-server/pkg/util"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (mc *ModelController) ParameterRules(c *gin.Context) {
	paramsUri := &dto.CreateModelCredentialUri{}
	paramsQuery := &dto.ParameterRulesQuery{}

	if err := c.ShouldBindUri(paramsUri); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	if err := c.ShouldBind(paramsQuery); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	userID, err := util.GetUserIDFromGin(c)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	parameterRules, err := mc.ModelProviderService.GetModelParameterRules(c, userID, paramsUri.Provider, paramsQuery.Model)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, parameterRules)
}
