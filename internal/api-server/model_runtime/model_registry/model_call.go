// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model_registry

import (
	"context"
	"fmt"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider/model_provider"
	"github.com/lunarianss/Luna/infrastructure/log"
	biz_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
)

type modelRegistryCall struct {
	Model       string
	Provider    string
	Credentials map[string]interface{}
	ModelType   string
	ModelRuntime biz_entity.IAIModelRuntime
}

func NewModelRegisterCaller(model, modelType, provider string, credentials map[string]interface{}, modelRuntime biz_entity.IAIModelRuntime) *modelRegistryCall {
	return &modelRegistryCall{
		Model:       model,
		ModelType:   modelType,
		Provider:    provider,
		Credentials: credentials,
		ModelRuntime: modelRuntime,
	}

}
func (ac *modelRegistryCall) InvokeLLM(ctx context.Context, promptMessage []*po_entity.PromptMessage, queueManager *biz_entity_chat.StreamGenerateQueue, modelParameters map[string]interface{}, tools interface{}, stop []string, stream bool, user string, callbacks interface{}) {

	modelKeyMapInvoke := fmt.Sprintf("%s/%s", ac.Provider, ac.ModelType)

	log.Infof("invoke %s", modelKeyMapInvoke)

	AIModelIns, err := ModelRuntimeRegistry.Acquire(modelKeyMapInvoke)

	if err != nil {
		queueManager.PushErr(err)
		return
	}
	AIModelIns.Invoke(ctx, queueManager, ac.Model, ac.Credentials, modelParameters, stop, stream, user, promptMessage, ac.ModelRuntime)
}
