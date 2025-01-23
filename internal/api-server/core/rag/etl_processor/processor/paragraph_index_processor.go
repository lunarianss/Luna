package etl_processor_impl

import (
	"context"

	interface_etl_processor "github.com/lunarianss/Luna/internal/api-server/core/rag/etl_processor/interface"
	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/po_entity"
)

type paragraphProcessor struct {
}

func NewParagraphProcessor() interface_etl_processor.IIndexProcessor {
	return &paragraphProcessor{}
}

func (p *paragraphProcessor) Extract(ctx context.Context, setting *biz_entity.ExtractSetting, mode string) ([]*biz_entity.Document, error) {
	return nil, nil
}

func (p *paragraphProcessor) Transform(originDocuments []*biz_entity.Document) ([]*biz_entity.Document, error) {
	return nil, nil
}

func (p *paragraphProcessor) Load(dataset *po_entity.Dataset, documents []*biz_entity.Document, withKeyword bool) error {
	return nil
}

func (p *paragraphProcessor) Clean(dataset *po_entity.Dataset, nodesID []string, withKeyword bool) error {
	return nil
}
