package provider

import (
	"github.com/lunarianss/Hurricane/internal/api-server/model/v1"
	modelRuntimeEntities "github.com/lunarianss/Hurricane/internal/api-server/model-runtime/entities"
)

type ProviderConfiguration struct {
	TenantId             string                              `json:"tenant_id"`
	Provider             modelRuntimeEntities.ProviderEntity `json:"provider"`
	PreferedProviderType model.ProviderType                  `json:"preferred_provider_type"`
	UsingProviderType    model.ProviderType                  `json:"using_provider_type"`
}
