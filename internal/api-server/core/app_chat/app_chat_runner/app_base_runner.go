package app_chat_runner

import (
	"context"
	"time"

	"github.com/lunarianss/Luna/internal/api-server/core/app_chat/token_buffer_memory"
	prompt "github.com/lunarianss/Luna/internal/api-server/core/app_prompt"
	biz_entity_app_config "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/provider_app_config"
	po_entity_app "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	po_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	biz_entity_app_generate "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_generate"
	biz_entity_provider_config "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_configuration"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
)

type AppBaseChatRunner struct {
}

func NewAppBaseChatRunner() *AppBaseChatRunner {
	return &AppBaseChatRunner{}
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

func (r *AppBaseChatRunner) DirectOutStream(applicationGenerateEntity biz_entity_app_generate.BasedAppGenerateEntity, message *po_entity_chat.Message, conversation *po_entity_chat.Conversation, queueManager biz_entity.IStreamGenerateQueue, text string, promptMessages []*po_entity_chat.PromptMessage) {

	index := 0
	for i, token := range text {
		tokenStr := string(token)
		llmResultChunk := &biz_entity.LLMResultChunk{
			Model: applicationGenerateEntity.GetModel(),
			PromptMessage: util.ConvertToInterfaceSlice(promptMessages, func(pm *po_entity_chat.PromptMessage) po_entity_chat.IPromptMessage {
				return pm
			}),
			Delta: &biz_entity.LLMResultChunkDelta{
				Index:   index,
				Message: biz_entity.NewAssistantPromptMessage(tokenStr),
			},
		}
		event := biz_entity.NewAppQueueEvent(biz_entity.LLMChunk)
		queueManager.Push(&biz_entity.QueueLLMChunkEvent{
			AppQueueEvent: event,
			Chunk:         llmResultChunk})
		i++
		time.Sleep(10 * time.Millisecond)
	}

	queueManager.Final(&biz_entity.QueueMessageEndEvent{
		AppQueueEvent: biz_entity.NewAppQueueEvent(biz_entity.MessageEnd),
		LLMResult: &biz_entity.LLMResult{
			Model: applicationGenerateEntity.GetModel(),
			PromptMessage: util.ConvertToInterfaceSlice(promptMessages, func(pm *po_entity_chat.PromptMessage) po_entity_chat.IPromptMessage {
				return pm
			}),
			Message: biz_entity.NewAssistantPromptMessage(text),
			Usage:   biz_entity.NewEmptyLLMUsage(),
		},
	})
}
