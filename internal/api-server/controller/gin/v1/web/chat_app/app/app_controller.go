package controller

import "github.com/lunarianss/Luna/internal/api-server/service"

type WebAppController struct {
	webAppService *service.WebAppService
}

func NewWebAppController(webAppService *service.WebAppService) *WebAppController {
	return &WebAppController{
		webAppService: webAppService,
	}
}
