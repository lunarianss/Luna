package apps

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/core/app"
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config"
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config/entities"
	"github.com/lunarianss/Luna/internal/api-server/core/prompt"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	po_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/entities/llm"
	"github.com/lunarianss/Luna/internal/api-server/entities/message"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime/model_registry"
)

type AppRunner struct {
	AppDomain *domain_service.AppDomain
}

func (runner *AppRunner) HandleInvokeResultStream(ctx context.Context, invokeResult *llm.LLMResultChunk, streamGenerator *model_runtime.StreamGenerateQueue, end bool, err error) {

	if err != nil && invokeResult == nil {
		streamGenerator.Final(&entities.QueueErrorEvent{
			AppQueueEvent: entities.NewAppQueueEvent(entities.Error),
			Err:           err,
		})
		return
	}

	if end {
		llmResult := &llm.LLMResult{
			Model:         invokeResult.Model,
			PromptMessage: invokeResult.PromptMessage,
			Reason:        invokeResult.Delta.FinishReason,
			Message: &message.AssistantPromptMessage{
				PromptMessage: &message.PromptMessage{
					Content: invokeResult.Delta.Message.Content,
				},
			},
		}

		event := entities.NewAppQueueEvent(entities.MessageEnd)
		streamGenerator.Final(&entities.QueueMessageEndEvent{
			AppQueueEvent: event,
			LLMResult:     llmResult,
		})
		return
	}

	event := entities.NewAppQueueEvent(entities.LLMChunk)
	streamGenerator.Push(&entities.QueueLLMChunkEvent{
		AppQueueEvent: event,
		Chunk:         invokeResult})

}

func (r *AppRunner) Run(ctx context.Context, applicationGenerateEntity *app.ChatAppGenerateEntity, message *po_entity_chat.Message, conversation *po_entity_chat.Conversation, queueManager *model_runtime.StreamGenerateQueue) {

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

func (r *AppRunner) OrganizePromptMessage(ctx context.Context, appRecord *po_entity.App, modelConfig *app.ModelConfigWithCredentialsEntity, promptTemplateEntity *app_config.PromptTemplateEntity, inputs map[string]interface{}, files []string, query string, context string, memory any) ([]*message.PromptMessage, []string, error) {

	var (
		promptMessages []*message.PromptMessage
		stop           []string
		err            error
	)
	if promptTemplateEntity.PromptType == string(app_config.SIMPLE) {
		simplePrompt := prompt.SimplePromptTransform{}

		promptMessages, stop, err = simplePrompt.GetPrompt(po_entity.AppMode(appRecord.Mode), promptTemplateEntity, inputs, query, files, context, nil, modelConfig)

		if err != nil {
			return nil, nil, err
		}
	}

	return promptMessages, stop, err
}
