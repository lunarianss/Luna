package task_pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"

	"github.com/lunarianss/Luna/infrastructure/log"
	repo_agent "github.com/lunarianss/Luna/internal/api-server/domain/agent/repository"
	biz_entity_app_generate "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_generate"
)

type agentChatFlusher struct {
	biz_entity_app_generate.BasedAppGenerateEntity

	message   *po_entity.Message
	agentRepo repo_agent.AgentRepo
	flusher   http.Flusher
	sender    io.Writer
}

func NewAgentChatFlusher(
	applicationGenerateEntity biz_entity_app_generate.BasedAppGenerateEntity,
	agentRepo repo_agent.AgentRepo, message *po_entity.Message) *agentChatFlusher {
	return &agentChatFlusher{
		BasedAppGenerateEntity: applicationGenerateEntity,
		agentRepo:              agentRepo,
		message:                message,
	}
}

func (tpp *agentChatFlusher) InitFlusher(ctx context.Context) {
	if !tpp.setFlush(ctx) {
		return
	}
}

func (tpp *agentChatFlusher) flush(streamString string) error {
	if _, err := fmt.Fprintf(tpp.sender, "data: %s\n\n", streamString); err != nil {
		return err
	}
	tpp.flusher.Flush()
	return nil
}

func (tpp *agentChatFlusher) setFlush(c context.Context) bool {

	ginContext, ok := c.(*gin.Context)

	if !ok {
		log.Infof("context is not a gin context")
		return false
	}

	ginContext.Writer.Header().Set("Content-Type", "text/event-stream")
	ginContext.Writer.Header().Set("Cache-Control", "no-cache")
	ginContext.Writer.Header().Set("Connection", "keep-alive")

	tpp.sender = ginContext.Writer
	flusher := ginContext.Writer.(http.Flusher)

	tpp.flusher = flusher
	return true
}

func (tpp *agentChatFlusher) AgentMessageToStreamResponse(answer string) error {
	messageChunkResponse := &biz_entity.AgentMessageStreamResponse{
		ID:     tpp.message.ID,
		Answer: answer,
		StreamResponse: &biz_entity.StreamResponse{
			TaskID: tpp.GetTaskID(),
			Event:  biz_entity.StreamEventAgentMessage,
		},
	}

	chatBotResponse := biz_entity.NewAgentChatBotAppStreamResponse(tpp.GetConversationID(), tpp.message.ID, tpp.message.CreatedAt, messageChunkResponse)

	streamBytes, err := json.Marshal(chatBotResponse)

	if err != nil {
		return err
	}

	if err := tpp.flush(string(streamBytes)); err != nil {
		return err
	}

	return nil
}

func (tpp *agentChatFlusher) AgentThoughtToStreamResponse(ctx context.Context, agentThoughtID string) error {
	agentThought, err := tpp.agentRepo.GetAgentThoughtByID(ctx, agentThoughtID)

	if err != nil {
		return err
	}

	thoughtResp := &biz_entity.AgentThoughtStreamResponse{
		StreamResponse: &biz_entity.StreamResponse{
			TaskID: tpp.GetTaskID(),
			Event:  biz_entity.StreamEventAgentThought,
		},
		Position:     agentThought.Position,
		Thought:      agentThought.Thought,
		Observation:  agentThought.Observation,
		Tool:         agentThought.Tool,
		ToolLabels:   agentThought.ToolLabelsStr,
		ToolInputs:   agentThought.ToolInput,
		MessageFiles: agentThought.MessageFiles,
	}

	streamBytes, err := json.Marshal(thoughtResp)

	if err != nil {
		return err
	}

	if err := tpp.flush(string(streamBytes)); err != nil {
		return err
	}

	return nil
}

func (tpp *agentChatFlusher) ManualFlush(streamString string) error {
	if _, err := fmt.Fprintf(tpp.sender, "data: %s\n\n", streamString); err != nil {
		return err
	}
	tpp.flusher.Flush()
	return nil
}
