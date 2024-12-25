package repository

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/po_entity"
	"gorm.io/gorm"
)

type DatasetRepo interface {
	CreateDatasetCollectionBinding(ctx context.Context, collectionBinding *po_entity.DatasetCollectionBinding, tb *gorm.DB) (*po_entity.DatasetCollectionBinding, error)
	CreateProviderHashEmbedding(ctx context.Context, embedding *po_entity.Embedding, tx *gorm.DB) (*po_entity.Embedding, error)
	BatchCreateProviderHashEmbedding(ctx context.Context, embeddings []*po_entity.Embedding, tx *gorm.DB) ([]*po_entity.Embedding, error)
	GetDatasetCollectionBinding(ctx context.Context, providerName string, modelName string, collectionType string, tb *gorm.DB) (*po_entity.DatasetCollectionBinding, error)
	GetProviderHashEmbedding(ctx context.Context, model string, hash string, provider string) (*po_entity.Embedding, error)
}
