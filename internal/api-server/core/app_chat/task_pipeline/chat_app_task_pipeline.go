package task_pipeline

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
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/repository"
	biz_entity_app_generate "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_generate"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
)

type chatAppTaskPipeline struct {
	ApplicationGenerateEntity *biz_entity_app_generate.ChatAppGenerateEntity
	StreamResultChunkQueue    chan *biz_entity.MessageQueueMessage
	StreamFinalChunkQueue     chan *biz_entity.MessageQueueMessage
	Message                   *po_entity.Message
	MessageRepo               repository.MessageRepo
	flusher                   http.Flusher
	sender                    io.Writer
	taskState                 *biz_entity.ChatAppTaskState
}

func NewChatAppTaskPipeline(
	applicationGenerateEntity *biz_entity_app_generate.ChatAppGenerateEntity,
	streamResultChunkQueue chan *biz_entity.MessageQueueMessage,
	streamFinalChunkQueue chan *biz_entity.MessageQueueMessage,
	messageRepo repository.MessageRepo, message *po_entity.Message) *chatAppTaskPipeline {
	return &chatAppTaskPipeline{
		ApplicationGenerateEntity: applicationGenerateEntity,
		StreamResultChunkQueue:    streamResultChunkQueue,
		StreamFinalChunkQueue:     streamFinalChunkQueue,
		Message:                   message,
		MessageRepo:               messageRepo,
		taskState: &biz_entity.ChatAppTaskState{
			LLMResult: biz_entity.NewEmptyLLMResult(),
		},
	}
}

func (tpp *chatAppTaskPipeline) Process(ctx context.Context) {
	if !tpp.setFlush(ctx) {
		return
	}

	tpp.process_stream_response(ctx)

}

func (tpp *chatAppTaskPipeline) flush(streamString string) error {
	if _, err := fmt.Fprintf(tpp.sender, "data: %s\n\n", streamString); err != nil {
		return err
	}
	tpp.flusher.Flush()
	return nil
}

func (tpp *chatAppTaskPipeline) setFlush(c context.Context) bool {

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

func (tpp *chatAppTaskPipeline) sendFallBackMessageEnd() {
	if err := tpp.flush("data: {\"event\": \"message_end\"}\n\n"); err != nil {
		log.Errorf("failed to send fallback message end to stream response: %v", err)
	}
}

func (tpp *chatAppTaskPipeline) process_stream_chunk_queue() {

	for v := range tpp.StreamResultChunkQueue {
		if chunkEvent, ok := v.Event.(*biz_entity.QueueLLMChunkEvent); ok {
			deltaText := chunkEvent.Chunk.Delta.Message.Content

			if content, ok := tpp.taskState.LLMResult.Message.Content.(string); ok {
				tpp.taskState.LLMResult.Message.Content = deltaText.(string) + content
			}
			if err := tpp.messageChunkToStreamResponse(deltaText.(string)); err != nil {
				log.Errorf("failed to flush message to stream response: %v", err)
				tpp.sendFallBackMessageEnd()
			}
		}
	}
}

func (tpp *chatAppTaskPipeline) process_stream_end_chunk_queue(c context.Context) {
	for v := range tpp.StreamFinalChunkQueue {
		if mc, ok := v.Event.(*biz_entity.QueueMessageEndEvent); ok {
			tpp.taskState.LLMResult = mc.LLMResult
			if err := tpp.saveMessage(c); err != nil {
				log.Errorf("failed to save message: %v", err)
			}

			if err := tpp.messageEndToStreamResponse(); err != nil {
				log.Errorf("failed to flush message to stream response: %v", err)
				tpp.sendFallBackMessageEnd()
			}
		} else if mc, ok := v.Event.(*biz_entity.QueueErrorEvent); ok {
			log.Errorf("found queue error event: %#+v", mc.Err)

			if err := tpp.messageErrToStreamResponse(c, mc.Err); err != nil {
				log.Errorf("failed to flush err message to stream response: %v", err)
				tpp.sendFallBackMessageEnd()
			}
		}
	}
}

func (tpp *chatAppTaskPipeline) process_stream_response(c context.Context) {
	tpp.process_stream_chunk_queue()
	tpp.process_stream_end_chunk_queue(c)
}

func (tpp *chatAppTaskPipeline) messageChunkToStreamResponse(answer string) error {
	messageChunkResponse := &biz_entity.MessageStreamResponse{
		ID:                   tpp.Message.ID,
		Answer:               answer,
		FromVariableSelector: make([]string, 0),
		StreamResponse: &biz_entity.StreamResponse{
			TaskID: tpp.ApplicationGenerateEntity.TaskID,
			Event:  biz_entity.StreamEventMessage,
		},
	}

	chatBotResponse := biz_entity.NewChatBotAppStreamResponse(tpp.ApplicationGenerateEntity.ConversationID, tpp.Message.ID, tpp.Message.CreatedAt, messageChunkResponse)

	streamBytes, err := json.Marshal(chatBotResponse)

	if err != nil {
		return err
	}

	if err := tpp.flush(string(streamBytes)); err != nil {
		return err
	}

	return nil
}

func (tpp *chatAppTaskPipeline) messageErrToStreamResponse(ctx context.Context, err error) error {

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
			TaskID: tpp.ApplicationGenerateEntity.TaskID,
			Event:  biz_entity.StreamEventError,
		},
		Err:     errStr,
		Code:    errStr,
		Message: errStr,
		Status:  500,
	}

	chatBotResponse := biz_entity.NewChatBotAppErrStreamResponse(tpp.ApplicationGenerateEntity.ConversationID, tpp.Message.ID, tpp.Message.CreatedAt, messageErrResponse)

	errorStreamBytes, err := json.Marshal(chatBotResponse)

	if err != nil {
		return err
	}

	if err := tpp.flush(string(errorStreamBytes)); err != nil {
		return err
	}

	return nil
}

