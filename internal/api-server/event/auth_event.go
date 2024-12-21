package event

import (
	"context"
	"encoding/json"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/infrastructure/shutdown"
	"github.com/lunarianss/Luna/internal/api-server/config"
	"github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	"github.com/lunarianss/Luna/internal/infrastructure/email"
	"github.com/lunarianss/Luna/internal/infrastructure/mq"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
)

type SendEmailCodeMessage struct {
	Language  string `json:"language"`
	Email     string `json:"email"`
	EmailCode string `json:"email_code"`
}

type A struct {
	Name string `json:"name"`
}

type AuthEvent struct {
	mq rocketmq.PushConsumer
}

func (ae *AuthEvent) GetModule() string {
	return "mq_consumer_auth_event"
}
func (ae *AuthEvent) Subscribe(c context.Context, sd *shutdown.GracefulShutdown) error {
	email, err := email.GetEmailSMTPIns(nil)

	if err != nil {
		return err
	}

	// config
	config, err := config.GetLunaRuntimeConfig()

	if err != nil {
		return err
	}
	sig := make(chan struct{})

	sd.AddShutdownCallback(shutdown.ShutdownFunc(func(s string) error {
		sig <- struct{}{}

		if ae.mq != nil {
			ae.mq.Shutdown()
		}

		return nil
	}))

	accountDomain := domain_service.NewAccountDomain(nil, nil, config, email, nil)
	mqConsumer, err := mq.GetMQConsumerIns(nil)

	ae.mq = mqConsumer

	if err != nil {
		return err
	}

	go func() {
		mqConsumer.Subscribe(AuthTopic, consumer.MessageSelector{
			Type:       consumer.TAG,
			Expression: SendEmailCodeTag,
		}, func(ctx context.Context, me ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

			util.LogCompleteInfo(me)

			for _, message := range me {
				sendMessageBody := SendEmailCodeMessage{}
				if err := json.Unmarshal(message.Body, &sendMessageBody); err != nil {
					return consumer.ConsumeRetryLater, err
				}

				accountDomain.SendEmailHtml(ctx, sendMessageBody.Language, sendMessageBody.Email, sendMessageBody.EmailCode)
			}

			return consumer.ConsumeSuccess, nil
		})

		err = mqConsumer.Start()

		if err != nil {
			log.Infof("start producer error: %+v", err.Error())
		}

		<-sig
		log.Infof("consumer %s exit", ae.GetModule())
	}()

	return nil

}
