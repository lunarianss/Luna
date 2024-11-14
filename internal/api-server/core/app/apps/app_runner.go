package apps

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/core/app"
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config"
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config/entities"
	"github.com/lunarianss/Luna/internal/api-server/core/prompt"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/app"
	"github.com/lunarianss/Luna/internal/api-server/entities/llm"
	"github.com/lunarianss/Luna/internal/api-server/entities/message"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime/model_registry"
)

type AppRunner struct {
	AppDomain *domain.AppDomain
}

func (runner *AppRunner) HandleInvokeResultStream(ctx context.Context, invokeResult *llm.LLMResultChunk, streamGenerator *model_runtime.StreamGenerateQueue, end bool, err error) {

	if err != nil && invokeResult == nil {
		streamGenerator.Final(&entities.QueueErrorEvent{
			AppQueueEvent: entities.NewAppQueueEvent(entities.Error),
			Err:           err,
		})
		return
	}

	var (
		model         string
		promptMessage []*message.PromptMessage
		text          string
		event         *entities.AppQueueEvent
	)

	event = entities.NewAppQueueEvent(entities.LLMChunk)
	streamGenerator.Push(&entities.QueueLLMChunkEvent{
		AppQueueEvent: event,
		Chunk:         invokeResult})

	if contentStr, ok := invokeResult.Delta.Message.Content.(string); ok {
		text += contentStr
	}

	if model == "" {
		model = invokeResult.Model
	}

	promptMessage = invokeResult.PromptMessage

	if end {
		llmResult := &llm.LLMResult{
			Model:         model,
			PromptMessage: promptMessage,
			Reason:        invokeResult.Delta.FinishReason,
			Message: &message.AssistantPromptMessage{
				PromptMessage: &message.PromptMessage{
					Content: text,
				},
			},
		}

		event := entities.NewAppQueueEvent(entities.MessageEnd)

		streamGenerator.Final(&entities.QueueMessageEndEvent{
			AppQueueEvent: event,
			LLMResult:     llmResult,
		})
	}

}

func (r *AppRunner) Run(ctx context.Context, applicationGenerateEntity *app.ChatAppGenerateEntity, message *model.Message, conversation *model.Conversation, queueManager *model_runtime.StreamGenerateQueue) {

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

func (r *AppRunner) OrganizePromptMessage(ctx context.Context, appRecord *model.App, modelConfig *app.ModelConfigWithCredentialsEntity, promptTemplateEntity *app_config.PromptTemplateEntity, inputs map[string]interface{}, files []string, query string, context string, memory any) ([]*message.PromptMessage, []string, error) {

	var (
		promptMessages []*message.PromptMessage
		stop           []string
		err            error
	)
	if promptTemplateEntity.PromptType == string(app_config.SIMPLE) {
		simplePrompt := prompt.SimplePromptTransform{}

		promptMessages, stop, err = simplePrompt.GetPrompt(model.AppMode(appRecord.Mode), promptTemplateEntity, inputs, query, files, context, nil, modelConfig)

		if err != nil {
			return nil, nil, err
		}
	}

	return promptMessages, stop, err
}
