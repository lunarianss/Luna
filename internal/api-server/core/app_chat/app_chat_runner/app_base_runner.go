package app_chat_runner

import (
	"context"
	"time"

	"github.com/lunarianss/Luna/internal/api-server/core/app_chat/token_buffer_memory"
	prompt "github.com/lunarianss/Luna/internal/api-server/core/app_prompt"
	biz_entity_app_config "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/provider_app_config"
	po_entity_app "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	po_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	biz_entity_app_generate "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_generate"
	biz_entity_provider_config "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_configuration"
	"github.com/lunarianss/Luna/internal/infrastructure/util"

	biz_entity_chat_prompt_message "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity/chat_prompt_message"
	biz_entity_base_stream_generator "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity/stream_base_generator"
)

type AppBaseChatRunner struct {
}

func NewAppBaseChatRunner() *AppBaseChatRunner {
	return &AppBaseChatRunner{}
}

func (r *AppBaseChatRunner) OrganizePromptMessage(ctx context.Context, appRecord *po_entity_app.App, modelConfig *biz_entity_provider_config.ModelConfigWithCredentialsEntity, promptTemplateEntity *biz_entity_app_config.PromptTemplateEntity, inputs map[string]interface{}, files []string, query string, context string, memory token_buffer_memory.ITokenBufferMemory) ([]*biz_entity_chat_prompt_message.PromptMessage, []string, error) {

	var (
		promptMessages []*biz_entity_chat_prompt_message.PromptMessage
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

func (r *AppBaseChatRunner) DirectOutStream(applicationGenerateEntity biz_entity_app_generate.BasedAppGenerateEntity, message *po_entity_chat.Message, conversation *po_entity_chat.Conversation, queueManager biz_entity_base_stream_generator.IStreamGenerateQueue, text string, promptMessages []*biz_entity_chat_prompt_message.PromptMessage) {

	index := 0
	for i, token := range text {
		tokenStr := string(token)
		llmResultChunk := &biz_entity_base_stream_generator.LLMResultChunk{
			Model: applicationGenerateEntity.GetModel(),
			PromptMessage: util.ConvertToInterfaceSlice(promptMessages, func(pm *biz_entity_chat_prompt_message.PromptMessage) biz_entity_chat_prompt_message.IPromptMessage {
				return pm
			}),
			Delta: &biz_entity_base_stream_generator.LLMResultChunkDelta{
				Index:   index,
				Message: biz_entity_chat_prompt_message.NewAssistantToolPromptMessage(tokenStr),
			},
		}
		event := biz_entity_base_stream_generator.NewAppQueueEvent(biz_entity_base_stream_generator.LLMChunk)
		queueManager.Push(&biz_entity_base_stream_generator.QueueLLMChunkEvent{
			AppQueueEvent: event,
			Chunk:         llmResultChunk})
		i++
		time.Sleep(10 * time.Millisecond)
	}

	queueManager.Final(&biz_entity_base_stream_generator.QueueMessageEndEvent{
		AppQueueEvent: biz_entity_base_stream_generator.NewAppQueueEvent(biz_entity_base_stream_generator.MessageEnd),
		LLMResult: &biz_entity_base_stream_generator.LLMResult{
			Model: applicationGenerateEntity.GetModel(),
			PromptMessage: util.ConvertToInterfaceSlice(promptMessages, func(pm *biz_entity_chat_prompt_message.PromptMessage) biz_entity_chat_prompt_message.IPromptMessage {
				return pm
			}),
			Message: biz_entity_chat_prompt_message.NewAssistantToolPromptMessage(text),
			Usage:   biz_entity_base_stream_generator.NewEmptyLLMUsage(),
		},
	})
}
