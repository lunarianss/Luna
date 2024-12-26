package cache_embedding

import (
	"context"
	"errors"
	"math"
	"slices"

	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/po_entity"
	common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_configuration"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime/model_registry"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
	"gorm.io/gorm"
)

type ICacheEmbedding interface {
	EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error)
}

type cacheEmbedding struct {
	modelAllIntegrate *biz_entity.ModelIntegratedInstance
	user              string
	datasetDomain     *domain_service.DatasetDomain
	tx                *gorm.DB
}

func NewCacheEmbedding(modelAllIntegrate *biz_entity.ModelIntegratedInstance, user string, datasetDomain *domain_service.DatasetDomain, tx *gorm.DB) *cacheEmbedding {
	return &cacheEmbedding{
		modelAllIntegrate: modelAllIntegrate,
		user:              user,
		tx:                tx,
		datasetDomain:     datasetDomain,
	}
}

func (ce *cacheEmbedding) EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error) {
	var textEmbeddings = make([][]float32, len(texts))
	var (
		maxChunks int
		// 需要向量化的数组索引
		embeddingQueueIndices []int
		// 存放 embedding 之后的向量
		embeddingQueueEmbeddingsCache [][]float32
		hashCacheEmbeddings           []string
		embeddingQueueTexts           []string
		batchCreatedEmbedding         []*po_entity.Embedding
	)

	for i, text := range texts {
		textHash := util.GenerateTextHash(text)

		embedding, err := ce.datasetDomain.DatasetRepo.GetProviderHashEmbedding(ctx, ce.modelAllIntegrate.Model, textHash, ce.modelAllIntegrate.Provider)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				embeddingQueueIndices = append(embeddingQueueIndices, i)
				continue
			} else {
				return nil, err
			}
		}

		decodeEmbedding, err := embedding.GetEmbeddings()
		if err != nil {
			return nil, err
		}
		// 数据库已经有的，直接占位填充
		textEmbeddings[i] = decodeEmbedding
	}

	// 处理之前数据库仍有没有添加过的向量
	if len(embeddingQueueIndices) > 0 {
		caller := model_registry.NewModelRegisterCaller(ce.modelAllIntegrate.Model, string(common.TEXT_EMBEDDING), ce.modelAllIntegrate.Provider, ce.modelAllIntegrate.Credentials, ce.modelAllIntegrate.ModelTypeInstance)

		for _, embeddingIndex := range embeddingQueueIndices {
			embeddingQueueTexts = append(embeddingQueueTexts, texts[embeddingIndex])
		}

		modelSchema, err := ce.modelAllIntegrate.ModelTypeInstance.GetModelSchema(ce.modelAllIntegrate.Model, ce.modelAllIntegrate.Credentials)

		if err != nil {
			return nil, err
		}

		maxChunksAny := modelSchema.ModelProperties[common.MAX_CHUNKS]

		if maxChunksAny == nil {
			maxChunks = 1
		} else {
			maxChunks = maxChunksAny.(int)
		}

		for i := 0; i < len(embeddingQueueTexts); i += maxChunks {
			batchTexts := embeddingQueueTexts[i:int(math.Min(float64(i+maxChunks), float64(len(embeddingQueueTexts))))]

			embeddingResults, err := caller.InvokeTextEmbedding(ctx, nil, ce.user, "document", batchTexts)

			if err != nil {
				return nil, err
			}

			//对 embedding 后的向量做数据库缓存
			for _, vector := range embeddingResults.Embeddings {
				normalizedEmbedding, err := util.NormalizeVector(vector)
				if err != nil {
					log.Errorf("occur error: %s when normalize vector", err)
				}
				embeddingQueueEmbeddingsCache = append(embeddingQueueEmbeddingsCache, normalizedEmbedding)
			}

			for i, embeddingIndex := range embeddingQueueIndices {
				hash := util.GenerateTextHash(texts[embeddingIndex])
				embeddingObject := embeddingQueueEmbeddingsCache[i]
				textEmbeddings[embeddingIndex] = embeddingObject

				if err != nil {
					log.Errorf("occur error: %s when marshal embeddingObject %v", err, embeddingObject)
					continue
				}

				if !slices.Contains(hashCacheEmbeddings, hash) {
					createEmbedding := &po_entity.Embedding{
						ModelName:    ce.modelAllIntegrate.Model,
						ProviderName: ce.modelAllIntegrate.Provider,
						Hash:         hash,
					}

					createEmbedding.Embedding, err = createEmbedding.SetEmbedding(embeddingObject)

					if err != nil {
						log.Errorf("occurred error: %s when convert embedding to gob")
						continue
					}
					hashCacheEmbeddings = append(hashCacheEmbeddings, hash)
					batchCreatedEmbedding = append(batchCreatedEmbedding, createEmbedding)
				}
			}

			if batchCreated, err := ce.datasetDomain.DatasetRepo.BatchCreateProviderHashEmbedding(ctx, batchCreatedEmbedding, ce.tx); err != nil {
				util.LogCompleteInfo(batchCreated)
				return nil, err
			}
		}
	}

	log.Info("Text embedding result by AI and DB cache:  %v", textEmbeddings)
	return textEmbeddings, nil
}
