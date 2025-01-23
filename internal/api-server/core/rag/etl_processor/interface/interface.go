package interface_etl_processor

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/po_entity"
)

type IIndexProcessor interface {
	Extract(ctx context.Context, setting *biz_entity.ExtractSetting, mode string) ([]*biz_entity.Document, error)
	Transform(originDocuments []*biz_entity.Document) ([]*biz_entity.Document, error)
	Load(dataset *po_entity.Dataset, documents []*biz_entity.Document, withKeyword bool) error
	Clean(dataset *po_entity.Dataset, nodesID []string, withKeyword bool) error
}
