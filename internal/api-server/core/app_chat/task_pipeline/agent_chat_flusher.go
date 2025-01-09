package task_pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	biz_agent "github.com/lunarianss/Luna/internal/api-server/domain/agent/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"

	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/internal/api-server/domain/agent/domain_service"
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

var _ biz_agent.AgentFlusher = (*agentChatFlusher)(nil)

func NewAgentChatFlusher(
	applicationGenerateEntity biz_entity_app_generate.BasedAppGenerateEntity,
	agentRepo repo_agent.AgentRepo, message *po_entity.Message) biz_agent.AgentFlusher {
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

func (tpp *agentChatFlusher) AgentMessageFileToStreamResponse(ctx context.Context, messageFileID string, secretKey string, baseUrl string) error {

	messageFile, err := tpp.agentRepo.GetMessageFileByID(ctx, messageFileID)

	if err != nil {
		return err
	}

	messageFileUrls := strings.Split(messageFile.URL, "/")

	messageFileUrl := messageFileUrls[len(messageFileUrls)-1]

	toolFilename := strings.Split(messageFileUrl, ".")

	toolFileID := toolFilename[0]

	extension := toolFilename[1]

	url := ""

	if extension == "" || len(extension) > 10 {
		extension = ".bin"
	}

	if strings.HasPrefix(messageFile.URL, "http") || strings.HasPrefix(messageFile.URL, "https") {
		url = messageFile.URL
	} else {

		signedUrl, err := domain_service.NewToolFileManager(nil).SignFile(toolFileID, extension, "", "")

		if err != nil {
			return err
		}

		url = signedUrl

	}

	messageFileChunkResponse := &biz_entity.MessageFileStreamResponse{
		ID:        messageFile.ID,
		Type:      messageFile.Type,
		BelongsTo: messageFile.BelongsTo,
		URL:       url,
		StreamResponse: &biz_entity.StreamResponse{
			TaskID: tpp.GetTaskID(),
			Event:  biz_entity.StreamEventMessageFile,
		},
	}

	chatBotResponse := biz_entity.NewChatBotAppMessageFileStreamResponse(tpp.GetConversationID(), tpp.message.ID, tpp.message.CreatedAt, messageFileChunkResponse)

	streamBytes, err := json.Marshal(chatBotResponse)

	if err != nil {
		return err
	}

	if err := tpp.flush(string(streamBytes)); err != nil {
		return err
	}

	return nil
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
