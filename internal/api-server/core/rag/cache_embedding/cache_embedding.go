package cache_embedding

import (
	"context"

	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_configuration"
)

type ICacheEmbedding interface {
	EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error)
}

type cacheEmbedding struct {
	modelAllIntegrate *biz_entity.ModelIntegratedInstance
	user              string
}

func NewCacheEmbedding(modelAllIntegrate *biz_entity.ModelIntegratedInstance, user string) *cacheEmbedding {
	return &cacheEmbedding{
		modelAllIntegrate: modelAllIntegrate,
		user:              user,
	}
}

func (ce *cacheEmbedding) EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error) {

	return nil, nil

}
