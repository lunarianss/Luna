package repo

import (
	"github.com/lunarianss/Hurricane/internal/apiServer/model/v1"
	"github.com/lunarianss/Hurricane/internal/apiServer/model_runtime/entities"
)

type ModelProviderRepo interface {
	// Get tenant's model providers
	GetTenantModelProviders(tenantId int64) ([]*model.Provider, error)
	// Get tenant's model providers mapped by provider name
	GetMapTenantModelProviders(tenantId int64) (map[string]*model.Provider, error)
	// Get all inner Providers
	GetSystemProviders() ([]*entities.ProviderEntity, error)
	// Get all inner Providers mapped by provider name
	GetMapSystemProviders() (map[string]*entities.ProviderEntity, error)
}
