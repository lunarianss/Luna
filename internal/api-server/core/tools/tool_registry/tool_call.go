// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package tool_registry

import (
	"context"
	"fmt"

	"github.com/lunarianss/Luna/infrastructure/log"

	"github.com/lunarianss/Luna/internal/api-server/domain/agent/entity/biz_entity"
)

type IModelRegistryCall interface {
	Invoke(ctx context.Context, toolParameters []byte) ([]*biz_entity.ToolInvokeMessage, error)
}

type modelRegistryCall struct {
	toolRuntime *biz_entity.ToolRuntimeConfiguration
	userID      string
}

func NewModelRegisterCaller(userID string, toolRuntime *biz_entity.ToolRuntimeConfiguration) IModelRegistryCall {
	return &modelRegistryCall{
		toolRuntime: toolRuntime,
		userID:      userID,
	}
}

func (ac *modelRegistryCall) Invoke(ctx context.Context, toolParameters []byte) ([]*biz_entity.ToolInvokeMessage, error) {

	toolKeyMapInvoke := fmt.Sprintf("%s/%s", ac.toolRuntime.Identity.Provider, ac.toolRuntime.Identity.Name)

	log.Infof("invoke %s", toolKeyMapInvoke)

	toolIns, err := ToolRuntimeRegistry.Acquire(toolKeyMapInvoke)

	if err != nil {
		return nil, err
	}

	msgs, err := toolIns.Invoke(ctx, ac.userID, toolParameters, ac.toolRuntime)

	if err != nil {
		log.Errorf("tool invoke error %#+v", err)
	}

	return msgs, err
}
