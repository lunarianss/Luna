package domain_service

import (
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_configuration"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/repository"
)

type providerConfigurationsManager struct {
	*biz_entity.ProviderConfigurations
}

func NewProviderConfigurationsManager(providerRepo repository.ProviderRepo, modelRepo repository.ModelRepo, tenantID string, configs map[string]*biz_entity.ProviderConfiguration) *providerConfigurationsManager {
	providerConfigurationsManager := &providerConfigurationsManager{}
	providerConfigurationsManager.ProviderConfigurations = &biz_entity.ProviderConfigurations{
		ProviderRepo:   providerRepo,
		ModelRepo:      modelRepo,
		Configurations: configs,
		TenantId:       tenantID,
	}
	return providerConfigurationsManager
}
