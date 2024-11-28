// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package entities

import (
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config/entities"
	"github.com/lunarianss/Luna/internal/api-server/core/app/file"
)

type InvokeForm string

const (
	SERVICE_API InvokeForm = "service-api"
	WEB_APP     InvokeForm = "web-app"
	EXPLORE     InvokeForm = "explore"
	DEBUGGER    InvokeForm = "debugger"
)

func (i InvokeForm) ToSource() string {
	if i == WEB_APP {
		return "web_app"
	} else if i == DEBUGGER {

		return "dev"
	} else if i == EXPLORE {
		return "explore_app"
	} else if i == SERVICE_API {
		return "api"
	}

	return "dev"
}

type AppGenerateEntity struct {
	TaskID     string                 `json:"task_id"`
	AppConfig  *entities.AppConfig    `json:"app_config"`
	Inputs     map[string]interface{} `json:"inputs"`
	Files      []*file.File           `json:"files"`
	UseID      string                 `json:"use_id"`
	Stream     bool                   `json:"stream"`
	InvokeFrom InvokeForm             `json:"invoke_from"`
	CallDepth  int                    `json:"call_depth"`
	Extras     map[string]interface{} `json:"extras"`
}
