package domain_service

import (
	biz_entity "github.com/lunarianss/Luna/internal/api-server/_domain/provider/entity/biz_entity/provider_configuration"
	"github.com/lunarianss/Luna/internal/api-server/_domain/provider/repository"
)

type providerConfigurationsManager struct {
	*biz_entity.ProviderConfigurationsManager
}

func NewProviderConfigurationsManager(providerRepo repository.ProviderRepo, modelRepo repository.ModelRepo, tenantID string, configs map[string]*ProviderConfigurationManager) *providerConfigurationsManager {
	providerConfigurationsManager = &providerConfigurationsManager{}
	providerConfigurationsManager.ProviderConfigurationsManager = &biz_entity.ProviderConfigurationsManager{
		ProviderRepo:   providerRepo,
		ModelRepo:      modelRepo,
		Configurations: configs,
		TenantId:       tenantID,
	}
	return providerConfigurationsManager
}
