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

var _ biz_entity.IStreamGenerateQueue = (*AgentStreamGenerateQueue)(nil)

type AgentStreamGenerateQueue struct {
	// Input
	StreamResultChunkQueue chan *biz_entity.MessageQueueMessage
	StreamFinalChunkQueue  chan *biz_entity.MessageQueueMessage
	StreamErrorQueue       chan *biz_entity.MessageQueueMessage

	// Output
	OutStreamResultChunkQueue chan *biz_entity.MessageQueueMessage
	OutStreamFinalChunkQueue  chan *biz_entity.MessageQueueMessage
	OutStreamErrorChunkQueue  chan *biz_entity.MessageQueueMessage

	// Message Info
	TaskID         string
	UserID         string
	ConversationID string
	MessageID      string
	AppMode        po_entity.AppMode
	InvokeFrom     string

	// Runtime Parameters
	isOccurredErr bool
	isNormalQuit  bool
}

func NewAgentStreamGenerateQueue(taskID, userID, conversationID, messageId string, appMode po_entity.AppMode, invokeFrom string) biz_entity.IStreamGenerateQueue {

	streamResultChan := make(chan *biz_entity.MessageQueueMessage, STREAM_BUFFER_SIZE)
	streamFinalChan := make(chan *biz_entity.MessageQueueMessage, STREAM_BUFFER_SIZE)
	streamErrorQueue := make(chan *biz_entity.MessageQueueMessage, ERROR_BUFFER_SIZE)

	return &AgentStreamGenerateQueue{
		StreamResultChunkQueue:    make(chan *biz_entity.MessageQueueMessage, STREAM_BUFFER_SIZE),
		StreamFinalChunkQueue:     make(chan *biz_entity.MessageQueueMessage, STREAM_BUFFER_SIZE),
		StreamErrorQueue:          make(chan *biz_entity.MessageQueueMessage, ERROR_BUFFER_SIZE),
		OutStreamResultChunkQueue: streamResultChan,
		OutStreamFinalChunkQueue:  streamFinalChan,
		OutStreamErrorChunkQueue:  streamErrorQueue,
		TaskID:                    taskID,
		UserID:                    userID,
		ConversationID:            conversationID,
		MessageID:                 messageId,
		AppMode:                   appMode,
		InvokeFrom:                invokeFrom,
	}
}

func (sgq *AgentStreamGenerateQueue) PushErr(err error) {
	defer sgq.CloseErrChan()

	errEvent := biz_entity.NewAppQueueEvent(biz_entity.Error)

	sgq.StreamErrorQueue <- sgq.constructMessageQueue(&biz_entity.QueueErrorEvent{
		AppQueueEvent: errEvent,
		Err:           err,
	})
}

func (sgq *AgentStreamGenerateQueue) Fork() biz_entity.IStreamGenerateQueue {

	streamResultChan := make(chan *biz_entity.MessageQueueMessage, STREAM_BUFFER_SIZE)
	streamFinalChan := make(chan *biz_entity.MessageQueueMessage, STREAM_BUFFER_SIZE)
	streamErrorQueue := make(chan *biz_entity.MessageQueueMessage, ERROR_BUFFER_SIZE)

	return &AgentStreamGenerateQueue{
		StreamResultChunkQueue:    make(chan *biz_entity.MessageQueueMessage, STREAM_BUFFER_SIZE),
		StreamFinalChunkQueue:     make(chan *biz_entity.MessageQueueMessage, STREAM_BUFFER_SIZE),
		StreamErrorQueue:          make(chan *biz_entity.MessageQueueMessage, ERROR_BUFFER_SIZE),
		OutStreamResultChunkQueue: streamResultChan,
		OutStreamFinalChunkQueue:  streamFinalChan,
		OutStreamErrorChunkQueue:  streamErrorQueue,
		TaskID:                    sgq.TaskID,
		UserID:                    sgq.UserID,
		ConversationID:            sgq.ConversationID,
		MessageID:                 sgq.MessageID,
		AppMode:                   sgq.AppMode,
		InvokeFrom:                sgq.InvokeFrom,
	}
}

func (sgq *AgentStreamGenerateQueue) printInfo(ch chan *biz_entity.MessageQueueMessage, name string) {

	chanLen := len(ch)
	v, ok := <-ch
	log.Infof(color.GreenString("%s: 是否关闭: %v, 剩余容量: %d, 值: %+v", name, !ok, chanLen, v))

}
func (sgq *AgentStreamGenerateQueue) Debug() {

	log.Infof("=========== QUEUE ===========")
	sgq.printInfo(sgq.StreamResultChunkQueue, "StreamResultQueue")
	sgq.printInfo(sgq.StreamFinalChunkQueue, "StreamFinalQueue")
	sgq.printInfo(sgq.StreamErrorQueue, "StreamFinalQueue")

	log.Infof("=========== END QUEUE ===========")
	sgq.printInfo(sgq.OutStreamResultChunkQueue, "OutStreamResultQueue")
	sgq.printInfo(sgq.OutStreamFinalChunkQueue, "OutStreamFinalQueue")
	sgq.printInfo(sgq.OutStreamErrorChunkQueue, "OutStreamFinalQueue")
}

