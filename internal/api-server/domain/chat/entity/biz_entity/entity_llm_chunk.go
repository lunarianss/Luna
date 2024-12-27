package biz_entity

import "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"

// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

type LLMUsage struct {
	PromptTokens        int64   `json:"prompt_tokens"`
	PromptUnitPrice     float64 `json:"prompt_unit_price"`
	PromptPriceUnit     float64 `json:"prompt_price_unit"`
	PromptPrice         float64 `json:"prompt_price"`
	CompletionTokens    int64   `json:"completion_tokens"`
	CompletionUnitPrice float64 `json:"completion_unit_price"`
	CompletionPriceUnit float64 `json:"completion_price_unit"`
	CompletionPrice     float64 `json:"completion_price"`
	TotalTokens         int64   `json:"total_tokens"`
	TotalPrice          float64 `json:"total_price"`
	Currency            string  `json:"currency"`
	Latency             float64 `json:"latency"`
}

func NewEmptyLLMUsage() *LLMUsage {
	return &LLMUsage{
		PromptTokens:        0,
		PromptUnitPrice:     0,
		PromptPriceUnit:     0,
		PromptPrice:         0,
		CompletionTokens:    0,
		CompletionUnitPrice: 0,
		CompletionPriceUnit: 0,
		CompletionPrice:     0,
		TotalTokens:         0,
		TotalPrice:          0,
		Currency:            "",
		Latency:             0,
	}
}

type LLMResultChunkDelta struct {
	Index        int                     `json:"index"`
	Message      *AssistantPromptMessage `json:"message"`
	Usage        *LLMUsage               `json:"usage"`
	FinishReason string                  `json:"finish_reason"`
}

type LLMResultChunk struct {
	ID                string                     `json:"id"`
	Model             string                     `json:"model"`
	PromptMessage     []*po_entity.PromptMessage `json:"prompt_message"`
	SystemFingerprint string                     `json:"system_fingerprint"`
	Delta             *LLMResultChunkDelta       `json:"delta"`
}

type LLMResult struct {
	ID                string                     `json:"id"`
	Model             string                     `json:"model"`
	Message           *AssistantPromptMessage    `json:"message"`
	PromptMessage     []*po_entity.PromptMessage `json:"prompt_message"`
	Usage             *LLMUsage                  `json:"usage"`
	SystemFingerprint string                     `json:"system_fingerprint"`
	Reason            string                     `json:"reason"`
}

func NewEmptyLLMResult() *LLMResult {
	return &LLMResult{
		Message:       NewEmptyAssistantPromptMessage(),
		PromptMessage: make([]*po_entity.PromptMessage, 0),
	}
}
