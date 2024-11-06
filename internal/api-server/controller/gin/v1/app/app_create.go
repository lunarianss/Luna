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

	app, err := ac.AppService.CreateApp(c, "9ecdc361-cbc1-4c9b-8fb9-827dff4c145a", "8ecdc361-cbc1-4c9b-8fb9-827dff4c145a", params)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, app)
}
