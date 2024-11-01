// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"slices"

	domain "github.com/lunarianss/Luna/internal/api-server/domain/model-provider"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/provider"
	"github.com/lunarianss/Luna/internal/api-server/model-runtime/entities"
	model_providers "github.com/lunarianss/Luna/internal/api-server/model-runtime/model-providers"
)

type ModelProviderService struct {
	ModelProviderDomain *domain.ModelProviderDomain
}

func NewModelProviderService(modelProviderDomain *domain.ModelProviderDomain) *ModelProviderService {
	return &ModelProviderService{ModelProviderDomain: modelProviderDomain}
}

func (mpSrv *ModelProviderService) GetProviderList(tenantId int64, modelType string) ([]*dto.ProviderResponse, error) {
	var customConfigurationStatus dto.CustomConfigurationStatus

	providerConfigurations, err := mpSrv.ModelProviderDomain.GetConfigurations(tenantId)
	if err != nil {
		return nil, err
	}

	providerListResponse := make([]*dto.ProviderResponse, 0, model_providers.PROVIDER_COUNT)

	for _, providerConfiguration := range providerConfigurations.Configurations {

		if modelType != "" {
			if !slices.Contains(providerConfiguration.Provider.SupportedModelTypes, entities.ModelType(modelType)) {
				continue
			}
		}

		if providerConfiguration.Provider != nil {
			customConfigurationStatus = dto.ACTIVE
		} else {
			customConfigurationStatus = dto.NO_CONFIGURE
		}

		providerListResponse = append(providerListResponse, &dto.ProviderResponse{
			Provider:                 providerConfiguration.Provider.Provider,
			Label:                    providerConfiguration.Provider.Label,
			Description:              providerConfiguration.Provider.Description,
			IconSmall:                providerConfiguration.Provider.IconSmall,
			IconLarge:                providerConfiguration.Provider.IconLarge,
			Background:               providerConfiguration.Provider.Background,
			Help:                     providerConfiguration.Provider.Help,
			SupportedModelTypes:      providerConfiguration.Provider.SupportedModelTypes,
			ConfigurationMethods:     providerConfiguration.Provider.ConfigurationMethods,
			ProviderCredentialSchema: providerConfiguration.Provider.ProviderCredentialSchema,
			ModelCredentialSchema:    providerConfiguration.Provider.ModelCredentialSchema,
			PreferredProviderType:    providerConfiguration.PreferredProviderType,
			CustomConfiguration: &dto.CustomConfigurationResponse{
				Status: customConfigurationStatus,
			},
		})
	}
	return providerListResponse, nil
}
