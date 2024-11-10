package llm

import "github.com/lunarianss/Luna/internal/api-server/entities/message"

type LLMUsage struct {
	PromptTokens        int     `json:"prompt_tokens"`
	PromptUnitPrice     float64 `json:"prompt_unit_price"`
	PromptPriceUnit     float64 `json:"prompt_price_unit"`
	PromptPrice         float64 `json:"prompt_price"`
	CompletionTokens    int     `json:"completion_tokens"`
	CompletionUnitPrice float64 `json:"completion_unit_price"`
	CompletionPriceUnit float64 `json:"completion_price_unit"`
	CompletionPrice     float64 `json:"completion_price"`
	TotalTokens         int     `json:"total_tokens"`
	TotalPrice          float64 `json:"total_price"`
	Currency            string  `json:"currency"`
	Latency             float64 `json:"latency"`
}

type LLMResultChunkDelta struct {
	Index        int                             `json:"index"`
	Message      *message.AssistantPromptMessage `json:"message"`
	Usage        *LLMUsage                       `json:"usage"`
	FinishReason string                          `json:"finish_reason"`
}

type LLMResultChunk struct {
	ID                string                   `json:"id"`
	Model             string                   `json:"model"`
	PromptMessage     []*message.PromptMessage `json:"prompt_message"`
	SystemFingerprint string                   `json:"system_fingerprint"`
	Delta             *LLMResultChunkDelta     `json:"delta"`
}

type LLMResult struct {
	ID                string                          `json:"id"`
	Model             string                          `json:"model"`
	Message           *message.AssistantPromptMessage `json:"message"`
	PromptMessage     []*message.PromptMessage        `json:"prompt_message"`
	Usage             *LLMUsage                       `json:"usage"`
	SystemFingerprint string                          `json:"system_fingerprint"`
	Reason            string                          `json:"reason"`
}
