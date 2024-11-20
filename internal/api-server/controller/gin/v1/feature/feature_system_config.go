package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (fc *FeatureController) GetSystemConfigs(c *gin.Context) {
	systemConfig, err := fc.FeatureService.GetSystemConfig()

	if err != nil {
		core.WriteResponse(c, err, nil)
	}

	core.WriteResponse(c, nil, systemConfig)
}
