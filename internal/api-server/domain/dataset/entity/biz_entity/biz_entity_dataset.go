package biz_entity

import "context"

type VdbType string

var (
	WEAVIATE VdbType = "weaviate"
)

type IVectorDB interface {
	GetType() string
	Create(ctx context.Context, texts []*Document, embeddings [][]float64, other any) error
	DeleteByMetadataField(ctx context.Context, key string, value string) error
}

type IEmbeddings interface {
	EmbedDocuments(ctx context.Context, texts []*Document) ([][]float64, error)
}

type Document struct {
	PageContent string
	Vector      []float64
	Metadata    map[string]string
	Provider    string
}
