// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model_registry

import (
	"context"
	"fmt"

	"github.com/lunarianss/Luna/infrastructure/log"
	biz_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	biz_entity_model "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider/model_provider"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_configuration"
)

type ModelInstance struct {
	ProviderModelBundle *biz_entity.ProviderModelBundleRuntime
	Model               string
	Provider            string
	Credentials         map[string]interface{}
	ModelTypeInstance   *biz_entity_model.AIModelRuntime
}

func (ac *ModelInstance) InvokeLLM(ctx context.Context, promptMessage []*po_entity.PromptMessage, queueManager *biz_entity_chat.StreamGenerateQueue, modelParameters map[string]interface{}, tools interface{}, stop []string, stream bool, user string, callbacks interface{}) {

	modelKeyMapInvoke := fmt.Sprintf("%s/%s", ac.Provider, ac.ProviderModelBundle.ModelTypeInstance.ModelType)

	log.Infof("invoke %s", modelKeyMapInvoke)

	AIModelIns, err := ModelRuntimeRegistry.Acquire(modelKeyMapInvoke)

	if err != nil {
		queueManager.PushErr(err)
	}

	AIModelIns.Invoke(ctx, queueManager, ac.Model, ac.Credentials, modelParameters, stop, stream, user, promptMessage)
}
