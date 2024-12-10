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
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider/model_provider"
)

type modelRegistryCall struct {
	Model        string
	Provider     string
	Credentials  map[string]interface{}
	ModelType    string
	ModelRuntime biz_entity.IAIModelRuntime
}

func NewModelRegisterCaller(model, modelType, provider string, credentials map[string]interface{}, modelRuntime biz_entity.IAIModelRuntime) *modelRegistryCall {
	return &modelRegistryCall{
		Model:        model,
		ModelType:    modelType,
		Provider:     provider,
		Credentials:  credentials,
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

func (ac *modelRegistryCall) InvokeSpeechToText(ctx context.Context, audioFileContent []byte, user string, filename string) (string, error) {

	modelKeyMapInvoke := fmt.Sprintf("%s/%s", ac.Provider, ac.ModelType)

	log.Infof("invoke %s", modelKeyMapInvoke)

	AIModelIns, err := AudioModelRuntimeRegistry.Acquire(modelKeyMapInvoke)

	if err != nil {
		return "", err
	}

	resp, err := AIModelIns.Invoke(ctx, ac.Model, ac.Credentials, nil, user, filename, audioFileContent, ac.ModelRuntime)

	if err != nil {
		return "", err
	}

	return resp.Text, nil
}

func (ac *modelRegistryCall) InvokeTextToSpeech(ctx context.Context, modelParameters map[string]interface{}, user string, voice string, format string, texts []string) error {

	modelKeyMapInvoke := fmt.Sprintf("%s/%s", ac.Provider, ac.ModelType)

	log.Infof("invoke %s", modelKeyMapInvoke)

	AIModelIns, err := TTSModelRuntimeRegistry.Acquire(modelKeyMapInvoke)

	if err != nil {
		return err
	}

	err = AIModelIns.Invoke(ctx, ac.Model, ac.Credentials, nil, user, "", voice, ac.ModelRuntime, format, texts)

	if err != nil {
		return err
	}
	return nil
}
