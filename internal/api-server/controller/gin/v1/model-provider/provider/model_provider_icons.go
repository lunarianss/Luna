package controller

import (
	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/provider"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (mc *ModelProviderController) ListIcons(c *gin.Context) {

	params := &dto.ListIconRequest{}

	if err := c.ShouldBindUri(params); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	filePath, err := mc.modelProviderService.GetProviderIconPath(c, params.Provider, params.IconType, params.Lang)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	c.File(filePath)
}
