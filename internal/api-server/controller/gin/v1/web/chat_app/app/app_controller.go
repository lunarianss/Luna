// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

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
