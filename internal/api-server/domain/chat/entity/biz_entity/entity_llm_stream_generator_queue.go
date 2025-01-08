package biz_entity

// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

import (
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
)

const (
	STREAM_BUFFER_SIZE = 17
	ERROR_BUFFER_SIZE  = 7
)

type IStreamGenerateQueue interface {
	PushErr(err error)
	Push(chunk IQueueEvent)
	Final(chunk IQueueEvent)
	Close()
}

type StreamGenerateQueue struct {
	// Input
	StreamResultChunkQueue chan *MessageQueueMessage
	StreamFinalChunkQueue  chan *MessageQueueMessage

	// Output
	OutStreamResultChunkQueue chan *MessageQueueMessage
	OutStreamFinalChunkQueue  chan *MessageQueueMessage

	// Message Info
	TaskID         string
	UserID         string
	ConversationID string
	MessageID      string
	AppMode        po_entity.AppMode
	InvokeFrom     string
}

func NewStreamGenerateQueue(taskID, userID, conversationID, messageId string, appMode po_entity.AppMode, invokeFrom string) (*StreamGenerateQueue, chan *MessageQueueMessage, chan *MessageQueueMessage) {

	streamResultChan := make(chan *MessageQueueMessage, STREAM_BUFFER_SIZE)
	streamFinalChan := make(chan *MessageQueueMessage, STREAM_BUFFER_SIZE)

	return &StreamGenerateQueue{
		StreamResultChunkQueue:    make(chan *MessageQueueMessage, STREAM_BUFFER_SIZE),
		StreamFinalChunkQueue:     make(chan *MessageQueueMessage, STREAM_BUFFER_SIZE),
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

	errEvent := NewAppQueueEvent(Error)

	sgq.StreamFinalChunkQueue <- sgq.constructMessageQueue(&QueueErrorEvent{
		AppQueueEvent: errEvent,
		Err:           err,
	})
}

func (sgq *StreamGenerateQueue) Push(chunk IQueueEvent) {
	sgq.StreamResultChunkQueue <- sgq.constructMessageQueue(chunk)
}

func (sgq *StreamGenerateQueue) Final(chunk IQueueEvent) {
	defer sgq.Close()
	sgq.StreamFinalChunkQueue <- sgq.constructMessageQueue(chunk)
}

func (sgq *StreamGenerateQueue) FinalManual(chunk IQueueEvent) {
	defer sgq.Close()
	sgq.StreamFinalChunkQueue <- sgq.constructMessageQueue(chunk)
}

func (sgq *StreamGenerateQueue) constructMessageQueue(chunk IQueueEvent) *MessageQueueMessage {
	return &MessageQueueMessage{
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
