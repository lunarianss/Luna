package controller

import (
	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/auth"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (ac *AuthController) RefreshToken(c *gin.Context) {

	params := &dto.RefreshTokenRequest{}

	if err := c.ShouldBind(params); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	sendResp, err := ac.AuthService.RefreshToken(c, params.RefreshToken)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, sendResp)
}
