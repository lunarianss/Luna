package controller

import "github.com/lunarianss/Luna/internal/api-server/service"

type MessageController struct {
	webAppService *service.WebMessageService
}

func NewMessageController(webMessageService *service.WebMessageService) *MessageController {
	return &MessageController{
		webAppService: webMessageService,
	}
}
