package controller

import "github.com/lunarianss/Luna/internal/api-server/service"

type AccountController struct {
	AccountService *service.AccountService
}

func NewAccountController(accountService *service.AccountService) *AccountController {
	return &AccountController{
		AccountService: accountService,
	}
}
