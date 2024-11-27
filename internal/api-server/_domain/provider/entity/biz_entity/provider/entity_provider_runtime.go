package biz_entity

import (
	provider_model_entity "github.com/lunarianss/Luna/internal/api-server/_domain/provider/entity/biz_entity/provider/provider_model"
)

type ModelProviderRuntime struct {
	ProviderSchema   *ProviderStaticConfiguration
	ModelConfPath    string
	ModelInstanceMap map[string]*provider_model_entity.AIModel
}
