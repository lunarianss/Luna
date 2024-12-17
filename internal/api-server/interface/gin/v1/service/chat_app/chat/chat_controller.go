// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import service "github.com/lunarianss/Luna/internal/api-server/application"

type ServiceChatController struct {
	webAppService *service.WebChatService
}

func NewServiceChatController(webChatService *service.WebChatService) *ServiceChatController {
	return &ServiceChatController{
		webAppService: webChatService,
	}
}
