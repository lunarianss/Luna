package model_runtime

import (
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config/entities"
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
}

func NewStreamGenerateQueue(streamResultChan chan entities.IQueueEvent, streamFinalChan chan entities.IQueueEvent) *StreamGenerateQueue {
	return &StreamGenerateQueue{
		StreamResultChunkQueue:    make(chan entities.IQueueEvent, STREAM_BUFFER_SIZE),
		StreamFinalChunkQueue:     make(chan entities.IQueueEvent, STREAM_BUFFER_SIZE),
		OutStreamResultChunkQueue: streamResultChan,
		OutStreamFinalChunkQueue:  streamFinalChan,
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
