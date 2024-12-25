// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model_registry

import (
	"context"
	"sync"

	"github.com/lunarianss/Luna/infrastructure/errors"
	biz_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider/model_provider"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
)

const (
	PROVIDER_NUMBER = 52
)

type IModelRegistry interface {
	Invoke(ctx context.Context, queueManager *biz_entity_chat.StreamGenerateQueue, model string, credentials map[string]interface{}, modelParameters map[string]interface{}, stop []string, user string, promptMessages []*po_entity.PromptMessage, modelRuntime biz_entity.IAIModelRuntime)
	InvokeNonStream(ctx context.Context, model string, credentials map[string]interface{}, modelParameters map[string]interface{}, stop []string, user string, promptMessages []*po_entity.PromptMessage, modelRuntime biz_entity.IAIModelRuntime) (*biz_entity_chat.LLMResult, error)

	RegisterName() string
}

type IAudioModelRegistry interface {
	Invoke(ctx context.Context, model string, credentials map[string]interface{}, modelParameters map[string]interface{}, user, filename string, fileContent []byte, modelRuntime biz_entity.IAIModelRuntime) (*biz_entity_chat.Speech2TextResp, error)
	RegisterName() string
}

type ITTSModelRegistry interface {
	Invoke(ctx context.Context, model string, credentials map[string]interface{}, modelParameters map[string]interface{}, user, tenantID string, voice string, modelRuntime biz_entity.IAIModelRuntime, format string, texts []string) error
	RegisterName() string
}

type ITextEmbeddingRegistry interface {
	Embedding(ctx context.Context, model string, credentials map[string]interface{}, modelParameters map[string]interface{}, user string, modelRuntime biz_entity.IAIModelRuntime, inputType string, texts []string) (*biz_entity_chat.TextEmbeddingResult, error)
	RegisterName() string
}

var (
	ModelRuntimeRegistry = &ModelRegistries[IModelRegistry]{
		ModelRegistry: make(map[string]IModelRegistry, PROVIDER_NUMBER),
		RWMutex:       &sync.RWMutex{},
	}

	AudioModelRuntimeRegistry = &ModelRegistries[IAudioModelRegistry]{
		ModelRegistry: make(map[string]IAudioModelRegistry, 12),
		RWMutex:       &sync.RWMutex{},
	}

	TTSModelRuntimeRegistry = &ModelRegistries[ITTSModelRegistry]{
		ModelRegistry: make(map[string]ITTSModelRegistry, 5),
		RWMutex:       &sync.RWMutex{},
	}

	TextEmbeddingRegistry = &ModelRegistries[ITextEmbeddingRegistry]{
		ModelRegistry: make(map[string]ITextEmbeddingRegistry, 5),
		RWMutex:       &sync.RWMutex{},
	}
)

type ModelRegistries[T any] struct {
	ModelRegistry map[string]T
	*sync.RWMutex
}

func (mr *ModelRegistries[T]) RegisterLargeModelInstance(modelRegistry T) {
	defer mr.Unlock()
	mr.Lock()

	switch v := any(modelRegistry).(type) {
	case IModelRegistry:
		mr.ModelRegistry[v.RegisterName()] = modelRegistry
	case IAudioModelRegistry:
		mr.ModelRegistry[v.RegisterName()] = modelRegistry
	case ITTSModelRegistry:
		mr.ModelRegistry[v.RegisterName()] = modelRegistry
	case ITextEmbeddingRegistry:
		mr.ModelRegistry[v.RegisterName()] = modelRegistry
	default:
		panic("AI mulit model registry error: ")
	}
}

func (mr *ModelRegistries[T]) Acquire(name string) (T, error) {
	defer mr.RUnlock()
	mr.RLock()

	var zeroT T

	AIModelIns, ok := mr.ModelRegistry[name]

	if !ok {
		return zeroT, errors.WithCode(code.ErrNotFoundModelRegistry, "registry %s not found", name)
	}

	return AIModelIns, nil
}
