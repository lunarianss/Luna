package controller

import (
	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/auth"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (ac *AuthController) EmailValidity(c *gin.Context) {
	params := &dto.EmailCodeValidityRequest{}

	if err := c.ShouldBind(params); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	tokenPair, err := ac.authService.EmailCodeValidity(c, params.Email, params.Code, params.Token)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, tokenPair)
}
