package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (ac *AppController) Detail(c *gin.Context) {

	appID := c.Param("appID")

	appDetail, err := ac.appService.AppDetail(c, appID)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, appDetail)
}
