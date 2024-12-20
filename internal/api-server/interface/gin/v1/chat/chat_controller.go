package controller

// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

import (
	service "github.com/lunarianss/Luna/internal/api-server/application"
)

type ChatController struct {
	chatService *service.ChatService
	annotation  *service.AnnotationService
}

func NewChatController(chatSrv *service.ChatService, annotation *service.AnnotationService) *ChatController {
	return &ChatController{chatService: chatSrv, annotation: annotation}
}
