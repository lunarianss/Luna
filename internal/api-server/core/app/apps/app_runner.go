package apps

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/core/app"
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config/entities"
	"github.com/lunarianss/Luna/internal/api-server/entities/llm"
	"github.com/lunarianss/Luna/internal/api-server/entities/message"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime/model_registry"
)

type AppRunner struct{}

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

	credentials, err := applicationGenerateEntity.ModelConf.ProviderModelBundle.Configuration.GetCurrentCredentials(applicationGenerateEntity.ModelConf.ProviderModelBundle.ModelTypeInstance.ModelType, applicationGenerateEntity.AppConfig.Model.Model)

	if err != nil {
		queueManager.PushErr(err)
	}

	modelInstance := model_registry.ModelInstance{
		Model:               applicationGenerateEntity.AppConfig.Model.Model,
		ProviderModelBundle: applicationGenerateEntity.ModelConf.ProviderModelBundle,
		ModelTypeInstance:   applicationGenerateEntity.ModelConf.ProviderModelBundle.ModelTypeInstance,
		Credentials:         credentials,
		Provider:            applicationGenerateEntity.ModelConf.ProviderModelBundle.Configuration.Provider.Provider,
	}

	modelInstance.InvokeLLM(ctx, nil, queueManager, applicationGenerateEntity.ModelConf.Parameters, nil, nil, applicationGenerateEntity.Stream, applicationGenerateEntity.UserID, nil)
}
