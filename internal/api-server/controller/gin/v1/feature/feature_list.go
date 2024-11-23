package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (fc *FeatureController) List(c *gin.Context) {
	features, err := fc.featureService.ListFeatures()

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, features)
}
