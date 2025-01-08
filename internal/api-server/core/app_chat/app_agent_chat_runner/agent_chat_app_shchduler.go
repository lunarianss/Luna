package app_agent_chat_runner

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/infrastructure/log"
	repo_agent "github.com/lunarianss/Luna/internal/api-server/domain/agent/repository"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/repository"
	biz_entity_app_generate "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_generate"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
)

type IAgentChatAppTaskScheduler interface {
	Process(ctx context.Context)
	SetFunctionCallRunner(*FunctionCallAgentRunner)
}

type agentChatAppTaskScheduler struct {
	biz_entity_app_generate.BasedAppGenerateEntity
	StreamResultChunkQueue chan *biz_entity.MessageQueueMessage
	StreamFinalChunkQueue  chan *biz_entity.MessageQueueMessage
	StreamErrorChunkQueue  chan *biz_entity.MessageQueueMessage
	Message                *po_entity.Message
	MessageRepo            repository.MessageRepo
	AnnotationRepo         repository.AnnotationRepo
	AgentRepo              repo_agent.AgentRepo

	flusher   http.Flusher
	sender    io.Writer
	taskState *biz_entity.ChatAppTaskState
	runner    *FunctionCallAgentRunner
}

func NewAgentChatAppTaskScheduler(
	applicationGenerateEntity biz_entity_app_generate.BasedAppGenerateEntity,
	messageRepo repository.MessageRepo, message *po_entity.Message, annotationRepo repository.AnnotationRepo, runner *FunctionCallAgentRunner) *agentChatAppTaskScheduler {
	return &agentChatAppTaskScheduler{
		BasedAppGenerateEntity: applicationGenerateEntity,
		Message:                message,
		MessageRepo:            messageRepo,
		AnnotationRepo:         annotationRepo,
		runner:                 runner,
		taskState: &biz_entity.ChatAppTaskState{
			LLMResult: biz_entity.NewEmptyLLMResult(),
		},
	}
}

func (tpp *agentChatAppTaskScheduler) SetFunctionCallRunner(runner *FunctionCallAgentRunner) {
	tpp.runner = runner
}

func (tpp *agentChatAppTaskScheduler) Process(ctx context.Context) {
	if !tpp.setFlush(ctx) {
		return
	}

	var err error

	tpp.taskState, err = tpp.runner.Run(ctx, tpp.Message, tpp.Message.Query)

	if err != nil {
		if err := tpp.messageErrToStreamResponse(ctx, err); err != nil {
			log.Errorf("failed to flush message to stream response: %v", err)
			tpp.sendFallBackMessageEnd()
		}
	}
	if err := tpp.saveMessage(ctx); err != nil {
		log.Errorf("failed to save message: %v", err)
	}

	if err := tpp.messageEndToStreamResponse(); err != nil {
		log.Errorf("failed to flush message to stream response: %v", err)
		tpp.sendFallBackMessageEnd()
	}
}

func (tpp *agentChatAppTaskScheduler) flush(streamString string) error {
	if _, err := fmt.Fprintf(tpp.sender, "data: %s\n\n", streamString); err != nil {
		return err
	}
	tpp.flusher.Flush()
	return nil
}

