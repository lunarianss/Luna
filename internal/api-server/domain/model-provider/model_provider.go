// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model_provider

import (
	"context"
	"encoding/json"
	"strings"

	providerEntities "github.com/lunarianss/Luna/internal/api-server/entities/provider"
	"github.com/lunarianss/Luna/internal/api-server/model-runtime/entities"
	model_providers "github.com/lunarianss/Luna/internal/api-server/model-runtime/model-providers"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/repo"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
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

// GetConfigurations Get all providers, models config for tenant
func (mpd *ModelProviderDomain) GetConfigurations(ctx context.Context, tenantId string) (*providerEntities.ProviderConfigurations, error) {
	providerNameMapRecords, err := mpd.ModelProviderRepo.GetMapTenantModelProviders(ctx, tenantId)

	if err != nil {
		return nil, err
	}

	providerNameMapEntities, err := mpd.ModelProviderRepo.GetSystemProviders(ctx)

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

func (mpd *ModelProviderDomain) validateProviderCredentials(ctx context.Context, providerConfiguration *providerEntities.ProviderConfiguration, credentials map[string]interface{}) (*model.Provider, map[string]interface{}, error) {

	provider, err := mpd.ModelProviderRepo.GetTenantProvider(ctx, providerConfiguration.TenantId, providerConfiguration.Provider.Provider, string(model.CUSTOM))

	if err != nil {
		return nil, nil, err
	}
	// credentials 对 apikey 进行 validate and encrypt
	return provider, credentials, nil
}

func (mpd *ModelProviderDomain) AddOrUpdateCustomProviderCredentials(ctx context.Context, providerConfiguration *providerEntities.ProviderConfiguration, credentialParam map[string]interface{}) error {

	providerRecord, credentials, err := mpd.validateProviderCredentials(ctx, providerConfiguration, credentialParam)

	if err != nil {
		return err
	}

	byteCredentials, err := json.Marshal(credentials)

	if err != nil {
		return errors.WithCode(code.ErrEncodingJSON, err.Error())
	}

	if providerRecord != nil {
		providerRecord.EncryptedConfig = string(byteCredentials)
		providerRecord.IsValid = 1

		if err := mpd.ModelProviderRepo.UpdateProvider(ctx, providerRecord); err != nil {
			return err
		}

	} else {
		provider := &model.Provider{
			ProviderName:    providerConfiguration.Provider.Provider,
			ProviderType:    string(model.CUSTOM),
			EncryptedConfig: string(byteCredentials),
			IsValid:         1,
			TenantID:        providerConfiguration.TenantId,
		}

		if err := mpd.ModelProviderRepo.CreateProvider(ctx, provider); err != nil {
			return err
		}
	}
	return nil
}

func (mpd *ModelProviderDomain) toCustomConfiguration(
	_ string,
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
