package repo_impl

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/repository"
	"gorm.io/gorm"
)

type DatasetRepoImpl struct {
	db *gorm.DB
}

var _ repository.DatasetRepo = (*DatasetRepoImpl)(nil)

func NewDatasetRepoImpl(db *gorm.DB) repository.DatasetRepo {
	return &DatasetRepoImpl{
		db: db,
	}
}

func (dr *DatasetRepoImpl) CreateDatasetCollectionBinding(ctx context.Context, collectionBinding *po_entity.DatasetCollectionBinding, tx *gorm.DB) (*po_entity.DatasetCollectionBinding, error) {

	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = dr.db
	}

	if err := dbIns.Create(collectionBinding).Error; err != nil {
		return nil, err
	}

	return collectionBinding, nil

}

func (dr *DatasetRepoImpl) GetProviderHashEmbedding(ctx context.Context, model string, hash string, provider string) (*po_entity.Embedding, error) {

	var embedding po_entity.Embedding

	if err := dr.db.Where("model_name = ? AND hash = ? AND provider_name = ?", model, hash, provider).First(&embedding).Error; err != nil {
		return nil, err
	}

	return &embedding, nil
}

func (dr *DatasetRepoImpl) GetDatasetCollectionBinding(ctx context.Context, providerName string, modelName string, collectionType string, tx *gorm.DB) (*po_entity.DatasetCollectionBinding, error) {

	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = dr.db
	}

	var datasetCollection *po_entity.DatasetCollectionBinding

	err := dbIns.Where("provider_name = ? AND model_name = ? AND type = ?", providerName, modelName, collectionType).First(&datasetCollection).Error

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	dataset := &po_entity.Dataset{}

	if datasetCollection.ID == "" {
		datasetCollection = &po_entity.DatasetCollectionBinding{
			ProviderName:   providerName,
			ModelName:      modelName,
			Type:           collectionType,
			CollectionName: dataset.GetCollectionNameByID(uuid.NewString()),
		}

		datasetCollection, err = dr.CreateDatasetCollectionBinding(ctx, datasetCollection, tx)

		if err != nil {
			return nil, err
		}
	}

	return datasetCollection, nil
}

func (dr *DatasetRepoImpl) CreateProviderHashEmbedding(ctx context.Context, embedding *po_entity.Embedding, tx *gorm.DB) (*po_entity.Embedding, error) {
	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = dr.db
	}

	if err := dbIns.Create(embedding).Error; err != nil {
		return nil, err
	}

	return embedding, nil
}

func (dr *DatasetRepoImpl) BatchCreateProviderHashEmbedding(ctx context.Context, embeddings []*po_entity.Embedding, tx *gorm.DB) ([]*po_entity.Embedding, error) {

	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = dr.db
	}

	if err := dbIns.CreateInBatches(&embeddings, 100).Error; err != nil {
		return nil, err
	}
	return embeddings, nil
}
