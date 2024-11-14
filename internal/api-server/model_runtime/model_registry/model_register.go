package model_registry

import (
	"context"
	"fmt"
	"sync"

	"github.com/lunarianss/Luna/internal/api-server/entities/message"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
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
	Invoke(ctx context.Context, queueManager *model_runtime.StreamGenerateQueue, model string, credentials map[string]interface{}, modelParameters map[string]interface{}, stop []string, stream bool, user string, promptMessages []*message.PromptMessage)
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
