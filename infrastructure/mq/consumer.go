package mq

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/internal/infrastructure/options"
)

type MQConsumer struct {
	mqc rocketmq.PushConsumer
}

func NewConsumer(opt *options.RocketMQOptions) (*MQConsumer, error) {
	rlog.SetLogLevel("warn")

	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(opt.GroupName),
		consumer.WithNsResolver(primitive.NewPassthroughResolver(opt.Endpoint)),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithMaxReconsumeTimes(int32(opt.ConsumerRetry)),
		consumer.WithCredentials(primitive.Credentials{
			SecretKey: opt.SecretKey,
			AccessKey: opt.AccessKey,
		}),
		consumer.WithNamespace(opt.Namespace),
	)

	if err != nil {
		log.Infof("init producer error: %+v", err.Error())
		return nil, err
	}

	return &MQConsumer{
		mqc: c,
	}, nil

}

func (mq *MQConsumer) Shutdown() {
	err := mq.mqc.Shutdown()
	log.Info("shutdown consumer success")
	log.Errorf("shutdown consumer error: %+v", err.Error())
}

func (mq *MQConsumer) GetConsumer() rocketmq.PushConsumer {
	return mq.mqc
}
