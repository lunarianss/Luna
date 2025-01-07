package biz_entity

import "context"

type AgentFlusher interface {
	ManualFlush(streamString string) error
	AgentThoughtToStreamResponse(ctx context.Context, agentThoughtID string) error
	AgentMessageToStreamResponse(answer string) error
}
