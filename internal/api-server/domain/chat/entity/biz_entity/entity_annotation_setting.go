package biz_entity

import (
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	po_dataset "github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/po_entity"
)

type CollectionBingDetail struct {
	ID             string `json:"id"`
	ProviderName   string `json:"provider_name"`
	ModelName      string `json:"model_name"`
	Type           string `json:"type"`
	CollectionName string `json:"collection_name"`
	CreatedAt      int64  `json:"created_at"`
}

type AnnotationSettingWithBinding struct {
	ID                      string                `json:"id"`
	AppID                   string                `json:"app_id"`
	ScoreThreshold          float64               `json:"score_threshold"`
	CollectionBindingID     string                `json:"collection_binding_id"`
	CreatedUserID           string                `json:"created_user_id"`
	CreatedAt               int64                 `json:"created_at"`
	UpdatedUserID           string                `json:"updated_user_id"`
	UpdatedAt               int64                 `json:"updated_at"`
	CollectionBindingDetail *CollectionBingDetail `json:"collection_binding_detail"`
}

func ConvertPoAnnotationSetting(setting *po_entity.AppAnnotationSetting, binding *po_dataset.DatasetCollectionBinding) *AnnotationSettingWithBinding {
	convertedBinding := &CollectionBingDetail{
		ID:             binding.ID,
		ProviderName:   binding.ProviderName,
		ModelName:      binding.ModelName,
		Type:           binding.Type,
		CollectionName: binding.CollectionName,
		CreatedAt:      binding.CreatedAt,
	}

	return &AnnotationSettingWithBinding{
		ID:                      setting.ID,
		AppID:                   setting.AppID,
		ScoreThreshold:          setting.ScoreThreshold,
		CollectionBindingID:     setting.CollectionBindingID,
		CreatedUserID:           setting.CreatedUserID,
		CreatedAt:               setting.CreatedAt,
		UpdatedUserID:           setting.UpdatedUserID,
		UpdatedAt:               setting.UpdatedAt,
		CollectionBindingDetail: convertedBinding,
	}
}
