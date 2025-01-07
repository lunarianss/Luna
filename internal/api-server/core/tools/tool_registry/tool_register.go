// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package tool_registry

import (
	"context"
	"fmt"
	"sync"

	"github.com/lunarianss/Luna/infrastructure/errors"

	"github.com/lunarianss/Luna/internal/api-server/domain/agent/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
)

const (
	TOOL_NUMBER = 100
)

type IToolCallRegistry interface {
	Invoke(ctx context.Context, userID string, toolParameters []byte, toolRuntime *biz_entity.ToolRuntimeConfiguration) ([]*biz_entity.ToolInvokeMessage, error)
	Register() string
}

var (
	ToolRuntimeRegistry = &ModelRegistries[IToolCallRegistry]{
		toolRegistry: make(map[string]IToolCallRegistry, TOOL_NUMBER),
		RWMutex:      &sync.RWMutex{},
	}
)

type ModelRegistries[T any] struct {
	toolRegistry map[string]T
	*sync.RWMutex
}

func (mr *ModelRegistries[T]) RegisterAgentToolInstance(toolRegistry T) {
	defer mr.Unlock()
	mr.Lock()

	switch v := any(toolRegistry).(type) {
	case IToolCallRegistry:
		mr.toolRegistry[v.Register()] = toolRegistry
	default:
		panic(fmt.Sprintf("AI multiply tool registry error: %+v is not implement IToolCallRegistry", biz_entity.ToolLabelEducation))
	}
}

func (mr *ModelRegistries[T]) Acquire(name string) (T, error) {
	defer mr.RUnlock()
	mr.RLock()

	var zeroT T

	toolIns, ok := mr.toolRegistry[name]

	if !ok {
		return zeroT, errors.WithCode(code.ErrNotFoundToolRegistry, "registry %s not found", name)
	}

	return toolIns, nil
}
