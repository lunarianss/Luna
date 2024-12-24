package repository

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/po_entity"
	"gorm.io/gorm"
)

type DatasetRepo interface {
	CreateDatasetCollectionBinding(ctx context.Context, collectionBinding *po_entity.DatasetCollectionBinding, tb *gorm.DB) (*po_entity.DatasetCollectionBinding, error)
	GetDatasetCollectionBinding(ctx context.Context, providerName string, modelName string, collectionType string, tb *gorm.DB) (*po_entity.DatasetCollectionBinding, error)
}
