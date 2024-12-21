package mq

import (
	"fmt"
	"sync"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/infrastructure/mq"
	"github.com/lunarianss/Luna/internal/infrastructure/options"
)

var (
	once          sync.Once
	consumerOnce  sync.Once
	MQProducerIns rocketmq.Producer
	MQConsumerIns rocketmq.PushConsumer
)

func GetMQProducerIns(opt *options.RocketMQOptions) (rocketmq.Producer, error) {
	var (
		err        error
		mqProducer *mq.MQProducer
	)

	once.Do(func() {
		mqProducer, err = mq.NewProducer(opt)
		if err != nil {
			log.Error(err)
		}
		MQProducerIns = mqProducer.GetProducer()
	})

	if MQProducerIns == nil || err != nil {
		return nil, fmt.Errorf("failed to get mq producer factory, mqFactory: %+v, error: %w", MQProducerIns, err)
	}

	return MQProducerIns, nil

}

func GetMQConsumerIns(opt *options.RocketMQOptions) (rocketmq.PushConsumer, error) {
	var (
		err        error
		mqConsumer *mq.MQConsumer
	)

	consumerOnce.Do(func() {
		mqConsumer, err = mq.NewConsumer(opt)
		if err != nil {
			log.Error(err)
		}
		MQConsumerIns = mqConsumer.GetConsumer()
	})

	if MQConsumerIns == nil || err != nil {
		return nil, fmt.Errorf("failed to get mq consumer factory, mqFactory: %+v, error: %w", MQConsumerIns, err)
	}

	return MQConsumerIns, nil

}