func (sgq *AgentStreamGenerateQueue) Close() {

}

func (sgq *AgentStreamGenerateQueue) Push(chunk biz_entity.IQueueEvent) {
	sgq.StreamResultChunkQueue <- sgq.constructMessageQueue(chunk)
}

func (sgq *AgentStreamGenerateQueue) Final(chunk biz_entity.IQueueEvent) {
	defer sgq.CloseFinalChan()
	sgq.StreamFinalChunkQueue <- sgq.constructMessageQueue(chunk)
}

func (sgq *AgentStreamGenerateQueue) GetQueues() (chan *biz_entity.MessageQueueMessage, chan *biz_entity.MessageQueueMessage, chan *biz_entity.MessageQueueMessage) {
	return sgq.OutStreamResultChunkQueue, sgq.OutStreamFinalChunkQueue, sgq.OutStreamErrorChunkQueue
}

func (sgq *AgentStreamGenerateQueue) constructMessageQueue(chunk biz_entity.IQueueEvent) *biz_entity.MessageQueueMessage {
	return &biz_entity.MessageQueueMessage{
		Event:          chunk,
		TaskID:         sgq.TaskID,
		ConversationID: sgq.ConversationID,
		MessageID:      sgq.MessageID,
		AppMode:        string(sgq.AppMode),
	}
}

func (sgq *AgentStreamGenerateQueue) CloseErrChan() {
	close(sgq.StreamErrorQueue)
}

func (sgq *AgentStreamGenerateQueue) CloseFinalChan() {
	close(sgq.StreamFinalChunkQueue)
}

func (sgq *AgentStreamGenerateQueue) CloseOutFinalChan() {
	close(sgq.OutStreamFinalChunkQueue)
}

func (sgq *AgentStreamGenerateQueue) CloseOutErrChan() {
	close(sgq.OutStreamErrorChunkQueue)
}

func (sgq *AgentStreamGenerateQueue) closeErr() {
	close(sgq.StreamFinalChunkQueue)
	close(sgq.StreamResultChunkQueue)
}

func (sgq *AgentStreamGenerateQueue) CloseOutErr() {
	close(sgq.OutStreamFinalChunkQueue)
	close(sgq.OutStreamResultChunkQueue)
}

func (sgq *AgentStreamGenerateQueue) CloseOutNormalExit() {
	close(sgq.OutStreamErrorChunkQueue)
	close(sgq.OutStreamResultChunkQueue)
}

func (sgq *AgentStreamGenerateQueue) closeNormalExit() {
	close(sgq.StreamErrorQueue)
	close(sgq.StreamResultChunkQueue)
}

func (sgq *AgentStreamGenerateQueue) Listen() {
QuitLoop:
	for {
		select {
		case resultChunk := <-sgq.StreamResultChunkQueue:
			sgq.OutStreamResultChunkQueue <- resultChunk
		case finalChunk, ok := <-sgq.StreamFinalChunkQueue:
			sgq.isNormalQuit = true
			if !ok {
				sgq.CloseOutFinalChan()
				break QuitLoop
			}
			sgq.OutStreamFinalChunkQueue <- finalChunk
		case errChunk, ok := <-sgq.StreamErrorQueue:
			sgq.isOccurredErr = true
			if !ok {
				sgq.CloseOutErrChan()
				break QuitLoop
			}
			sgq.OutStreamErrorChunkQueue <- errChunk
		}
	}

	if sgq.isOccurredErr {
		sgq.handleErrFallback()
		return
	}

	if sgq.isNormalQuit {
		for len(sgq.StreamResultChunkQueue) > 0 {
			resultChunk := <-sgq.StreamResultChunkQueue
			sgq.OutStreamResultChunkQueue <- resultChunk
		}
		sgq.closeNormalExit()
	}
}

func (sgq *AgentStreamGenerateQueue) handleErrFallback() {
	defer sgq.closeErr()

	for len(sgq.StreamResultChunkQueue) > 0 {
		resultChunk := <-sgq.StreamResultChunkQueue
		sgq.OutStreamResultChunkQueue <- resultChunk
	}

	for len(sgq.StreamFinalChunkQueue) > 0 {
		resultChunk := <-sgq.StreamResultChunkQueue
		sgq.OutStreamResultChunkQueue <- resultChunk
	}
}
