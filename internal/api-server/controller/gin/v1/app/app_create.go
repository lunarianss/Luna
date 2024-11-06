package controller

import (
	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/app"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (ac *AppController) Create(c *gin.Context) {

	params := &dto.CreateAppRequest{}

	if err := c.ShouldBind(params); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	// ac.AppService.CreateApp(c)

}
