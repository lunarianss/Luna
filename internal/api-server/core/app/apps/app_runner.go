package apps

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config/entities"
	"github.com/lunarianss/Luna/internal/api-server/entities/llm"
	"github.com/lunarianss/Luna/internal/api-server/entities/message"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime"
)

type AppRunner struct{}

func (runner *AppRunner) HandleInvokeResultStream(ctx context.Context, invokeResult *llm.LLMResultChunk, streamGenerator *model_runtime.StreamGenerateQueue, end bool) {

	var (
		model         string
		promptMessage []*message.PromptMessage
		text          string
	)

	event := entities.NewAppQueueEvent(entities.LLMChunk)
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
