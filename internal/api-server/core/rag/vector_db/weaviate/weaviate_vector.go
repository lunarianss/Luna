package weaviate_vector

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bsm/redislock"
	"github.com/go-openapi/strfmt"
	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/po_entity"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/redis/go-redis/v9"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
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

func (wv *WeaviateVector) SearchByVector(ctx context.Context, queryFloat []float32, topK int, scoreThreshold float32) ([]*biz_entity.Document, error) {

	var (
		documents   []*biz_entity.Document
		hitDocument []*biz_entity.Document
	)

	if topK == 0 {
		topK = 4
	}

	response, err := wv.client.GraphQL().Get().WithClassName(wv.collectionName).WithFields(
		graphql.Field{Name: "text"},
		graphql.Field{Name: "doc_id"},
		graphql.Field{Name: "app_id"},
		graphql.Field{Name: "annotation_id"},
		graphql.Field{Name: "_additional", Fields: []graphql.Field{
			{Name: "vector"},
			{Name: "distance"},
		}}).WithLimit(topK).WithNearVector(wv.client.GraphQL().NearVectorArgBuilder().
		WithVector(queryFloat)).Do(ctx)

	if err != nil {
		return nil, err
	}

	if len(response.Errors) > 0 {
		var errorMessage []string
		for _, errGraphQL := range response.Errors {
			errorMessage = append(errorMessage, errGraphQL.Message)
		}
		return nil, errors.WithCode(code.ErrVDBQueryError, "occurred error when similarity search %s", strings.Join(errorMessage, " | "))
	}

	collectionInfoMap, ok := response.Data["GET"].(map[string]any)

	if !ok {
		return nil, errors.WithSCode(code.ErrVDBConstructError, "data['GET'] is not map[string]any")
	}

	collectionInfoMapByte, err := json.Marshal(collectionInfoMap)

	if err != nil {
		return nil, errors.WithSCode(code.ErrEncodingJSON, err.Error())
	}
	var collectionInfo biz_entity.SimilaritySearchVDBResponse

	if err := json.Unmarshal(collectionInfoMapByte, &collectionInfo); err != nil {
		return nil, errors.WithSCode(code.ErrDecodingJSON, err.Error())
	}

	objectsInfo, ok := collectionInfo[wv.collectionName]

	if !ok {
		return nil, errors.WithSCode(code.ErrVDBConstructError, "data['collection_name'] is not  exist")
	}

	for _, objectInfo := range objectsInfo {
		score := 1 - objectInfo.Additional.Distance
		documents = append(documents, &biz_entity.Document{
			PageContent: objectInfo.Text,
			Vector:      objectInfo.Additional.Vector,
			Score:       score,
			Metadata: map[string]string{
				"annotation_id": objectInfo.AnnotationID,
				"app_id":        objectInfo.AppID,
				"doc_id":        objectInfo.AnnotationID,
			},
		})
	}

	for _, documentScore := range documents {
		if documentScore.Score > scoreThreshold {
			hitDocument = append(hitDocument, documentScore)
		}
	}

	sort.Slice(hitDocument, func(i, j int) bool {
		return hitDocument[i].Score > hitDocument[j].Score
	})
	return hitDocument, nil
}

func (wv WeaviateVector) ExistsCollection(ctx context.Context) (bool, error) {

	exist, err := wv.client.Schema().ClassExistenceChecker().WithClassName(wv.collectionName).Do(ctx)

	if err != nil {
		return true, errors.WithSCode(code.ErrVDB, err.Error())
	}

	return exist, nil

}
func (wv *WeaviateVector) DeleteByMetadataField(ctx context.Context, key string, value string) error {

	exist, err := wv.ExistsCollection(ctx)

	if err != nil {
		return err
	}

	if exist {
		_, err := wv.client.Batch().ObjectsBatchDeleter().WithClassName(wv.collectionName).WithOutput("minimal").WithWhere(filters.Where().WithPath([]string{key}).WithOperator(filters.Equal).WithValueText(value)).Do(ctx)
		if err != nil {
			return errors.WithSCode(code.ErrVDB, err.Error())
		}
	}
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
			return errors.WithSCode(code.ErrRedis, err.Error())
		}
	}

	defer lock.Release(ctx)

	collectionExistCacheKey := fmt.Sprintf("vector_indexing_%s", wv.collectionName)

	val, err := wv.redisIns.Get(ctx, collectionExistCacheKey).Result()

	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return errors.WithSCode(code.ErrRedis, err.Error())
		}
	}

	if val != "" && val == "1" {
		return nil
	}

	exist, err := wv.client.Schema().ClassExistenceChecker().WithClassName(wv.collectionName).Do(ctx)

	if err != nil {
		return errors.WithSCode(code.ErrVDB, err.Error())
	}

	if !exist {
		if err = wv.client.Schema().ClassCreator().WithClass(wv.defaultSchema(wv.collectionName)).Do(ctx); err != nil {
			return errors.WithSCode(code.ErrVDB, err.Error())
		}
	}

	if err := wv.redisIns.Set(ctx, collectionExistCacheKey, "1", time.Hour*1).Err(); err != nil {
		return errors.WithSCode(code.ErrRedis, err.Error())
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
		return errors.WithSCode(code.ErrVDB, err.Error())
	}

	return nil
}
