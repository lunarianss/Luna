package etl_processor

import (
	"errors"

	interface_etl_processor "github.com/lunarianss/Luna/internal/api-server/core/rag/etl_processor/interface"
	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/biz_entity"
)

type etlProcessorFactory struct {
	indexType string
}

func NewEtlProcessorFactory(indexType string) *etlProcessorFactory {
	return &etlProcessorFactory{
		indexType: indexType,
	}
}

func (e *etlProcessorFactory) InitETLProcessor() (interface_etl_processor.IIndexProcessor, error) {
	if e.indexType == "" {
		return nil, errors.New("")
	}

	if e.indexType == string(biz_entity.PARAGRAPH_INDEX) {

	}

	return nil, nil
}
