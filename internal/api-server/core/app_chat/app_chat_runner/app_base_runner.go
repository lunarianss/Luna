package app_chat_runner

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/core/app_chat/token_buffer_memory"
	prompt "github.com/lunarianss/Luna/internal/api-server/core/app_prompt"
	biz_entity_app_config "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/provider_app_config"
	po_entity_app "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	po_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	biz_entity_provider_config "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_configuration"
)

type AppBaseChatRunner struct {
}

func NewAppBaseChatRunner() *AppBaseChatRunner {
	return &AppBaseChatRunner{}
}

func (runner *AppBaseChatRunner) HandleInvokeResultStream(ctx context.Context, invokeResult *biz_entity.LLMResultChunk, streamGenerator *biz_entity.StreamGenerateQueue, end bool, err error) {

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
			Usage: invokeResult.Delta.Usage,
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

func (r *AppBaseChatRunner) OrganizePromptMessage(ctx context.Context, appRecord *po_entity_app.App, modelConfig *biz_entity_provider_config.ModelConfigWithCredentialsEntity, promptTemplateEntity *biz_entity_app_config.PromptTemplateEntity, inputs map[string]interface{}, files []string, query string, context string, memory token_buffer_memory.ITokenBufferMemory) ([]*po_entity_chat.PromptMessage, []string, error) {

	var (
		promptMessages []*po_entity_chat.PromptMessage
		stop           []string
		err            error
	)
	if promptTemplateEntity.PromptType == string(biz_entity_app_config.SIMPLE) {
		simplePrompt := prompt.SimplePromptTransform{}

		promptMessages, stop, err = simplePrompt.GetPrompt(po_entity_app.AppMode(appRecord.Mode), promptTemplateEntity, inputs, query, files, context, memory, modelConfig)

		if err != nil {
			return nil, nil, err
		}
	}

	return promptMessages, stop, err
}
