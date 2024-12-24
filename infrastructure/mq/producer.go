package mq

import (
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/internal/infrastructure/options"
)

type MQProducer struct {
	mq rocketmq.Producer
}

func NewProducer(opt *options.RocketMQOptions) (*MQProducer, error) {
	rlog.SetLogLevel("warn")
	p, err := rocketmq.NewProducer(producer.WithNsResolver(
		primitive.NewPassthroughResolver(opt.Endpoint)),
		producer.WithRetry(opt.ProducerRetry),
		producer.WithCredentials(primitive.Credentials{
			SecretKey: opt.SecretKey,
			AccessKey: opt.AccessKey,
		}),
		producer.WithSendMsgTimeout(5*time.Second),
		producer.WithNamespace(opt.Namespace))

	if err != nil {
		log.Infof("init producer error: %+v", err.Error())
		return nil, err
	}

	err = p.Start()

	if err != nil {
		log.Infof("start producer error: %+v", err.Error())
		return nil, err
	}

	return &MQProducer{
		mq: p,
	}, nil
}

func (mq *MQProducer) Shutdown() {
	err := mq.mq.Shutdown()
	log.Info("shutdown producer success")
	log.Errorf("shutdown producer error: %+v", err.Error())
}

func (mq *MQProducer) GetProducer() rocketmq.Producer {
	return mq.mq
}
