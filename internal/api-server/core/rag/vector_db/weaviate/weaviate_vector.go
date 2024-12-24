package weaviate_vector

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/go-openapi/strfmt"
	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/po_entity"
	"github.com/redis/go-redis/v9"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

type WeaviateVector struct {
	client         *weaviate.Client
	collectionName string
	redisIns       *redis.Client
}

var _ biz_entity.IVectorDB = (*WeaviateVector)(nil)

func NewWeaviateVector(dataset *po_entity.Dataset, attributes []string, client *weaviate.Client, redisIns *redis.Client) *WeaviateVector {
	var collectionName string

	if dataset.IndexStruct != "" {
		// todo 待补充
	} else {
		collectionName = dataset.GetCollectionNameByID(dataset.ID)
	}
	return &WeaviateVector{
		client:         client,
		collectionName: collectionName,
		redisIns:       redisIns,
	}
}

func (wv *WeaviateVector) GetType() string {
	return "weaviate"
}

func (wv *WeaviateVector) Create(ctx context.Context, texts []*biz_entity.Document, embeddings [][]float32, other any) error {
	if err := wv.CreateCollection(ctx); err != nil {
		return err
	}

	if err := wv.addTexts(ctx, texts, embeddings); err != nil {
		return err
	}
	return nil
}

func (wv *WeaviateVector) DeleteByMetadataField(ctx context.Context, key string, value string) error {
	return nil
}

func (wv *WeaviateVector) CreateCollection(ctx context.Context) error {

	lockName := fmt.Sprintf("vector_indexing_lock_%s", wv.collectionName)
	// Create a new lock client.
	locker := redislock.New(wv.redisIns)

	lock, err := locker.Obtain(ctx, lockName, 20*time.Second, nil)

	if err != nil {
		if errors.Is(err, redislock.ErrNotObtained) {
			return nil
		} else {
			return err
		}
	}

	defer lock.Release(ctx)

	collectionExistCacheKey := fmt.Sprintf("vector_indexing_%s", wv.collectionName)

	val, err := wv.redisIns.Get(ctx, collectionExistCacheKey).Result()

	if err != nil {
		return err
	}

	if val != "" && val == "1" {
		return nil
	}

	exist, err := wv.client.Schema().ClassExistenceChecker().WithClassName(wv.collectionName).Do(ctx)

	if err != nil {
		return err
	}

	if !exist {
		if err = wv.client.Schema().ClassCreator().WithClass(wv.defaultSchema(wv.collectionName)).Do(ctx); err != nil {
			return err
		}
	}

	if err := wv.redisIns.Set(ctx, collectionExistCacheKey, "1", time.Hour*1).Err(); err != nil {
		return err
	}
	return nil
}

func (wv *WeaviateVector) defaultSchema(indexName string) *models.Class {
	return &models.Class{
		Class: indexName,
		Properties: []*models.Property{
			{
				DataType: []string{"text"},
				Name:     "text",
			},
		},
	}
}

func (wv *WeaviateVector) addTexts(ctx context.Context, documents []*biz_entity.Document, embeddings [][]float32) error {

	var vectorObjects []*models.Object

	for i, document := range documents {
		dataProperties := map[string]string{
			"text": document.PageContent,
		}

		for mk, mv := range document.Metadata {
			dataProperties[mk] = mv
		}

		vectorObject := &models.Object{
			ID:         strfmt.UUID(document.Metadata["doc_id"]),
			Vector:     embeddings[i],
			Class:      wv.collectionName,
			Properties: dataProperties,
		}
		vectorObjects = append(vectorObjects, vectorObject)
	}

	batcher := wv.client.Batch().ObjectsBatcher()

	batcher.WithObjects(vectorObjects...)

	_, err := batcher.Do(ctx)

	if err != nil {
		return err
	}

	return nil
}