func (tpp *agentChatAppTaskScheduler) setFlush(c context.Context) bool {

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

func (tpp *agentChatAppTaskScheduler) sendFallBackMessageEnd() {
	if err := tpp.flush("data: {\"event\": \"message_end\"}\n\n"); err != nil {
		log.Errorf("failed to send fallback message end to stream response: %v", err)
	}
}

func (tpp *agentChatAppTaskScheduler) messageErrToStreamResponse(ctx context.Context, err error) error {

	var errStr = "Internal Server Error, please contact support."

	if errors.IsCode(err, code.ErrQuotaExceed) {
		errStr = "Your quota for Luna Hosted Model Provider has been exhausted. Please go to Settings -> Model Provider to complete your own provider credentials."
	}

	messageRecord, err := tpp.MessageRepo.GetMessageByID(ctx, tpp.Message.ID)

	if err != nil {
		return err
	}

	messageRecord.Status = "error"
	messageRecord.Error = errStr

	if err := tpp.MessageRepo.UpdateMessage(ctx, messageRecord); err != nil {
		return err
	}

	messageErrResponse := &biz_entity.ErrorStreamResponse{
		StreamResponse: &biz_entity.StreamResponse{
			TaskID: tpp.GetTaskID(),
			Event:  biz_entity.StreamEventError,
		},
		Err:     errStr,
		Code:    errStr,
		Message: errStr,
		Status:  500,
	}

	chatBotResponse := biz_entity.NewChatBotAppErrStreamResponse(tpp.GetConversationID(), tpp.Message.ID, tpp.Message.CreatedAt, messageErrResponse)

	errorStreamBytes, err := json.Marshal(chatBotResponse)

	if err != nil {
		return err
	}

	if err := tpp.flush(string(errorStreamBytes)); err != nil {
		return err
	}

	return nil
}

func (tpp *agentChatAppTaskScheduler) messageEndToStreamResponse() error {
	messageEndResponse := &biz_entity.MessageEndStreamResponse{
		ID: tpp.Message.ID,
		StreamResponse: &biz_entity.StreamResponse{
			TaskID: tpp.GetTaskID(),
			Event:  biz_entity.StreamEventMessageEnd,
		},
		Metadata: &biz_entity.MetaDataUsage{
			Usage: tpp.taskState.LLMResult.Usage,
		},
		MessageId:      tpp.Message.ID,
		ConversationID: tpp.Message.ConversationID,
	}

	if annotation_reply, ok := tpp.taskState.Metadata["annotation_reply"]; ok {
		if annotation_replyMap, ok := annotation_reply.(map[string]any); ok {
			messageEndResponse.Metadata.AnnotationReply = annotation_replyMap
		}
	}

	chatBotResponse := biz_entity.NewChatBotAppEndStreamResponse(tpp.GetConversationID(), tpp.Message.ID, tpp.Message.CreatedAt, messageEndResponse)

	endStreamBytes, err := json.Marshal(chatBotResponse)

	if err != nil {
		return err
	}

	if err := tpp.flush(string(endStreamBytes)); err != nil {
		return err
	}
	return nil
}

func (tpp *agentChatAppTaskScheduler) saveMessage(c context.Context) error {
	messageRecord, err := tpp.MessageRepo.GetMessageByID(c, tpp.Message.ID)

	if err != nil {
		return err
	}

	messageRecord.Answer = tpp.taskState.LLMResult.Message.Content.(string)

	messageRecord.Message = util.ConvertToInterfaceSlice(tpp.taskState.LLMResult.PromptMessage, func(v po_entity.IPromptMessage) any {
		return any(v)
	})
	messageRecord.MessageTokens = tpp.taskState.LLMResult.Usage.PromptTokens
	messageRecord.MessagePriceUnit = tpp.taskState.LLMResult.Usage.PromptPriceUnit
	messageRecord.MessageUnitPrice = tpp.taskState.LLMResult.Usage.PromptUnitPrice
	messageRecord.AnswerTokens = tpp.taskState.LLMResult.Usage.CompletionTokens
	messageRecord.AnswerPriceUnit = tpp.taskState.LLMResult.Usage.CompletionPriceUnit
	messageRecord.AnswerUnitPrice = tpp.taskState.LLMResult.Usage.CompletionUnitPrice
	messageRecord.TotalPrice = tpp.taskState.LLMResult.Usage.TotalPrice
	messageRecord.Currency = tpp.taskState.LLMResult.Usage.Currency
	messageRecord.MessageMetadata = tpp.taskState.Metadata

	if err := tpp.MessageRepo.UpdateMessage(c, messageRecord); err != nil {
		return err
	}

	return nil
}
