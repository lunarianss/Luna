package event

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/infrastructure/shutdown"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	datasetDomain "github.com/lunarianss/Luna/internal/api-server/domain/dataset/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/event/event_handler"
	repo_impl "github.com/lunarianss/Luna/internal/api-server/repository"
	"github.com/lunarianss/Luna/internal/infrastructure/mq"
	"github.com/lunarianss/Luna/internal/infrastructure/mysql"
	"github.com/lunarianss/Luna/internal/infrastructure/redis"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
)

type AnnotationEvent struct {
	mq rocketmq.PushConsumer
}

func (ae *AnnotationEvent) GetModule() string {
	return "mq_consumer_annotation_event"
}

func (ae *AnnotationEvent) Subscribe(c context.Context, sd *shutdown.GracefulShutdown) error {

	gormIns, err := mysql.GetMySQLIns(nil)

	if err != nil {
		return err
	}

	redisIns, err := redis.GetRedisIns(nil)

	if err != nil {
		return err
	}

	sig := make(chan struct{})

	sd.AddShutdownCallback(shutdown.ShutdownFunc(func(s string) error {
		sig <- struct{}{}
		return nil
	}))

	// repos
	tenantRepo := repo_impl.NewTenantRepoImpl(gormIns)
	appRepo := repo_impl.NewAppRepoImpl(gormIns)
	messageRepo := repo_impl.NewMessageRepoImpl(gormIns)
	providerRepo := repo_impl.NewProviderRepoImpl(gormIns)
	webAppRepo := repo_impl.NewWebAppRepoImpl(gormIns)
	modelProviderRepo := repo_impl.NewModelProviderRepoImpl(gormIns)
	providerConfigurationsManager := domain_service.NewProviderConfigurationsManager(providerRepo, modelProviderRepo, "", nil)
	annotationRepo := repo_impl.NewAnnotationRepoImpl(gormIns)
	// domain
	providerDomain := domain_service.NewProviderDomain(providerRepo, modelProviderRepo, tenantRepo, providerConfigurationsManager)
	appDomain := appDomain.NewAppDomain(appRepo, webAppRepo, gormIns)
	chatDomainService := chatDomain.NewChatDomain(messageRepo, annotationRepo)
	datasetRepo := repo_impl.NewDatasetRepoImpl(gormIns)

	datasetDomain := datasetDomain.NewDatasetDomain(datasetRepo)

	mqConsumer, err := mq.GetMQAnnotationTopicConsumerIns(nil)

	ae.mq = mqConsumer

	if err != nil {
		return err
	}

	go func() {
		mqConsumer.Subscribe(AnnotationTopic, consumer.MessageSelector{
			Type:       consumer.TAG,
			Expression: EnableAnnotationReplyTag + "||" + AddAnnotationTag,
		}, func(ctx context.Context, me ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, message := range me {
				util.LogCompleteInfo(me)
				if message.GetTags() == EnableAnnotationReplyTag {
					return event_handler.NewEnableAnnotationHandler(appDomain, chatDomainService, datasetDomain, redisIns, gormIns, providerDomain).Handle(ctx, message)
				}

				if message.GetTags() == AddAnnotationTag {
					return event_handler.NewAddAnnotationHandler(datasetDomain, providerDomain, redisIns).Handle(ctx, message)
				}
			}
			return consumer.ConsumeSuccess, nil
		})

		err = mqConsumer.Start()
		if err != nil {
			log.Infof("start annotation consumer error: %+v", err.Error())
		}

		<-sig
		log.Infof("annotation consumer %s exit successfully", ae.GetModule())
	}()

	return nil

}
