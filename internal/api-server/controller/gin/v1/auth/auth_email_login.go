package controller

import (
	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/auth"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (ac *AuthController) SendEmailCode(c *gin.Context) {

	params := &dto.SendEmailCodeRequest{}

	if err := c.ShouldBind(params); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	ac.AuthService.GetUserThroughEmails(c, params.Email)

}
