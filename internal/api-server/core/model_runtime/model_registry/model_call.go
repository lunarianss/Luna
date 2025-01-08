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

type IModelRegistryCall interface {
	InvokeLLM(ctx context.Context, promptMessage []po_entity.IPromptMessage, queueManager biz_entity_chat.IStreamGenerateQueue, modelParameters map[string]interface{}, tools []*biz_entity_chat.PromptMessageTool, stop []string, user string, callbacks interface{})

	InvokeLLMNonStream(ctx context.Context, promptMessage []po_entity.IPromptMessage, modelParameters map[string]interface{}, tools interface{}, stop []string, user string, callbacks interface{}) (*biz_entity_chat.LLMResult, error)

	InvokeSpeechToText(ctx context.Context, audioFileContent []byte, user string, filename string) (string, error)

	InvokeTextToSpeech(ctx context.Context, modelParameters map[string]interface{}, user string, voice string, format string, texts []string) error

	InvokeTextEmbedding(ctx context.Context, modelParameters map[string]interface{}, user string, inputType string, texts []string) (*biz_entity_chat.TextEmbeddingResult, error)
}

type modelRegistryCall struct {
	Model        string
	Provider     string
	Credentials  map[string]interface{}
	ModelType    string
	ModelRuntime biz_entity.IAIModelRuntime
}

func NewModelRegisterCaller(model, modelType, provider string, credentials map[string]interface{}, modelRuntime biz_entity.IAIModelRuntime) IModelRegistryCall {
	return &modelRegistryCall{
		Model:        model,
		ModelType:    modelType,
		Provider:     provider,
		Credentials:  credentials,
		ModelRuntime: modelRuntime,
	}
}

func (ac *modelRegistryCall) InvokeLLM(ctx context.Context, promptMessage []po_entity.IPromptMessage, queueManager biz_entity_chat.IStreamGenerateQueue, modelParameters map[string]interface{}, tools []*biz_entity_chat.PromptMessageTool, stop []string, user string, callbacks interface{}) {

	modelKeyMapInvoke := fmt.Sprintf("%s/%s", ac.Provider, ac.ModelType)

	log.Infof("invoke %s", modelKeyMapInvoke)

	AIModelIns, err := ModelRuntimeRegistry.Acquire(modelKeyMapInvoke)

	if err != nil {
		queueManager.PushErr(err)
		return
	}

	AIModelIns.Invoke(ctx, queueManager, ac.Model, ac.Credentials, modelParameters, stop, user, promptMessage, ac.ModelRuntime, tools)
}

func (ac *modelRegistryCall) InvokeLLMNonStream(ctx context.Context, promptMessage []po_entity.IPromptMessage, modelParameters map[string]interface{}, tools interface{}, stop []string, user string, callbacks interface{}) (*biz_entity_chat.LLMResult, error) {

	modelKeyMapInvoke := fmt.Sprintf("%s/%s", ac.Provider, ac.ModelType)

	log.Infof("invoke %s", modelKeyMapInvoke)

	AIModelIns, err := ModelRuntimeRegistry.Acquire(modelKeyMapInvoke)

	if err != nil {
		return nil, err
	}

	return AIModelIns.InvokeNonStream(ctx, ac.Model, ac.Credentials, modelParameters, stop, user, promptMessage, ac.ModelRuntime)
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

func (ac *modelRegistryCall) InvokeTextEmbedding(ctx context.Context, modelParameters map[string]interface{}, user string, inputType string, texts []string) (*biz_entity_chat.TextEmbeddingResult, error) {

	modelKeyMapInvoke := fmt.Sprintf("%s/%s", ac.Provider, ac.ModelType)

	log.Infof("invoke %s", modelKeyMapInvoke)

	AIModelIns, err := TextEmbeddingRegistry.Acquire(modelKeyMapInvoke)

	if err != nil {
		return nil, err
	}
	return AIModelIns.Embedding(ctx, ac.Model, ac.Credentials, modelParameters, user, ac.ModelRuntime, inputType, texts)
}
