package po_entity

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"

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

func (a *Dataset) GetCollectionNameByID(id string) string {
	id = strings.ReplaceAll(id, "-", "_")
	return fmt.Sprintf("Vector_index_%s_Node", id)
}

type Embedding struct {
	ID           string `json:"id" gorm:"column:id"`
	ModelName    string `json:"model_name" gorm:"column:model_name"`
	Hash         string `json:"hash" gorm:"column:hash"`
	Embedding    []byte `json:"embedding" gorm:"column:embedding"`
	CreatedAt    int64  `json:"created_at" gorm:"column:created_at"`
	ProviderName string `json:"provider_name" gorm:"column:provider_name"`
}

func (a *Embedding) TableName() string {
	return "embeddings"
}

func (a *Embedding) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.NewString()
	return
}

func (a *Embedding) GetEmbeddings() ([]float32, error) {
	buf := bytes.NewBuffer(a.Embedding)
	decoder := gob.NewDecoder(buf)
	var embeddingData []float32

	if err := decoder.Decode(&embeddingData); err != nil {
		return nil, err
	}

	return embeddingData, nil
}

func (a *Embedding) SetEmbedding(embeddingVector []float32) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	err := encoder.Encode(embeddingVector)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type DatasetProcessRule struct {
	ID        string `json:"id" gorm:"column:id"`
	DatasetID string `json:"dataset_id" gorm:"column:dataset_id"`
	Mode      string `json:"mode" gorm:"column:mode;default:automatic"`
	Rules     string `json:"rules" gorm:"column:rules"`
	CreatedBy string `json:"created_by" gorm:"column:created_by"`
	CreatedAt int64  `json:"created_at" gorm:"column:created_at"`
}

func (a *DatasetProcessRule) TableName() string {
	return "dataset_process_rules"
}

type AppDatasetJoin struct {
	ID        string `json:"id" gorm:"column:id"`
	AppID     string `json:"app_id" gorm:"column:app_id"`
	DatasetID string `json:"dataset_id" gorm:"column:dataset_id"`
	CreatedAt int64  `json:"created_at" gorm:"column:created_at"`
}

func (a *AppDatasetJoin) TableName() string {
	return "app_dataset_joins"
}

type DatasetQuery struct {
	ID            string `json:"id" gorm:"column:id"`
	DatasetID     string `json:"dataset_id" gorm:"column:dataset_id"`
	Content       string `json:"content" gorm:"column:content"`
	Source        string `json:"source" gorm:"column:source"`
	SourceAppID   string `json:"source_app_id" gorm:"column:source_app_id"`
	CreatedByRole string `json:"created_by_role" gorm:"column:created_by_role"`
	CreatedBy     string `json:"created_by" gorm:"column:created_by"`
	CreatedAt     int64  `json:"created_at" gorm:"column:created_at"`
}

func (a *DatasetQuery) TableName() string {
	return "dataset_queries"
}

type DatasetKeywordTable struct {
	ID             string `json:"id" gorm:"column:id"`
	DatasetID      string `json:"dataset_id" gorm:"column:dataset_id"`
	KeywordTable   string `json:"keyword_table" gorm:"column:keyword_table"`
	DataSourceType string `json:"data_source_type" gorm:"column:data_source_type;default:database"`
}

func (a *DatasetKeywordTable) TableName() string {
	return "dataset_keyword_tables"
}
