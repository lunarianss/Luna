package biz_entity

import common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"

type SimpleModelProvider struct {
	Provider            string             `json:"provider"`
	Label               *common.I18nObject `json:"label"`
	IconSmall           *common.I18nObject `json:"icon_small"`
	IconLarge           *common.I18nObject `json:"icon_large"`
	SupportedModelTypes []common.ModelType `json:"supported_model_types"`
}

type ProviderModelWithStatus struct {
	Status ModelStatus `json:"status"`
	*common.ProviderModel
}
