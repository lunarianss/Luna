package entities

import (
	modelRuntimeEntities "github.com/lunarianss/Luna/internal/api-server/model-runtime/entities"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
)

type ProviderConfiguration struct {
	TenantId              int64                                `json:"tenant_id"`
	Provider              *modelRuntimeEntities.ProviderEntity `json:"provider"`
	PreferredProviderType model.ProviderType                   `json:"preferred_provider_type"`
	UsingProviderType     model.ProviderType                   `json:"using_provider_type"`
	SystemConfiguration   *SystemConfiguration                 `json:"system_configuration"`
	CustomConfiguration   *CustomConfiguration                 `json:"custom_configuration"`
	ModelSettings         *ModelSettings                       `json:"model_settings"`
}

type ProviderConfigurations struct {
	TenantId       int64                             `json:"tenant_id"`
	Configurations map[string]*ProviderConfiguration `json:"configurations"`
}
