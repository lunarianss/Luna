package model_config

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/core/app"
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/provider"
	"github.com/lunarianss/Luna/internal/api-server/entities/base"
)

type ModelConfigConverter struct {
	ProviderDomain *domain.ModelProviderDomain
}

func NewModelConfigConverter(providerDomain *domain.ModelProviderDomain) *ModelConfigConverter {
	return &ModelConfigConverter{
		ProviderDomain: providerDomain,
	}
}
func (c *ModelConfigConverter) Convert(ctx context.Context, appConfig *app_config.EasyUIBasedAppConfig, skipCheck bool) (*app.ModelConfigWithCredentialsEntity, error) {
	modelConfig := appConfig.Model

	providerModelBundle, err := c.ProviderDomain.GetProviderModelBundle(ctx, appConfig.TenantID, modelConfig.Provider, base.LLM)

	if err != nil {
		return nil, err
	}
	modelTypeInstance := providerModelBundle.ModelTypeInstance

	credentials, err := providerModelBundle.Configuration.GetCurrentCredentials(base.LLM, modelConfig.Model)

	if err != nil {
		return nil, err
	}

	AIModelEntity, err := c.ProviderDomain.GetModelSchema(ctx, modelConfig.Model, credentials, modelTypeInstance)

	if err != nil {
		return nil, err
	}

	return &app.ModelConfigWithCredentialsEntity{
		Provider:            modelConfig.Provider,
		Model:               modelConfig.Model,
		ModelSchema:         AIModelEntity,
		Mode:                modelConfig.Mode,
		ProviderModelBundle: providerModelBundle,
		Credentials:         credentials,
		Stop:                modelConfig.Stop,
		Parameters:          modelConfig.Parameters,
	}, nil

}
