// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model_runtime

import (
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config/entities"
	appEntities "github.com/lunarianss/Luna/internal/api-server/core/app/apps/entities"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
)

const (
	STREAM_BUFFER_SIZE = 17
	ERROR_BUFFER_SIZE  = 7
)

type StreamGenerateQueue struct {
	// Input
	StreamResultChunkQueue chan *entities.MessageQueueMessage
	StreamFinalChunkQueue  chan *entities.MessageQueueMessage

	// Output
	OutStreamResultChunkQueue chan *entities.MessageQueueMessage
	OutStreamFinalChunkQueue  chan *entities.MessageQueueMessage

	// Message Info
	TaskID         string
	UserID         string
	ConversationID string
	MessageID      string
	AppMode        po_entity.AppMode
	InvokeFrom     appEntities.InvokeForm
}

func NewStreamGenerateQueue(taskID, userID, conversationID, messageId string, appMode po_entity.AppMode, invokeFrom appEntities.InvokeForm) (*StreamGenerateQueue, chan *entities.MessageQueueMessage, chan *entities.MessageQueueMessage) {

	streamResultChan := make(chan *entities.MessageQueueMessage, STREAM_BUFFER_SIZE)
	streamFinalChan := make(chan *entities.MessageQueueMessage, STREAM_BUFFER_SIZE)

	return &StreamGenerateQueue{
		StreamResultChunkQueue:    make(chan *entities.MessageQueueMessage, STREAM_BUFFER_SIZE),
		StreamFinalChunkQueue:     make(chan *entities.MessageQueueMessage, STREAM_BUFFER_SIZE),
		OutStreamResultChunkQueue: streamResultChan,
		OutStreamFinalChunkQueue:  streamFinalChan,
		TaskID:                    taskID,
		UserID:                    userID,
		ConversationID:            conversationID,
		MessageID:                 messageId,
		AppMode:                   appMode,
		InvokeFrom:                invokeFrom,
	}, streamResultChan, streamFinalChan
}

func (sgq *StreamGenerateQueue) PushErr(err error) {
	defer sgq.Close()

	errEvent := entities.NewAppQueueEvent(entities.Error)

	sgq.Final(&entities.QueueErrorEvent{
		AppQueueEvent: errEvent,
		Err:           err,
	})
}

func (sgq *StreamGenerateQueue) Push(chunk entities.IQueueEvent) {

	sgq.StreamResultChunkQueue <- sgq.constructMessageQueue(chunk)
}

func (sgq *StreamGenerateQueue) Final(chunk entities.IQueueEvent) {
	sgq.StreamFinalChunkQueue <- sgq.constructMessageQueue(chunk)
}

func (sgq *StreamGenerateQueue) constructMessageQueue(chunk entities.IQueueEvent) *entities.MessageQueueMessage {
	return &entities.MessageQueueMessage{
		Event:          chunk,
		TaskID:         sgq.TaskID,
		ConversationID: sgq.ConversationID,
		MessageID:      sgq.MessageID,
		AppMode:        string(sgq.AppMode),
	}
}

func (sgq *StreamGenerateQueue) Close() {
	close(sgq.StreamResultChunkQueue)
	close(sgq.StreamFinalChunkQueue)
}

func (sgq *StreamGenerateQueue) CloseOut() {
	close(sgq.OutStreamResultChunkQueue)
	close(sgq.OutStreamFinalChunkQueue)
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
