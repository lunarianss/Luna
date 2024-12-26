package event_handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/internal/api-server/core/rag/vector_db"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	datasetDomain "github.com/lunarianss/Luna/internal/api-server/domain/dataset/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type EnableAnnotationReplyTask struct {
	JobID                 string  `json:"job_id"`
	AppID                 string  `json:"app_id"`
	AccountID             string  `json:"account_id"`
	TenantID              string  `json:"tenant_id"`
	ScoreThreshold        float32 `json:"score_threshold"`
	EmbeddingProviderName string  `json:"embedding_provider_name"`
	EmbeddingModelName    string  `json:"embedding_model_name"`
}

type EnableAnnotationHandler struct {
	appDomain      *appDomain.AppDomain
	chatDomain     *chatDomain.ChatDomain
	datasetDomain  *datasetDomain.DatasetDomain
	providerDomain *domain_service.ProviderDomain
	redisIns       *redis.Client
	gormIns        *gorm.DB
}

var _ MQEventHandler = (*EnableAnnotationHandler)(nil)

func NewEnableAnnotationHandler(appDomain *appDomain.AppDomain, chatDomain *chatDomain.ChatDomain, datasetDomain *datasetDomain.DatasetDomain, redisIns *redis.Client,
	gormIns *gorm.DB, providerDomain *domain_service.ProviderDomain) *EnableAnnotationHandler {
	return &EnableAnnotationHandler{
		appDomain:      appDomain,
		chatDomain:     chatDomain,
		datasetDomain:  datasetDomain,
		redisIns:       redisIns,
		gormIns:        gormIns,
		providerDomain: providerDomain,
	}
}

func (eh *EnableAnnotationHandler) Handle(ctx context.Context, message *primitive.MessageExt) (consumer.ConsumeResult, error) {
	var documents []*biz_entity.Document
	enableAnnotationBody := EnableAnnotationReplyTask{}

	if err := json.Unmarshal(message.Body, &enableAnnotationBody); err != nil {
		return consumer.ConsumeRetryLater, err
	}

	log.Infof("============ enableAnnotationBody %+v ========", enableAnnotationBody)

	app, err := eh.appDomain.AppRepo.GetTenantApp(ctx, enableAnnotationBody.AppID, enableAnnotationBody.TenantID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("app not fond when execute annotation job: %s", err.Error())
			return consumer.ConsumeSuccess, err
		} else {
			return consumer.ConsumeRetryLater, err
		}
	}

	messageAnnotations, err := eh.chatDomain.AnnotationRepo.FindAppAnnotations(ctx, app.ID)

	if err != nil {
		return consumer.ConsumeRetryLater, err
	}

	enableAppAnnotationKey := fmt.Sprintf("enable_app_annotation_%s", enableAnnotationBody.AppID)
	enableAppAnnotationJobKey := fmt.Sprintf("enable_app_annotation_job_%s", enableAnnotationBody.JobID)
	enableAppAnnotationErrorKey := fmt.Sprintf("enable_app_annotation_error_%s", enableAnnotationBody.JobID)

	tx := eh.gormIns.Begin()

	datasetCollection, err := eh.datasetDomain.DatasetRepo.GetDatasetCollectionBinding(ctx, enableAnnotationBody.EmbeddingProviderName, enableAnnotationBody.EmbeddingModelName, "annotation", tx)

	if err != nil {
		eh.HandleEnableAnnotationError(ctx, tx, eh.redisIns, enableAppAnnotationJobKey, enableAppAnnotationErrorKey, enableAppAnnotationKey, err)
		return consumer.ConsumeRetryLater, err
	}

	_, err = eh.chatDomain.AnnotationRepo.GetAnnotationSettingWithCreate(ctx, enableAnnotationBody.AppID, enableAnnotationBody.ScoreThreshold, datasetCollection.ID, enableAnnotationBody.AccountID, tx)

	if err != nil {
		eh.HandleEnableAnnotationError(ctx, tx, eh.redisIns, enableAppAnnotationJobKey, enableAppAnnotationErrorKey, enableAppAnnotationKey, err)
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
		if err := eh.HandleEnableAnnotationSuccess(ctx, tx, eh.redisIns, enableAppAnnotationJobKey, enableAppAnnotationErrorKey, enableAppAnnotationKey); err != nil {
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

	vector, err := vector_db.NewVector(ctx, dataset, []string{"doc_id", "annotation_id", "app_id"}, biz_entity.WEAVIATE, eh.redisIns, eh.providerDomain, tx, eh.datasetDomain, enableAnnotationBody.AccountID)

	if err != nil {
		eh.HandleEnableAnnotationError(ctx, tx, eh.redisIns, enableAppAnnotationJobKey, enableAppAnnotationErrorKey, enableAppAnnotationKey, err)
		return consumer.ConsumeRetryLater, err
	}

	if err := vector.DeleteByMetadataField(ctx, "app_id", enableAnnotationBody.AppID); err != nil {
		eh.HandleEnableAnnotationError(ctx, tx, eh.redisIns, enableAppAnnotationJobKey, enableAppAnnotationErrorKey, enableAppAnnotationKey, err)
		return consumer.ConsumeRetryLater, err
	}

	if err := vector.Create(ctx, documents); err != nil {
		eh.HandleEnableAnnotationError(ctx, tx, eh.redisIns, enableAppAnnotationJobKey, enableAppAnnotationErrorKey, enableAppAnnotationKey, err)
		return consumer.ConsumeRetryLater, err
	}

	if err := eh.HandleEnableAnnotationSuccess(ctx, tx, eh.redisIns, enableAppAnnotationJobKey, enableAppAnnotationErrorKey, enableAppAnnotationKey); err != nil {
		return consumer.ConsumeRetryLater, err
	}

	return consumer.ConsumeSuccess, nil
}

func (ae *EnableAnnotationHandler) HandleEnableAnnotationSuccess(ctx context.Context, tx *gorm.DB, redisIns *redis.Client, jobKey string, errorKey string, appKey string) error {
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

func (ae *EnableAnnotationHandler) HandleEnableAnnotationError(ctx context.Context, tx *gorm.DB, redisIns *redis.Client, jobKey string, errorKey string, appKey string, err error) {
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