func (tpp *chatAppTaskPipeline) messageEndToStreamResponse() error {
	messageEndResponse := &biz_entity.MessageEndStreamResponse{
		ID: tpp.Message.ID,
		StreamResponse: &biz_entity.StreamResponse{
			TaskID: tpp.ApplicationGenerateEntity.TaskID,
			Event:  biz_entity.StreamEventMessageEnd,
		},
		Metadata: &biz_entity.MetaDataUsage{
			Usage: tpp.taskState.LLMResult.Usage,
		},
		MessageId:      tpp.Message.ID,
		ConversationID: tpp.Message.ConversationID,
	}

	chatBotResponse := biz_entity.NewChatBotAppEndStreamResponse(tpp.ApplicationGenerateEntity.ConversationID, tpp.Message.ID, tpp.Message.CreatedAt, messageEndResponse)

	endStreamBytes, err := json.Marshal(chatBotResponse)

	if err != nil {
		return err
	}

	if err := tpp.flush(string(endStreamBytes)); err != nil {
		return err
	}
	return nil
}

// func (t *chatAppTaskPipeline) toStreamResponse(c context.Context) {

// }

func (tpp *chatAppTaskPipeline) saveMessage(c context.Context) error {
	messageRecord, err := tpp.MessageRepo.GetMessageByID(c, tpp.Message.ID)

	if err != nil {
		return err
	}

	messageRecord.Answer = tpp.taskState.LLMResult.Message.Content.(string)

	messageRecord.Message = tpp.taskState.LLMResult.PromptMessage
	messageRecord.MessageTokens = tpp.taskState.LLMResult.Usage.PromptTokens
	messageRecord.MessagePriceUnit = tpp.taskState.LLMResult.Usage.PromptPriceUnit
	messageRecord.MessageUnitPrice = tpp.taskState.LLMResult.Usage.PromptUnitPrice
	messageRecord.AnswerTokens = tpp.taskState.LLMResult.Usage.CompletionTokens
	messageRecord.AnswerPriceUnit = tpp.taskState.LLMResult.Usage.CompletionPriceUnit
	messageRecord.AnswerUnitPrice = tpp.taskState.LLMResult.Usage.CompletionUnitPrice
	messageRecord.TotalPrice = tpp.taskState.LLMResult.Usage.TotalPrice
	messageRecord.Currency = tpp.taskState.LLMResult.Usage.Currency

	if err := tpp.MessageRepo.UpdateMessage(c, messageRecord); err != nil {
		return err
	}

	return nil
}
