package model_runtime

import (
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config/entities"
	appEntities "github.com/lunarianss/Luna/internal/api-server/core/app/apps/entities"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
)

const (
	STREAM_BUFFER_SIZE = 17
	ERROR_BUFFER_SIZE  = 7
)

type StreamGenerateQueue struct {
	// Input
	StreamResultChunkQueue chan entities.IQueueEvent
	StreamFinalChunkQueue  chan entities.IQueueEvent

	// Output
	OutStreamResultChunkQueue chan entities.IQueueEvent
	OutStreamFinalChunkQueue  chan entities.IQueueEvent

	// Message Info
	TaskID         string
	UserID         string
	ConversationID string
	MessageID      string
	AppMode        model.AppMode
	InvokeFrom     appEntities.InvokeForm
}

func NewStreamGenerateQueue(streamResultChan chan entities.IQueueEvent, streamFinalChan chan entities.IQueueEvent, taskID, userID, conversationID, messageId string, appMode model.AppMode, invokeFrom appEntities.InvokeForm) *StreamGenerateQueue {
	return &StreamGenerateQueue{
		StreamResultChunkQueue:    make(chan entities.IQueueEvent, STREAM_BUFFER_SIZE),
		StreamFinalChunkQueue:     make(chan entities.IQueueEvent, STREAM_BUFFER_SIZE),
		OutStreamResultChunkQueue: streamResultChan,
		OutStreamFinalChunkQueue:  streamFinalChan,
		TaskID:                    taskID,
		UserID:                    userID,
		ConversationID:            conversationID,
		MessageID:                 messageId,
		AppMode:                   appMode,
		InvokeFrom:                invokeFrom,
	}
}

func (sgq *StreamGenerateQueue) Push(chunk entities.IQueueEvent) {
	sgq.StreamResultChunkQueue <- chunk
}

func (sgq *StreamGenerateQueue) Final(chunk entities.IQueueEvent) {
	sgq.StreamFinalChunkQueue <- chunk
}

func (sgq *StreamGenerateQueue) Close() {
	close(sgq.StreamResultChunkQueue)
	close(sgq.StreamFinalChunkQueue)
}

func (sgq *StreamGenerateQueue) CloseOut() {
	close(sgq.OutStreamFinalChunkQueue)
	close(sgq.OutStreamResultChunkQueue)
}

func (sgq *StreamGenerateQueue) Listen() {
	defer sgq.CloseOut()

	for v := range sgq.StreamResultChunkQueue {
		sgq.OutStreamResultChunkQueue <- v
	}

	for v := range sgq.StreamFinalChunkQueue {
		sgq.OutStreamFinalChunkQueue <- v
	}

}
