// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

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
