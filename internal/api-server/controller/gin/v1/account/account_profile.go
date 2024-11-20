package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/api-server/pkg/util"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (ac *AccountController) GetAccountProfile(c *gin.Context) {

	userID, err := util.GetUserIDFromGin(c)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	accountProfile, err := ac.AccountService.GetAccountProfile(c, userID)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, accountProfile)
}
