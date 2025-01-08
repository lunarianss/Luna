package biz_entity

// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

import (
	"github.com/fatih/color"
	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
)

var _ IStreamGenerateQueue = (*AgentStreamGenerateQueue)(nil)

type AgentStreamGenerateQueue struct {
	// Input
	StreamResultChunkQueue chan *MessageQueueMessage
	StreamFinalChunkQueue  chan *MessageQueueMessage
	StreamErrorQueue       chan *MessageQueueMessage

	// Output
	OutStreamResultChunkQueue chan *MessageQueueMessage
	OutStreamFinalChunkQueue  chan *MessageQueueMessage
	OutStreamErrorChunkQueue  chan *MessageQueueMessage

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

func NewAgentStreamGenerateQueue(taskID, userID, conversationID, messageId string, appMode po_entity.AppMode, invokeFrom string) IStreamGenerateQueue {

	streamResultChan := make(chan *MessageQueueMessage, STREAM_BUFFER_SIZE)
	streamFinalChan := make(chan *MessageQueueMessage, STREAM_BUFFER_SIZE)
	streamErrorQueue := make(chan *MessageQueueMessage, ERROR_BUFFER_SIZE)

	return &AgentStreamGenerateQueue{
		StreamResultChunkQueue:    make(chan *MessageQueueMessage, STREAM_BUFFER_SIZE),
		StreamFinalChunkQueue:     make(chan *MessageQueueMessage, STREAM_BUFFER_SIZE),
		StreamErrorQueue:          make(chan *MessageQueueMessage, ERROR_BUFFER_SIZE),
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

	errEvent := NewAppQueueEvent(Error)

	sgq.StreamErrorQueue <- sgq.constructMessageQueue(&QueueErrorEvent{
		AppQueueEvent: errEvent,
		Err:           err,
	})
}

func (sgq *AgentStreamGenerateQueue) Fork() IStreamGenerateQueue {

	streamResultChan := make(chan *MessageQueueMessage, STREAM_BUFFER_SIZE)
	streamFinalChan := make(chan *MessageQueueMessage, STREAM_BUFFER_SIZE)
	streamErrorQueue := make(chan *MessageQueueMessage, ERROR_BUFFER_SIZE)

	return &AgentStreamGenerateQueue{
		StreamResultChunkQueue:    make(chan *MessageQueueMessage, STREAM_BUFFER_SIZE),
		StreamFinalChunkQueue:     make(chan *MessageQueueMessage, STREAM_BUFFER_SIZE),
		StreamErrorQueue:          make(chan *MessageQueueMessage, ERROR_BUFFER_SIZE),
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

func (sgq *AgentStreamGenerateQueue) printInfo(ch chan *MessageQueueMessage, name string) {

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

func (sgq *AgentStreamGenerateQueue) Push(chunk IQueueEvent) {
	sgq.StreamResultChunkQueue <- sgq.constructMessageQueue(chunk)
}

func (sgq *AgentStreamGenerateQueue) Final(chunk IQueueEvent) {
	defer sgq.CloseFinalChan()
	sgq.StreamFinalChunkQueue <- sgq.constructMessageQueue(chunk)
}

func (sgq *AgentStreamGenerateQueue) GetQueues() (chan *MessageQueueMessage, chan *MessageQueueMessage, chan *MessageQueueMessage) {
	return sgq.OutStreamResultChunkQueue, sgq.OutStreamFinalChunkQueue, sgq.OutStreamErrorChunkQueue
}

func (sgq *AgentStreamGenerateQueue) constructMessageQueue(chunk IQueueEvent) *MessageQueueMessage {
	return &MessageQueueMessage{
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
