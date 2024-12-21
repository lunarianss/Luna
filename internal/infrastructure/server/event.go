package server

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/infrastructure/shutdown"
)

type IMQEventConsumer interface {
	Subscribe(ctx context.Context, sd *shutdown.GracefulShutdown) error
	GetModule() string
}
type MQEventConsumer struct {
}

var mQEventConsumers []IMQEventConsumer

func RegisterConsumer(c IMQEventConsumer) {
	mQEventConsumers = append(mQEventConsumers, c)
}

func (s *BaseApiServer) InitMQConsumer(c context.Context, sd *shutdown.GracefulShutdown) error {
	for _, mQEventConsumer := range mQEventConsumers {
		if err := mQEventConsumer.Subscribe(c, sd); err != nil {
			return errors.WithMessage(err, fmt.Sprintf("mq consumer module %s error", mQEventConsumer.GetModule()))
		}
		log.Info(color.GreenString("MQConsumer %s init successfully.", mQEventConsumer.GetModule()))
	}
	return nil
}
