package controller

import (
	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/provider"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (mc *ModelController) GetAccountAvailableModels(c *gin.Context) {
	params := &dto.GetAccountAvailableModelsRequest{}

	if err := c.ShouldBindUri(params); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	// mc.ModelProviderService.

}
