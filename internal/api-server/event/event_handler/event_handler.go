package event_handler

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type MQEventHandler interface {
	Handle(ctx context.Context, message *primitive.MessageExt) (consumer.ConsumeResult, error)
}
