// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model_registry

import (
	"context"
	"fmt"
	"sync"

	"github.com/lunarianss/Luna/infrastructure/errors"
	biz_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
)

const (
	PROVIDER_NUMBER = 52
)

var (
	ModelRuntimeRegistry = &ModelRegistries{
		ModelRegistry: make(map[string]IModelRegistry, PROVIDER_NUMBER),
		RWMutex:       &sync.RWMutex{},
	}
)

type ModelRegistries struct {
	ModelRegistry map[string]IModelRegistry
	*sync.RWMutex
}

type IModelRegistry interface {
	Invoke(ctx context.Context, queueManager *biz_entity_chat.StreamGenerateQueue, model string, credentials map[string]interface{}, modelParameters map[string]interface{}, stop []string, stream bool, user string, promptMessages []*po_entity.PromptMessage)
	RegisterName() string
}

func (mr *ModelRegistries) RegisterLargeModelInstance(modelRegistry IModelRegistry) {
	defer mr.Unlock()
	mr.Lock()

	mr.ModelRegistry[modelRegistry.RegisterName()] = modelRegistry
}

func (mr *ModelRegistries) Acquire(name string) (IModelRegistry, error) {
	defer mr.RUnlock()
	mr.RLock()
	AIModelIns, ok := mr.ModelRegistry[name]

	if !ok {
		return nil, errors.WithCode(code.ErrNotFoundModelRegistry, fmt.Sprintf("registry %s not found", name))
	}

	return AIModelIns, nil
}
