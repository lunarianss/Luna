package biz_entity

import "context"

type AgentFlusher interface {
	AgentThoughtToStreamResponse(ctx context.Context, agentThoughtID string) error
	AgentMessageToStreamResponse(answer string) error
	AgentMessageFileToStreamResponse(ctx context.Context, messageFileID string, secretKey string, baseUrl string) error
	InitFlusher(ctx context.Context)
}
