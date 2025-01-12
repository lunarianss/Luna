package app_feature

import (
	"context"
	"errors"

	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/internal/api-server/core/rag/vector_db"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	po_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	datasetDomain "github.com/lunarianss/Luna/internal/api-server/domain/dataset/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/biz_entity"
	po_entity_dataset "github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/po_entity"

	providerDomain "github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type annotationReplyFeature struct {
	chatDomain     *chatDomain.ChatDomain
	datasetDomain  *datasetDomain.DatasetDomain
	redis          *redis.Client
	providerDomain *providerDomain.ProviderDomain
}

func NewAnnotationReplyFeature(chatDomain *chatDomain.ChatDomain, datasetDomain *datasetDomain.DatasetDomain, providerDomain *providerDomain.ProviderDomain, redis *redis.Client) *annotationReplyFeature {

	return &annotationReplyFeature{
		chatDomain:     chatDomain,
		datasetDomain:  datasetDomain,
		providerDomain: providerDomain,
		redis:          redis,
	}
}

func (arf *annotationReplyFeature) Query(ctx context.Context, app *po_entity.App, message *po_entity_chat.Message, query, accountID, invokeFrom string) (*po_entity_chat.MessageAnnotation, error) {
	annotationSetting, err := arf.chatDomain.AnnotationRepo.GetAnnotationSetting(ctx, app.ID, nil)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	collectionBinding := annotationSetting.CollectionBindingDetail

	scoreThreshold := annotationSetting.ScoreThreshold
	embeddingProviderName := collectionBinding.ProviderName
	embeddingModelName := collectionBinding.ModelName

	datasetCollectionBinding, err := arf.datasetDomain.DatasetRepo.GetDatasetCollectionBinding(ctx, embeddingProviderName, embeddingModelName, "annotation", nil)

	if err != nil {
		return nil, err
	}

	dataset := &po_entity_dataset.Dataset{
		ID:                     app.ID,
		TenantID:               app.TenantID,
		IndexingTechnique:      "high_quality",
		EmbeddingModelProvider: embeddingProviderName,
		EmbeddingModel:         embeddingModelName,
		CollectionBindingID:    datasetCollectionBinding.ID,
	}

	vector, err := vector_db.NewVector(ctx, dataset, []string{"doc_id", "annotation_id", "app_id"}, biz_entity.WEAVIATE, arf.redis, arf.providerDomain, nil, arf.datasetDomain, accountID)

	if err != nil {
		return nil, err
	}

	hitDocuments, err := vector.SearchByVector(ctx, query, 1, scoreThreshold)

	if err != nil {
		return nil, err
	}

	if len(hitDocuments) != 0 {
		annotationID := hitDocuments[0].Metadata["annotation_id"]
		score := hitDocuments[0].Score

		annotation, err := arf.chatDomain.AnnotationRepo.GetAnnotationByID(ctx, annotationID)

		if err != nil {
			return nil, err
		}

		fromResource := ""

		if invokeFrom == "service-api" || invokeFrom == "web-app" {
			fromResource = "api"
		} else {
			fromResource = "console"
		}

		annotationHistory := &po_entity_chat.AppAnnotationHitHistory{
			AnnotationID:       annotationID,
			AppID:              app.ID,
			AccountID:          accountID,
			Question:           query,
			Source:             fromResource,
			Score:              score,
			MessageID:          message.ID,
			AnnotationQuestion: annotation.Question,
			AnnotationContent:  annotation.Content,
		}

		_, err = arf.chatDomain.AnnotationRepo.CreateAppAnnotationHistory(ctx, annotationHistory)

		if err != nil {
			return nil, err
		}
		return annotation, nil
	}

	log.Infof("Query Similarity Vector %+v", hitDocuments)

	return nil, nil
}
