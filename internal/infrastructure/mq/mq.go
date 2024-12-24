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
	once                   sync.Once
	authConsumerOnce       sync.Once
	annotationConsumerOnce sync.Once

	MQProducerIns                rocketmq.Producer
	MQAuthTopicConsumerIns       rocketmq.PushConsumer
	MQAnnotationTopicConsumerIns rocketmq.PushConsumer
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

func GetMQAuthTopicConsumerIns(opt *options.RocketMQOptions) (rocketmq.PushConsumer, error) {
	var (
		err        error
		mqConsumer *mq.MQConsumer
	)

	authConsumerOnce.Do(func() {
		mqConsumer, err = mq.NewAuthTopicConsumer(opt)
		if err != nil {
			log.Error(err)
		}

		MQAuthTopicConsumerIns = mqConsumer.GetConsumer()
	})

	if MQAuthTopicConsumerIns == nil || err != nil {
		return nil, fmt.Errorf("failed to get mq consumer factory, mqFactory: %+v, error: %w", MQAuthTopicConsumerIns, err)
	}

	return MQAuthTopicConsumerIns, nil
}

func GetMQAnnotationTopicConsumerIns(opt *options.RocketMQOptions) (rocketmq.PushConsumer, error) {
	var (
		err        error
		mqConsumer *mq.MQConsumer
	)

	annotationConsumerOnce.Do(func() {
		mqConsumer, err = mq.NewAnnotationTopicConsumer(opt)
		if err != nil {
			log.Error(err)
		}

		MQAnnotationTopicConsumerIns = mqConsumer.GetConsumer()
	})

	if MQAnnotationTopicConsumerIns == nil || err != nil {
		return nil, fmt.Errorf("failed to get mq consumer factory, mqFactory: %+v, error: %w", MQAnnotationTopicConsumerIns, err)
	}
	return MQAnnotationTopicConsumerIns, nil
}
