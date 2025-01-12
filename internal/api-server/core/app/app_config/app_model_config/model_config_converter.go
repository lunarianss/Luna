// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package app_model_config

import (
	"context"

	biz_entity_app_config "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/provider_app_config"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_configuration"
)

type ModelConfigConverter struct {
	ProviderDomain *domain_service.ProviderDomain
}

func NewModelConfigConverter(providerDomain *domain_service.ProviderDomain) *ModelConfigConverter {
	return &ModelConfigConverter{
		ProviderDomain: providerDomain,
	}
}

func (c *ModelConfigConverter) Convert(ctx context.Context, appConfig *biz_entity_app_config.EasyUIBasedAppConfig, skipCheck bool) (*biz_entity.ModelConfigWithCredentialsEntity, error) {
	modelConfig := appConfig.Model

	providerModelBundle, err := c.ProviderDomain.GetProviderModelBundle(ctx, appConfig.TenantID, modelConfig.Provider, common.LLM)

	if err != nil {
		return nil, err
	}
	modelTypeInstance := providerModelBundle.ModelTypeInstance

	credentials, err := providerModelBundle.Configuration.GetCurrentCredentials(common.LLM, modelConfig.Model)

	if err != nil {
		return nil, err
	}

	AIModelEntity, err := c.ProviderDomain.GetModelSchema(ctx, modelConfig.Model, credentials, modelTypeInstance)

	if err != nil {
		return nil, err
	}

	return &biz_entity.ModelConfigWithCredentialsEntity{
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
