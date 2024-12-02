package controller

// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

import (
	service "github.com/lunarianss/Luna/internal/api-server/application"
)

type ChatController struct {
	chatService *service.ChatService
}

func NewAppController(chatSrv *service.ChatService) *ChatController {
	return &ChatController{chatService: chatSrv}
}
