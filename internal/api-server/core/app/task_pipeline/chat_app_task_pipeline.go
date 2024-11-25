package task_pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/api-server/core/app"
	appConfigEntities "github.com/lunarianss/Luna/internal/api-server/core/app/app_config/entities"
	"github.com/lunarianss/Luna/internal/api-server/core/app/task_pipeline/entities"
	"github.com/lunarianss/Luna/internal/api-server/entities/llm"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/repo"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
	"github.com/lunarianss/Luna/pkg/log"
)

type ChatAppTaskPipeline struct {
	ApplicationGenerateEntity *app.ChatAppGenerateEntity
	StreamResultChunkQueue    chan *appConfigEntities.MessageQueueMessage
	StreamFinalChunkQueue     chan *appConfigEntities.MessageQueueMessage
	Message                   *model.Message
	MessageRepo               repo.MessageRepo
	flusher                   http.Flusher
	sender                    io.Writer
	taskState                 *entities.ChatAppTaskState
}

func NewChatAppTaskPipeline(
	applicationGenerateEntity *app.ChatAppGenerateEntity,
	streamResultChunkQueue chan *appConfigEntities.MessageQueueMessage,
	streamFinalChunkQueue chan *appConfigEntities.MessageQueueMessage,
	messageRepo repo.MessageRepo, message *model.Message) *ChatAppTaskPipeline {
	return &ChatAppTaskPipeline{
		ApplicationGenerateEntity: applicationGenerateEntity,
		StreamResultChunkQueue:    streamResultChunkQueue,
		StreamFinalChunkQueue:     streamFinalChunkQueue,
		Message:                   message,
		MessageRepo:               messageRepo,
		taskState: &entities.ChatAppTaskState{
			LLMResult: llm.NewEmptyLLMResult(),
		},
	}
}

func (tpp *ChatAppTaskPipeline) Process(ctx context.Context, stream bool) {
	if !tpp.setFlush(ctx) {
		return
	}

	if stream {
		tpp.process_stream_response(ctx)
	}
}

func (tpp *ChatAppTaskPipeline) flush(streamString string) error {
	if _, err := fmt.Fprintf(tpp.sender, "data: %s\n\n", streamString); err != nil {
		return err
	}
	tpp.flusher.Flush()
	return nil
}

func (tpp *ChatAppTaskPipeline) setFlush(c context.Context) bool {

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

func (tpp *ChatAppTaskPipeline) sendFallBackMessageEnd() {
	if err := tpp.flush("data: {\"event\": \"message_end\"}\n\n"); err != nil {
		log.Errorf("failed to send fallback message end to stream response: %v", err)
	}
}

func (tpp *ChatAppTaskPipeline) process_stream_chunk_queue() {

	for v := range tpp.StreamResultChunkQueue {
		if chunkEvent, ok := v.Event.(*appConfigEntities.QueueLLMChunkEvent); ok {
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

func (tpp *ChatAppTaskPipeline) process_stream_end_chunk_queue(c context.Context) {
	for v := range tpp.StreamFinalChunkQueue {
		if mc, ok := v.Event.(*appConfigEntities.QueueMessageEndEvent); ok {
			tpp.taskState.LLMResult = mc.LLMResult
			if err := tpp.saveMessage(c); err != nil {
				log.Errorf("failed to save message: %v", err)
			}

			if err := tpp.messageEndToStreamResponse(); err != nil {
				log.Errorf("failed to flush message to stream response: %v", err)
				tpp.sendFallBackMessageEnd()
			}
		} else if mc, ok := v.Event.(*appConfigEntities.QueueErrorEvent); ok {
			log.Errorf("found queue error event: %v", mc.Err)

			if err := tpp.messageErrToStreamResponse(c, mc.Err); err != nil {
				log.Errorf("failed to flush err message to stream response: %v", err)
				tpp.sendFallBackMessageEnd()
			}
		}
	}
}

func (tpp *ChatAppTaskPipeline) process_stream_response(c context.Context) {
	tpp.process_stream_chunk_queue()
	tpp.process_stream_end_chunk_queue(c)
}

func (tpp *ChatAppTaskPipeline) messageChunkToStreamResponse(answer string) error {
	messageChunkResponse := &entities.MessageStreamResponse{
		ID:                   tpp.Message.ID,
		Answer:               answer,
		FromVariableSelector: make([]string, 0),
		StreamResponse: &entities.StreamResponse{
			TaskID: tpp.ApplicationGenerateEntity.TaskID,
			Event:  entities.StreamEventMessage,
		},
	}

	chatBotResponse := entities.NewChatBotAppStreamResponse(tpp.ApplicationGenerateEntity.ConversationID, tpp.Message.ID, tpp.Message.CreatedAt, messageChunkResponse)

	streamBytes, err := json.Marshal(chatBotResponse)

	if err != nil {
		return err
	}

	if err := tpp.flush(string(streamBytes)); err != nil {
		return err
	}

	return nil
}

func (tpp *ChatAppTaskPipeline) messageErrToStreamResponse(ctx context.Context, err error) error {

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

	messageErrResponse := &entities.ErrorStreamResponse{
		StreamResponse: &entities.StreamResponse{
			TaskID: tpp.ApplicationGenerateEntity.TaskID,
			Event:  entities.StreamEventError,
		},
		Err:     errStr,
		Code:    errStr,
		Message: errStr,
		Status:  500,
	}

	chatBotResponse := entities.NewChatBotAppErrStreamResponse(tpp.ApplicationGenerateEntity.ConversationID, tpp.Message.ID, tpp.Message.CreatedAt, messageErrResponse)

	errorStreamBytes, err := json.Marshal(chatBotResponse)

	if err != nil {
		return err
	}

	if err := tpp.flush(string(errorStreamBytes)); err != nil {
		return err
	}

	return nil
}

func (tpp *ChatAppTaskPipeline) messageEndToStreamResponse() error {
	messageEndResponse := &entities.MessageEndStreamResponse{
		ID: tpp.Message.ID,
		StreamResponse: &entities.StreamResponse{
			TaskID: tpp.ApplicationGenerateEntity.TaskID,
			Event:  entities.StreamEventMessageEnd,
		},
	}

	chatBotResponse := entities.NewChatBotAppEndStreamResponse(tpp.ApplicationGenerateEntity.ConversationID, tpp.Message.ID, tpp.Message.CreatedAt, messageEndResponse)

	endStreamBytes, err := json.Marshal(chatBotResponse)

	if err != nil {
		return err
	}

	if err := tpp.flush(string(endStreamBytes)); err != nil {
		return err
	}
	return nil
}

// func (t *ChatAppTaskPipeline) toStreamResponse(c context.Context) {

// }

func (tpp *ChatAppTaskPipeline) saveMessage(c context.Context) error {
	messageRecord, err := tpp.MessageRepo.GetMessageByID(c, tpp.Message.ID)

	if err != nil {
		return err
	}

	messageRecord.Answer = tpp.taskState.LLMResult.Message.Content.(string)

	messageRecord.Message = tpp.taskState.LLMResult.PromptMessage

	if err := tpp.MessageRepo.UpdateMessage(c, messageRecord); err != nil {
		return err
	}

	return nil
}
