package repo

import (
	"context"

	model "github.com/lunarianss/Luna/internal/api-server/model/v1"
)

type ModelRepo interface {
	// GetTenantModels get all models by searchModel
	GetTenantModel(ctx context.Context, tenantId, providerName, modelName, modelType string) (*model.ProviderModel, error)
	// Get the corresponding TenantDefaultModel record
	GetTenantDefaultModel(ctx context.Context, tenantId, modelType string) (*model.TenantDefaultModel, error)
	// CreateTenantDefaultModel create a default model for the tenant
	CreateTenantDefaultModel(ctx context.Context, tenantDefaultModel *model.TenantDefaultModel) (*model.TenantDefaultModel, error)
	// UpdateModel updates model
	UpdateModel(ctx context.Context, model *model.ProviderModel) error
	// CreateModel create model
	CreateModel(ctx context.Context, model *model.ProviderModel) error
}
