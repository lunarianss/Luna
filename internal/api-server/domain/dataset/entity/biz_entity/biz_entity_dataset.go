package biz_entity

import "context"

type VdbType string

var (
	WEAVIATE VdbType = "weaviate"
)

type IVectorDB interface {
	GetType() string
	Create(ctx context.Context, texts []*Document, embeddings [][]float32, other any) error
	DeleteByMetadataField(ctx context.Context, key string, value string) error
	SearchByVector(ctx context.Context, queryFloat []float32, topK int, scoreThreshold float32) ([]*Document, error)
}

type IEmbeddings interface {
	EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error)
	EmbedQuery(ctx context.Context, text string) ([]float32, error)
}

type Document struct {
	PageContent string
	Vector      []float32
	Metadata    map[string]string
	Provider    string
	Score       float32
}

type SimilaritySearchAdditional struct {
	Vector   []float32 `json:"vector"`
	Distance float32   `json:"distance"`
}
type SimilaritySearchCollectionInfo struct {
	Text         string                      `json:"text"`
	DocID        string                      `json:"doc_id"`
	AppID        string                      `json:"app_id"`
	AnnotationID string                      `json:"annotation_id"`
	Additional   *SimilaritySearchAdditional `json:"_additional"`
}

type SimilaritySearchVDBResponse map[string][]*SimilaritySearchCollectionInfo
