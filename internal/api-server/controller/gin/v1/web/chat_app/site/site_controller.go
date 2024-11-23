package controller

import "github.com/lunarianss/Luna/internal/api-server/service"

type WebSiteController struct {
	webSiteService *service.WebSiteService
}

func NewWebSiteController(webSiteService *service.WebSiteService) *WebSiteController {
	return &WebSiteController{
		webSiteService: webSiteService,
	}
}
