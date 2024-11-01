// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model_provider

import (
	"encoding/json"
	"strings"

	providerEntities "github.com/lunarianss/Luna/internal/api-server/entities/provider"
	"github.com/lunarianss/Luna/internal/api-server/model-runtime/entities"
	model_providers "github.com/lunarianss/Luna/internal/api-server/model-runtime/model-providers"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/repo"
	"github.com/lunarianss/Luna/pkg/log"
)

type ModelProviderDomain struct {
	ModelProviderRepo repo.ModelProviderRepo
}

func NewModelProviderDomain(modelProviderRepo repo.ModelProviderRepo) *ModelProviderDomain {
	return &ModelProviderDomain{
		ModelProviderRepo: modelProviderRepo,
	}
}

// Get all providers, models config for tenant
func (mpd *ModelProviderDomain) GetConfigurations(tenantId int64) (*providerEntities.ProviderConfigurations, error) {
	providerNameMapRecords, err := mpd.ModelProviderRepo.GetMapTenantModelProviders(tenantId)

	if err != nil {
		return nil, err
	}

	providerNameMapEntities, err := mpd.ModelProviderRepo.GetSystemProviders()

	if err != nil {
		return nil, err
	}

	providerConfigurations := &providerEntities.ProviderConfigurations{
		TenantId:       tenantId,
		Configurations: make(map[string]*providerEntities.ProviderConfiguration, model_providers.PROVIDER_COUNT),
	}

	for _, providerEntity := range providerNameMapEntities {
		providerName := providerEntity.Provider
		providerRecords := providerNameMapRecords[providerName]
		customConfiguration := mpd.toCustomConfiguration(tenantId, providerEntity, providerRecords)

		providerConfiguration := &providerEntities.ProviderConfiguration{
			TenantId:              tenantId,
			Provider:              providerEntity,
			UsingProviderType:     "system",
			PreferredProviderType: "system",
			CustomConfiguration:   customConfiguration,
		}

		providerConfigurations.Configurations[providerName] = providerConfiguration
	}

	return providerConfigurations, nil
}

func (mpd *ModelProviderDomain) toCustomConfiguration(
	_ int64,
	providerEntity *entities.ProviderEntity,
	providerRecords []*model.Provider,
) *providerEntities.CustomConfiguration {

	var (
		custom_provider_record *model.Provider
		// todo 从缓存中取 credentials information
		cache_provider_credentials  map[string]interface{}
		providerCredentials         map[string]interface{}
		customProviderConfiguration *providerEntities.CustomProviderConfiguration
	)

	// provider_credential_secret_variables := mpd.extractSecretVariables(
	// 	providerEntity.ProviderCredentialSchema.CredentialFormSchemas,
	// )

	for _, providerRecord := range providerRecords {
		if providerRecord.ProviderType == string(model.SYSTEM) {
			continue
		}

		if providerRecord.EncryptedConfig == "" {
			continue
		}
		custom_provider_record = providerRecord
	}

	if len(cache_provider_credentials) == 0 {
		if custom_provider_record != nil {
			if !strings.HasPrefix(custom_provider_record.EncryptedConfig, "{") {
				providerCredentials = map[string]interface{}{
					"openai_api_key": custom_provider_record.EncryptedConfig,
				}
			} else {
				if err := json.Unmarshal([]byte(custom_provider_record.EncryptedConfig), &providerCredentials); err != nil {
					log.Errorf("error occurred when unmarshal %s encryptedConfig", providerEntity.Provider)
					providerCredentials = map[string]interface{}{}
				}
			}
		}
	}

	// todo 对用户的 api key 进行加密｜解密

	if custom_provider_record != nil {
		customProviderConfiguration = &providerEntities.CustomProviderConfiguration{
			Credentials: providerCredentials,
		}
	}

	return &providerEntities.CustomConfiguration{
		Provider: customProviderConfiguration,
	}
}

// func (mpd *ModelProviderDomain) extractSecretVariables(
// 	credentialFromSchemas []*entities.CredentialFormSchema,
// ) []string {

// 	var secretInputFormVariables []string

// 	for _, credentialFromSchema := range credentialFromSchemas {
// 		if credentialFromSchema.Type == entities.SECRET_INPUT {
// 			secretInputFormVariables = append(secretInputFormVariables, credentialFromSchema.Variable)
// 		}
// 	}
// 	return secretInputFormVariables
// }
