package controller

import "github.com/lunarianss/Luna/internal/api-server/service"

type AccountController struct {
	accountService *service.AccountService
}

func NewAccountController(accountService *service.AccountService) *AccountController {
	return &AccountController{
		accountService: accountService,
	}
}
