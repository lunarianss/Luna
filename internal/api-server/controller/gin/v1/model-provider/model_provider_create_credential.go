package controller

import (
	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/provider"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (mc *ModelProviderController) CreateProviderCredential(c *gin.Context) {
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

	if err := mc.modelProviderService.CreateProviderCredentials("9ecdc361-cbc1-4c9b-8fb9-827dff4c145a", paramsUri.Provider, paramsBody.Credentials); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, core.GetSuccessResponse())
}
