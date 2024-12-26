package event

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/infrastructure/shutdown"
	"github.com/lunarianss/Luna/internal/api-server/core/rag/vector_db"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	datasetDomain "github.com/lunarianss/Luna/internal/api-server/domain/dataset/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	repo_impl "github.com/lunarianss/Luna/internal/api-server/repository"
	"github.com/lunarianss/Luna/internal/infrastructure/mq"
	"github.com/lunarianss/Luna/internal/infrastructure/mysql"
	"github.com/lunarianss/Luna/internal/infrastructure/redis"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
	redisV9 "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type EnableAnnotationReplyTask struct {
	JobID                 string  `json:"job_id"`
	AppID                 string  `json:"app_id"`
	AccountID             string  `json:"account_id"`
	TenantID              string  `json:"tenant_id"`
	ScoreThreshold        float64 `json:"score_threshold"`
	EmbeddingProviderName string  `json:"embedding_provider_name"`
	EmbeddingModelName    string  `json:"embedding_model_name"`
}

type AnnotationEvent struct {
	mq rocketmq.PushConsumer
}

func (ae *AnnotationEvent) GetModule() string {
	return "mq_consumer_annotation_event"
}

func (ae *AnnotationEvent) HandleEnableAnnotationSuccess(ctx context.Context, tx *gorm.DB, redisIns *redisV9.Client, jobKey string, errorKey string, appKey string) error {
	if err := redisIns.SetEx(ctx, jobKey, "completed", 600*time.Second).Err(); err != nil {
		ae.HandleEnableAnnotationError(ctx, tx, redisIns, jobKey, errorKey, appKey, err)
		return err
	}

	if err := redisIns.Del(ctx, appKey).Err(); err != nil {
		log.Errorf("delete appAnnotationKey %s in redis: %s", appKey, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		ae.HandleEnableAnnotationError(ctx, tx, redisIns, jobKey, errorKey, appKey, err)
		return err
	}

	return nil
}

func (ae *AnnotationEvent) HandleEnableAnnotationError(ctx context.Context, tx *gorm.DB, redisIns *redisV9.Client, jobKey string, errorKey string, appKey string, err error) {
	log.Errorf("%#+v", err)
	if err := redisIns.SetEx(ctx, jobKey, "error", 600*time.Second).Err(); err != nil {
		log.Errorf("set jobKey %s in redis when rollback enable-annotation process, error: %s", jobKey, err)
	}

	if err := redisIns.SetEx(ctx, errorKey, err.Error(), 600*time.Second).Err(); err != nil {
		log.Errorf("set errorKey %s in redis when rollback enable-annotation process, error: %s", errorKey, err)
	}

	if err := tx.Rollback().Error; err != nil {
		log.Errorf("rollback error in enable annotation: %s", err.Error())
	}

	if err := redisIns.Del(ctx, appKey).Err(); err != nil {
		log.Errorf("delete appAnnotationKey %s in redis: %s", appKey, err.Error())
	}
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
			Expression: EnableAnnotationReplyTag,
		}, func(ctx context.Context, me ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

			util.LogCompleteInfo(me)

			for _, message := range me {
				var documents []*biz_entity.Document
				enableAnnotationBody := EnableAnnotationReplyTask{}
				if err := json.Unmarshal(message.Body, &enableAnnotationBody); err != nil {
					return consumer.ConsumeRetryLater, err
				}

				app, err := appDomain.AppRepo.GetTenantApp(c, enableAnnotationBody.AppID, enableAnnotationBody.TenantID)

				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						log.Errorf("app not fond when execute annotation job: %s", err.Error())
						continue
					} else {
						return consumer.ConsumeRetryLater, err
					}
				}

				messageAnnotations, err := chatDomainService.AnnotationRepo.FindAppAnnotations(ctx, app.ID)

				if err != nil {
					return consumer.ConsumeRetryLater, err
				}

				enableAppAnnotationKey := fmt.Sprintf("enable_app_annotation_%s", enableAnnotationBody.AppID)
				enableAppAnnotationJobKey := fmt.Sprintf("enable_app_annotation_job_%s", enableAnnotationBody.JobID)
				enableAppAnnotationErrorKey := fmt.Sprintf("enable_app_annotation_error_%s", enableAnnotationBody.JobID)

				tx := gormIns.Begin()

				datasetCollection, err := datasetDomain.DatasetRepo.GetDatasetCollectionBinding(ctx, enableAnnotationBody.EmbeddingProviderName, enableAnnotationBody.EmbeddingModelName, "annotation", tx)

				if err != nil {
					ae.HandleEnableAnnotationError(ctx, tx, redisIns, enableAppAnnotationJobKey, enableAppAnnotationErrorKey, enableAppAnnotationKey, err)
					return consumer.ConsumeRetryLater, err
				}

				_, err = chatDomainService.AnnotationRepo.GetAnnotationSettingWithCreate(ctx, enableAnnotationBody.AppID, enableAnnotationBody.ScoreThreshold, datasetCollection.ID, enableAnnotationBody.AccountID, tx)

				if err != nil {
					ae.HandleEnableAnnotationError(ctx, tx, redisIns, enableAppAnnotationJobKey, enableAppAnnotationErrorKey, enableAppAnnotationKey, err)
					return consumer.ConsumeRetryLater, err
				}

				dataset := &po_entity.Dataset{
					ID:                     enableAnnotationBody.AppID,
					TenantID:               enableAnnotationBody.TenantID,
					IndexingTechnique:      "high_quality",
					EmbeddingModelProvider: enableAnnotationBody.EmbeddingProviderName,
					EmbeddingModel:         enableAnnotationBody.EmbeddingModelName,
					CollectionBindingID:    datasetCollection.ID,
				}

				if len(messageAnnotations) == 0 {
					if err := ae.HandleEnableAnnotationSuccess(ctx, tx, redisIns, enableAppAnnotationJobKey, enableAppAnnotationErrorKey, enableAppAnnotationKey); err != nil {
						return consumer.ConsumeRetryLater, err
					} else {
						return consumer.ConsumeSuccess, nil
					}
				}

				for _, messageAnnotation := range messageAnnotations {
					documents = append(documents, &biz_entity.Document{
						PageContent: messageAnnotation.Question,
						Metadata: map[string]string{
							"annotation_id": messageAnnotation.ID,
							"app_id":        enableAnnotationBody.AppID,
							"doc_id":        messageAnnotation.ID,
						},
					})
				}

				vector, err := vector_db.NewVector(c, dataset, []string{"doc_id", "annotation_id", "app_id"}, biz_entity.WEAVIATE, redisIns, providerDomain, tx, datasetDomain, enableAnnotationBody.AccountID)

				if err != nil {
					ae.HandleEnableAnnotationError(ctx, tx, redisIns, enableAppAnnotationJobKey, enableAppAnnotationErrorKey, enableAppAnnotationKey, err)
					return consumer.ConsumeRetryLater, err
				}

				if err := vector.DeleteByMetadataField(ctx, "app_id", enableAnnotationBody.AppID); err != nil {
					ae.HandleEnableAnnotationError(ctx, tx, redisIns, enableAppAnnotationJobKey, enableAppAnnotationErrorKey, enableAppAnnotationKey, err)
					return consumer.ConsumeRetryLater, err
				}

				if err := vector.Create(ctx, documents); err != nil {
					ae.HandleEnableAnnotationError(ctx, tx, redisIns, enableAppAnnotationJobKey, enableAppAnnotationErrorKey, enableAppAnnotationKey, err)
					return consumer.ConsumeRetryLater, err
				}

				if err := ae.HandleEnableAnnotationSuccess(ctx, tx, redisIns, enableAppAnnotationJobKey, enableAppAnnotationErrorKey, enableAppAnnotationKey); err != nil {
					return consumer.ConsumeRetryLater, err
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
