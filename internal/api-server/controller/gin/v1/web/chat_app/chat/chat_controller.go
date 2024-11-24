package controller

import "github.com/lunarianss/Luna/internal/api-server/service"

type WebChatController struct {
	webAppService *service.WebChatService
}

func NewWebChatController(webChatService *service.WebChatService) *WebChatController {
	return &WebChatController{
		webAppService: webChatService,
	}
}
