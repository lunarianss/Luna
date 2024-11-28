// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model_runtime

import (
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	biz_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	biz_entity_app_generate "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_generate"
)

const (
	STREAM_BUFFER_SIZE = 17
	ERROR_BUFFER_SIZE  = 7
)

type StreamGenerateQueue struct {
	// Input
	StreamResultChunkQueue chan *biz_entity_chat.MessageQueueMessage
	StreamFinalChunkQueue  chan *biz_entity_chat.MessageQueueMessage

	// Output
	OutStreamResultChunkQueue chan *biz_entity_chat.MessageQueueMessage
	OutStreamFinalChunkQueue  chan *biz_entity_chat.MessageQueueMessage

	// Message Info
	TaskID         string
	UserID         string
	ConversationID string
	MessageID      string
	AppMode        po_entity.AppMode
	InvokeFrom     biz_entity_app_generate.InvokeFrom
}

func NewStreamGenerateQueue(taskID, userID, conversationID, messageId string, appMode po_entity.AppMode, invokeFrom biz_entity_app_generate.InvokeFrom) (*StreamGenerateQueue, chan *biz_entity_chat.MessageQueueMessage, chan *biz_entity_chat.MessageQueueMessage) {

	streamResultChan := make(chan *biz_entity_chat.MessageQueueMessage, STREAM_BUFFER_SIZE)
	streamFinalChan := make(chan *biz_entity_chat.MessageQueueMessage, STREAM_BUFFER_SIZE)

	return &StreamGenerateQueue{
		StreamResultChunkQueue:    make(chan *biz_entity_chat.MessageQueueMessage, STREAM_BUFFER_SIZE),
		StreamFinalChunkQueue:     make(chan *biz_entity_chat.MessageQueueMessage, STREAM_BUFFER_SIZE),
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

	errEvent := biz_entity_chat.NewAppQueueEvent(biz_entity_chat.Error)

	sgq.Final(&biz_entity_chat.QueueErrorEvent{
		AppQueueEvent: errEvent,
		Err:           err,
	})
}

func (sgq *StreamGenerateQueue) Push(chunk biz_entity_chat.IQueueEvent) {

	sgq.StreamResultChunkQueue <- sgq.constructMessageQueue(chunk)
}

func (sgq *StreamGenerateQueue) Final(chunk biz_entity_chat.IQueueEvent) {
	sgq.StreamFinalChunkQueue <- sgq.constructMessageQueue(chunk)
}

func (sgq *StreamGenerateQueue) constructMessageQueue(chunk biz_entity_chat.IQueueEvent) *biz_entity_chat.MessageQueueMessage {
	return &biz_entity_chat.MessageQueueMessage{
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
