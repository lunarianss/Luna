package chat

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config/entities"
	appEntities "github.com/lunarianss/Luna/internal/api-server/core/app/apps/entities"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
	"github.com/lunarianss/Luna/internal/api-server/entities/message"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime"
	groqLLM "github.com/lunarianss/Luna/internal/api-server/model_runtime/model_providers/groq/llm"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime/model_providers/openai_api_compatible/llm"
	"github.com/lunarianss/Luna/pkg/log"
)

type ChatAppGenerator struct {
}

func (g *ChatAppGenerator) Generate(c context.Context, appModel *model.App, user interface{}, args *dto.CreateChatMessageBody, invokeFrom appEntities.InvokeForm, stream bool) error {

	// var conversation *model.Conversation
	// var message *model.Message

	// query := args.Query
	// inputs := args.Inputs

	// role := model.AccountCreatedByRole

	StreamResultChunkQueue := make(chan entities.IQueueEvent, model_runtime.STREAM_BUFFER_SIZE)
	StreamFinalChunkQueue := make(chan entities.IQueueEvent, model_runtime.STREAM_BUFFER_SIZE)

	queueManager := model_runtime.NewStreamGenerateQueue(
		StreamResultChunkQueue,
		StreamFinalChunkQueue,
		uuid.NewString(),
		"",
		"",
		"",
		model.AppMode("chat"),
		invokeFrom)

	var credentials = map[string]interface{}{
		"api_key": "xxx",
	}

	go func() {
		openApiCompactModel := &llm.OpenApiCompactLargeLanguageModel{Stream: true, Stop: nil, StreamGenerateQueue: queueManager, Model: "llama3-8b-8192"}

		groqLM := groqLLM.GroqLargeLanguageModel{
			OpenApiCompactLargeLanguageModel: openApiCompactModel,
		}

		msg := &message.PromptMessage{
			Content: "Explain the importance of fast language models",
			Role:    "user",
		}

		msgs := []*message.PromptMessage{msg}
		groqLM.Invoke(c, "llama-3.1-70b-versatile", credentials, nil, nil, true, "", msgs)
	}()

	go func() {
		queueManager.Listen()
	}()

	// 确保 Gin 使用 HTTP 流式传输
	c.(*gin.Context).Writer.Header().Set("Content-Type", "text/event-stream")
	c.(*gin.Context).Writer.Header().Set("Cache-Control", "no-cache")
	c.(*gin.Context).Writer.Header().Set("Connection", "keep-alive")

	// 确保 c.Writer 实现了 http.Flusher 接口
	flusher, ok := c.(*gin.Context).Writer.(http.Flusher)
	if !ok {
		c.(*gin.Context).String(http.StatusInternalServerError, "Streaming unsupported!")
		return nil
	}

	for v := range StreamResultChunkQueue {
		if cm, ok := v.(*entities.QueueLLMChunkEvent); ok {
			log.Info("event %s, answer %+v", cm.Event, cm.Chunk.Delta)
			// 将事件格式化为 SSE 格式发送给客户端
			fmt.Fprintf(c.(*gin.Context).Writer, "data: {answer: %s}\n\n", cm.Chunk.Delta.Message.Content)
			flusher.Flush() // 确保数据立即发送到客户端
		}
	}

	for v := range StreamFinalChunkQueue {
		if mc, ok := v.(*entities.QueueLLMChunkEvent); ok {
			log.Info("event %s, answer %+v", mc.Event, mc.Chunk.Delta)
			// 将事件格式化为 SSE 格式发送给客户端
			fmt.Fprintf(c.(*gin.Context).Writer, "data: %s\n\n", mc.Chunk.Delta.Message.Content)
			flusher.Flush() // 确保数据立即发送到客户端
		} else if mc, ok := v.(*entities.QueueMessageEndEvent); ok {
			log.Info("event %s, end %+v", mc.Event, mc.LLMResult)
		}
	}

	return nil
}
