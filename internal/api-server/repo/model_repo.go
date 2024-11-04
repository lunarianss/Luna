package repo

import model "github.com/lunarianss/Luna/internal/api-server/model/v1"

type ModelRepo interface {
	// GetTenantModels get all models by searchModel
	GetTenantModel(tenantId, providerName, modelName, modelType string) (*model.ProviderModel, error)
	// UpdateModel updates model
	UpdateModel(model *model.ProviderModel) error
	// CreateModel create model
	CreateModel(model *model.ProviderModel) error
}
