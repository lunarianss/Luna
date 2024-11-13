package chat

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/core/app"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/app"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime/model_providers/openai_api_compatible/llm"
)

type AppRunner struct {
	appDomain *domain.AppDomain
}

func (r *AppRunner) Run(ctx context.Context, applicationGenerateEntity *app.ChatAppGenerateEntity, message *model.Message, conversation *model.Conversation) {

	openApiCompactModel := &llm.OpenApiCompactLargeLanguageModel{Stream: true, Stop: nil, StreamGenerateQueue: queueManager, Model: "llama3-8b-8192"}

	groqLM := groqLLM.GroqLargeLanguageModel{
		OpenApiCompactLargeLanguageModel: openApiCompactModel,
	}

	msg := &message.PromptMessage{
		Content: "用中文阐述一下 AI 的发展史",
		Role:    "user",
	}

	msgs := []*message.PromptMessage{msg}
	groqLM.Invoke(c, "llama-3.1-70b-versatile", credentials, nil, nil, true, "", msgs)

}
