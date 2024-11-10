package chat

import (
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config/entities"
	appEntities "github.com/lunarianss/Luna/internal/api-server/core/app/apps/entities"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime"
)

type ChatAppGenerator struct {
}

func (g *ChatAppGenerator) Generate(appModel *model.App, user interface{}, args interface{}, invokeFrom appEntities.InvokeForm, stream bool) error {

	StreamResultChunkQueue := make(chan entities.IQueueEvent, model_runtime.STREAM_BUFFER_SIZE)
	StreamFinalChunkQueue := make(chan entities.IQueueEvent, model_runtime.STREAM_BUFFER_SIZE)
	streamChanQueue := model_runtime.NewStreamGenerateQueue(StreamResultChunkQueue, StreamFinalChunkQueue)

	go func() {
		streamChanQueue.Listen()
	}()

	return nil
}
