package repository

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/po_entity"
)

type ModelRepo interface {

	// CreateModel create model
	CreateModel(ctx context.Context, model *po_entity.ProviderModel) error

	// UpdateModel updates model
	UpdateModel(ctx context.Context, model *po_entity.ProviderModel) error

	// GetTenantModels get all models by searchModel
	GetTenantModel(ctx context.Context, tenantId, providerName, modelName, modelType string) (*po_entity.ProviderModel, error)
	// CreateTenantDefaultModel create a default model for the tenant
	CreateTenantDefaultModel(ctx context.Context, tenantDefaultModel *po_entity.TenantDefaultModel) (*po_entity.TenantDefaultModel, error)
	// Get the corresponding TenantDefaultModel record
	GetTenantDefaultModel(ctx context.Context, tenantId, modelType string) (*po_entity.TenantDefaultModel, error)
}
