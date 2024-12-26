package event_handler

import (
	"context"
	"encoding/json"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/internal/api-server/core/rag/vector_db"
	datasetDomain "github.com/lunarianss/Luna/internal/api-server/domain/dataset/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	"github.com/redis/go-redis/v9"
)

type AddAnnotationTask struct {
	AnnotationID        string `json:"annotation_id"`
	Question            string `json:"question"`
	TenantID            string `json:"tenant_id"`
	AppID               string `json:"app_id"`
	CollectionBindingID string `json:"collection_binding_id"`
	AccountID           string `json:"account_id"`
}

type AddAnnotationHandler struct {
	datasetDomain  *datasetDomain.DatasetDomain
	providerDomain *domain_service.ProviderDomain
	redisIns       *redis.Client
}

func NewAddAnnotationHandler(datasetDomain *datasetDomain.DatasetDomain, providerDomain *domain_service.ProviderDomain, redisIns *redis.Client) *AddAnnotationHandler {
	return &AddAnnotationHandler{
		datasetDomain:  datasetDomain,
		providerDomain: providerDomain,
		redisIns:       redisIns,
	}
}

var _ MQEventHandler = (*EnableAnnotationHandler)(nil)

func (eh *AddAnnotationHandler) Handle(ctx context.Context, message *primitive.MessageExt) (consumer.ConsumeResult, error) {
	var documents []*biz_entity.Document
	addAnnotationBody := AddAnnotationTask{}

	if err := json.Unmarshal(message.Body, &addAnnotationBody); err != nil {
		return consumer.ConsumeRetryLater, err
	}

	log.Infof("============ addAnnotationBody %+v ========", addAnnotationBody)

	collectionBinding, err := eh.datasetDomain.DatasetRepo.GetDatasetCollectionBindingByIDAndType(ctx, addAnnotationBody.CollectionBindingID, "annotation")

	if err != nil {
		return consumer.ConsumeRetryLater, err
	}

	dataset := &po_entity.Dataset{
		ID:                     addAnnotationBody.AppID,
		TenantID:               addAnnotationBody.TenantID,
		IndexingTechnique:      "high_quality",
		EmbeddingModelProvider: collectionBinding.ProviderName,
		EmbeddingModel:         collectionBinding.ModelName,
		CollectionBindingID:    collectionBinding.ID,
	}

	documents = append(documents, &biz_entity.Document{
		PageContent: addAnnotationBody.Question,
		Metadata: map[string]string{
			"annotation_id": addAnnotationBody.AnnotationID,
			"app_id":        addAnnotationBody.AppID,
			"doc_id":        addAnnotationBody.AnnotationID,
		},
	})

	vector, err := vector_db.NewVector(ctx, dataset, []string{"doc_id", "annotation_id", "app_id"}, biz_entity.WEAVIATE, eh.redisIns, eh.providerDomain, nil, eh.datasetDomain, addAnnotationBody.AccountID)

	if err != nil {
		log.Errorf("%#+v", err)
		return consumer.ConsumeRetryLater, err
	}

	if err := vector.Create(ctx, documents); err != nil {
		log.Errorf("%#+v", err)
		return consumer.ConsumeRetryLater, err
	}

	return consumer.ConsumeSuccess, nil
}
