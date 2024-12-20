package po_entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DatasetCollectionBinding struct {
	ID             string `json:"id" gorm:"column:id"`
	ProviderName   string `json:"provider_name" gorm:"column:provider_name"`
	ModelName      string `json:"model_name" gorm:"column:model_name"`
	Type           string `json:"type" gorm:"column:type"`
	CollectionName string `json:"collection_name" gorm:"column:collection_name"`
	CreatedAt      int64  `json:"created_at" gorm:"column:created_at"`
}

func (a *DatasetCollectionBinding) TableName() string {
	return "dataset_collection_bindings"
}

func (a *DatasetCollectionBinding) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.NewString()
	return
}

type Dataset struct {
	ID                     string `json:"id" gorm:"column:id"`
	TenantID               string `json:"tenant_id" gorm:"column:tenant_id"`
	Name                   string `json:"name" gorm:"column:name"`
	Description            string `json:"description" gorm:"column:description"`
	Provider               string `json:"provider" gorm:"column:provider"`
	Permission             string `json:"permission" gorm:"column:permission"`
	DataSourceType         string `json:"data_source_type" gorm:"column:data_source_type"`
	IndexingTechnique      string `json:"indexing_technique" gorm:"column:indexing_technique"`
	IndexStruct            string `json:"index_struct" gorm:"column:index_struct"`
	CreatedBy              string `json:"created_by" gorm:"column:created_by"`
	CreatedAt              int64  `json:"created_at" gorm:"column:created_at"`
	UpdatedBy              string `json:"updated_by" gorm:"column:updated_by"`
	UpdatedAt              int64  `json:"updated_at" gorm:"column:updated_at"`
	EmbeddingModel         string `json:"embedding_model" gorm:"column:embedding_model"`
	EmbeddingModelProvider string `json:"embedding_model_provider" gorm:"column:embedding_model_provider"`
	CollectionBindingID    string `json:"collection_binding_id" gorm:"column:collection_binding_id"`
	RetrievalModel         string `json:"retrieval_model" gorm:"column:retrieval_model"`
}

func (a *Dataset) TableName() string {
	return "datasets"
}

func (a *Dataset) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.NewString()
	return
}
