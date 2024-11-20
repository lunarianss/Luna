package controller

import (
	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/provider"
	"github.com/lunarianss/Luna/internal/api-server/pkg/util"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (mc *ModelProviderController) SaveProviderCredential(c *gin.Context) {
	paramsUri := &dto.CreateProviderCredentialUri{}
	paramsBody := &dto.CreateProviderCredentialBody{}

	if err := c.ShouldBindUri(paramsUri); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	if err := c.ShouldBindJSON(paramsBody); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	userID, err := util.GetUserIDFromGin(c)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	if err := mc.modelProviderService.SaveProviderCredentials(c, userID, paramsUri.Provider, paramsBody.Credentials); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, core.GetSuccessResponse())
}
