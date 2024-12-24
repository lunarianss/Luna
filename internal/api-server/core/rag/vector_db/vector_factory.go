package vector_db

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/core/rag/cache_embedding"
	weaviate_vector "github.com/lunarianss/Luna/internal/api-server/core/rag/vector_db/weaviate"
	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"
	"github.com/lunarianss/Luna/internal/infrastructure/weaviate"
	"github.com/redis/go-redis/v9"
)

type Vector struct {
	Attributes      []string
	Dataset         *po_entity.Dataset
	VectorProcessor biz_entity.IVectorDB
	Embeddings      biz_entity.IEmbeddings
	ProviderDomain  *domain_service.ProviderDomain
	Vdb             biz_entity.VdbType
	redis           *redis.Client
}

func NewVector(ctx context.Context, dataset *po_entity.Dataset, attributes []string, vdbName biz_entity.VdbType, redis *redis.Client, providerDomain *domain_service.ProviderDomain) (*Vector, error) {
	var (
		err error
	)

	vector := &Vector{
		ProviderDomain: providerDomain,
		Dataset:        dataset,
		Vdb:            vdbName,
		redis:          redis,
	}

	vector.Embeddings, err = vector.GetEmbeddings(ctx)

	if err != nil {
		return nil, err
	}

	vector.VectorProcessor, err = vector.InitVectorProcessor(ctx)

	if err != nil {
		return nil, err
	}
	return vector, nil
}

func (v *Vector) InitVectorProcessor(ctx context.Context) (biz_entity.IVectorDB, error) {

	if v.Vdb == biz_entity.WEAVIATE {
		weaviateClient, err := weaviate.GetWeaviateClient(nil)
		if err != nil {
			return nil, err
		}
		vectorProcessor := weaviate_vector.NewWeaviateVector(v.Dataset, v.Attributes, weaviateClient, v.redis)

		return vectorProcessor, nil
	}
	return nil, nil
}

func (v *Vector) Create(ctx context.Context, texts []*biz_entity.Document) error {
	if len(texts) > 0 {
		var textStr []string

		for _, text := range texts {
			textStr = append(textStr, text.PageContent)
		}

		embeddings, err := v.Embeddings.EmbedDocuments(ctx, textStr)
		if err != nil {
			return err
		}

		if err := v.VectorProcessor.Create(ctx, texts, embeddings, nil); err != nil {
			return err
		}
	}
	return nil
}

func (v *Vector) DeleteByMetadataField(ctx context.Context, key string, value string) error {
	return nil
}

func (v *Vector) GetEmbeddings(ctx context.Context) (cache_embedding.ICacheEmbedding, error) {
	embeddingModel, err := v.ProviderDomain.GetModelInstance(ctx, v.Dataset.TenantID, v.Dataset.EmbeddingModelProvider, v.Dataset.EmbeddingModel, common.TEXT_EMBEDDING)

	if err != nil {
		return nil, err
	}
	return cache_embedding.NewCacheEmbedding(embeddingModel, v.Dataset.UpdatedBy), nil
}
