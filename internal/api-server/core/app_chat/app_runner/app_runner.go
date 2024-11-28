// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package app_runner

import (
	"context"

	prompt "github.com/lunarianss/Luna/internal/api-server/core/app_prompt"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	po_entity_app "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	po_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	biz_entity_app_config "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_config"
	biz_entity_app_generate "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_generate"
	biz_entity_provider_config "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_configuration"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime/model_registry"
)

type AppRunner struct {
	AppDomain *domain_service.AppDomain
}

func (runner *AppRunner) HandleInvokeResultStream(ctx context.Context, invokeResult *biz_entity.LLMResultChunk, streamGenerator *biz_entity.StreamGenerateQueue, end bool, err error) {

	if err != nil && invokeResult == nil {
		streamGenerator.Final(&biz_entity.QueueErrorEvent{
			AppQueueEvent: biz_entity.NewAppQueueEvent(biz_entity.Error),
			Err:           err,
		})
		return
	}

	if end {
		llmResult := &biz_entity.LLMResult{
			Model:         invokeResult.Model,
			PromptMessage: invokeResult.PromptMessage,
			Reason:        invokeResult.Delta.FinishReason,
			Message: &biz_entity.AssistantPromptMessage{
				PromptMessage: &po_entity_chat.PromptMessage{
					Content: invokeResult.Delta.Message.Content,
				},
			},
		}

		event := biz_entity.NewAppQueueEvent(biz_entity.MessageEnd)
		streamGenerator.Final(&biz_entity.QueueMessageEndEvent{
			AppQueueEvent: event,
			LLMResult:     llmResult,
		})
		return
	}

	event := biz_entity.NewAppQueueEvent(biz_entity.LLMChunk)
	streamGenerator.Push(&biz_entity.QueueLLMChunkEvent{
		AppQueueEvent: event,
		Chunk:         invokeResult})

}

func (r *AppRunner) Run(ctx context.Context, applicationGenerateEntity *biz_entity_app_generate.ChatAppGenerateEntity, message *po_entity_chat.Message, conversation *po_entity_chat.Conversation, queueManager *biz_entity.StreamGenerateQueue) {

	appRecord, err := r.AppDomain.AppRepo.GetAppByID(ctx, applicationGenerateEntity.AppConfig.AppID)

	if err != nil {
		queueManager.PushErr(err)
		return
	}

	credentials, err := applicationGenerateEntity.ModelConf.ProviderModelBundle.Configuration.GetCurrentCredentials(applicationGenerateEntity.ModelConf.ProviderModelBundle.ModelTypeInstance.ModelType, applicationGenerateEntity.AppConfig.Model.Model)

	if err != nil {
		queueManager.PushErr(err)
		return
	}

	promptMessages, stop, err := r.OrganizePromptMessage(ctx, appRecord, applicationGenerateEntity.ModelConf, applicationGenerateEntity.AppConfig.PromptTemplate, applicationGenerateEntity.Inputs, nil, applicationGenerateEntity.Query, "", nil)

	if err != nil {
		queueManager.PushErr(err)
		return
	}

	modelInstance := model_registry.ModelInstance{
		Model:               applicationGenerateEntity.AppConfig.Model.Model,
		ProviderModelBundle: applicationGenerateEntity.ModelConf.ProviderModelBundle,
		ModelTypeInstance:   applicationGenerateEntity.ModelConf.ProviderModelBundle.ModelTypeInstance,
		Credentials:         credentials,
		Provider:            applicationGenerateEntity.ModelConf.ProviderModelBundle.Configuration.Provider.Provider,
	}

	modelInstance.InvokeLLM(ctx, promptMessages, queueManager, applicationGenerateEntity.ModelConf.Parameters, nil, stop, applicationGenerateEntity.Stream, applicationGenerateEntity.UserID, nil)
}

func (r *AppRunner) OrganizePromptMessage(ctx context.Context, appRecord *po_entity_app.App, modelConfig *biz_entity_provider_config.ModelConfigWithCredentialsEntity, promptTemplateEntity *biz_entity_app_config.PromptTemplateEntity, inputs map[string]interface{}, files []string, query string, context string, memory any) ([]*po_entity_chat.PromptMessage, []string, error) {

	var (
		promptMessages []*po_entity_chat.PromptMessage
		stop           []string
		err            error
	)
	if promptTemplateEntity.PromptType == string(biz_entity_app_config.SIMPLE) {
		simplePrompt := prompt.SimplePromptTransform{}

		promptMessages, stop, err = simplePrompt.GetPrompt(po_entity_app.AppMode(appRecord.Mode), promptTemplateEntity, inputs, query, files, context, nil, modelConfig)

		if err != nil {
			return nil, nil, err
		}
	}

	return promptMessages, stop, err
}
