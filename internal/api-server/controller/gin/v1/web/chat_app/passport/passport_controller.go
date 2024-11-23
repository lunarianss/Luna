package controller

import "github.com/lunarianss/Luna/internal/api-server/service"

type PassportController struct {
	passportService *service.PassportService
}

func NewPassportController(passportService *service.PassportService) *PassportController {
	return &PassportController{
		passportService: passportService,
	}
}
