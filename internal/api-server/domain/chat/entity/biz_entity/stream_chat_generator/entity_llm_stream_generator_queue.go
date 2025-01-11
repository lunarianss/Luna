package biz_entity

// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

import (
	"github.com/fatih/color"
	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity/stream_base_generator"
)

const (
	STREAM_BUFFER_SIZE = 17
	ERROR_BUFFER_SIZE  = 7
)

var _ biz_entity.IStreamGenerateQueue = (*StreamGenerateQueue)(nil)

type StreamGenerateQueue struct {
	// Input
	StreamResultChunkQueue chan *biz_entity.MessageQueueMessage
	StreamFinalChunkQueue  chan *biz_entity.MessageQueueMessage

	// Output
	OutStreamResultChunkQueue chan *biz_entity.MessageQueueMessage
	OutStreamFinalChunkQueue  chan *biz_entity.MessageQueueMessage

	// Message Info
	TaskID         string
	UserID         string
	ConversationID string
	MessageID      string
	AppMode        po_entity.AppMode
	InvokeFrom     string
}

func NewStreamGenerateQueue(taskID, userID, conversationID, messageId string, appMode po_entity.AppMode, invokeFrom string) (*StreamGenerateQueue, chan *biz_entity.MessageQueueMessage, chan *biz_entity.MessageQueueMessage) {

	streamResultChan := make(chan *biz_entity.MessageQueueMessage, STREAM_BUFFER_SIZE)
	streamFinalChan := make(chan *biz_entity.MessageQueueMessage, STREAM_BUFFER_SIZE)

	return &StreamGenerateQueue{
		StreamResultChunkQueue:    make(chan *biz_entity.MessageQueueMessage, STREAM_BUFFER_SIZE),
		StreamFinalChunkQueue:     make(chan *biz_entity.MessageQueueMessage, STREAM_BUFFER_SIZE),
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

	errEvent := biz_entity.NewAppQueueEvent(biz_entity.Error)

	sgq.StreamFinalChunkQueue <- sgq.constructMessageQueue(&biz_entity.QueueErrorEvent{
		AppQueueEvent: errEvent,
		Err:           err,
	})
}

func (sgq *StreamGenerateQueue) Push(chunk biz_entity.IQueueEvent) {
	sgq.StreamResultChunkQueue <- sgq.constructMessageQueue(chunk)
}

func (sgq *StreamGenerateQueue) Fork() biz_entity.IStreamGenerateQueue {
	return sgq
}

func (sgq *StreamGenerateQueue) CloseOutErr() {

}

func (sgq *StreamGenerateQueue) CloseOutNormalExit() {

}

func (sgq *StreamGenerateQueue) printInfo(ch chan *biz_entity.MessageQueueMessage, name string) {

	chanLen := len(ch)
	v, ok := <-ch
	log.Infof(color.GreenString("%s: 是否关闭: %v, 剩余容量: %d, 值: %+v", name, !ok, chanLen, v))

}
func (sgq *StreamGenerateQueue) Debug() {

	log.Infof("=========== QUEUE ===========")
	sgq.printInfo(sgq.StreamResultChunkQueue, "StreamResultQueue")
	sgq.printInfo(sgq.StreamFinalChunkQueue, "StreamFinalQueue")

	log.Infof("=========== END QUEUE ===========")
	sgq.printInfo(sgq.OutStreamResultChunkQueue, "OutStreamResultQueue")
	sgq.printInfo(sgq.OutStreamFinalChunkQueue, "OutStreamFinalQueue")
}
func (sgq *StreamGenerateQueue) GetQueues() (chan *biz_entity.MessageQueueMessage, chan *biz_entity.MessageQueueMessage, chan *biz_entity.MessageQueueMessage) {
	return sgq.OutStreamResultChunkQueue, sgq.StreamFinalChunkQueue, nil
}

func (sgq *StreamGenerateQueue) Final(chunk biz_entity.IQueueEvent) {
	defer sgq.Close()
	sgq.StreamFinalChunkQueue <- sgq.constructMessageQueue(chunk)
}

func (sgq *StreamGenerateQueue) constructMessageQueue(chunk biz_entity.IQueueEvent) *biz_entity.MessageQueueMessage {
	return &biz_entity.MessageQueueMessage{
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
