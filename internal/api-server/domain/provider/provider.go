// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model_provider

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/lunarianss/Luna/internal/api-server/entities/base"
	"github.com/lunarianss/Luna/internal/api-server/entities/model_provider"
	model_providers "github.com/lunarianss/Luna/internal/api-server/model-runtime/model-providers"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/repo"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
	"github.com/lunarianss/Luna/pkg/log"
)

type ModelProviderDomain struct {
	ModelProviderRepo repo.ModelProviderRepo
	ModelRepo         repo.ModelRepo
}

func NewModelProviderDomain(modelProviderRepo repo.ModelProviderRepo, modelRepo repo.ModelRepo) *ModelProviderDomain {
	return &ModelProviderDomain{
		ModelProviderRepo: modelProviderRepo,
		ModelRepo:         modelRepo,
	}
}

// GetConfigurations Get all providers, models config for tenant
func (mpd *ModelProviderDomain) GetConfigurations(ctx context.Context, tenantId string) (*model_provider.ProviderConfigurations, error) {
	providerNameMapRecords, err := mpd.ModelProviderRepo.GetMapTenantModelProviders(ctx, tenantId)

	if err != nil {
		return nil, err
	}

	providerNameMapEntities, err := mpd.ModelProviderRepo.GetSystemProviders(ctx)

	if err != nil {
		return nil, err
	}

	providerConfigurations := &model_provider.ProviderConfigurations{
		TenantId:       tenantId,
		Configurations: make(map[string]*model_provider.ProviderConfiguration, model_providers.PROVIDER_COUNT),
	}

	for _, providerEntity := range providerNameMapEntities {
		providerName := providerEntity.Provider
		providerRecords := providerNameMapRecords[providerName]
		customConfiguration := mpd.toCustomConfiguration(tenantId, providerEntity, providerRecords)

		providerConfiguration := &model_provider.ProviderConfiguration{
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

func (mpd *ModelProviderDomain) GetDefaultModel(ctx context.Context, tenantId, modelType string) (*model_provider.DefaultModelEntity, error) {
	defaultModel, err := mpd.ModelRepo.GetTenantDefaultModel(ctx, tenantId, modelType)
	var allModels []*model_provider.ModelWithProviderEntity

	if err != nil {
		return nil, err
	}

	if defaultModel == nil {
		providerConfigurations, err := mpd.GetConfigurations(ctx, tenantId)

		if err != nil {
			return nil, err
		}

		for _, providerConfiguration := range providerConfigurations.Configurations {
			providerModels, err := mpd.getProviderModels(ctx, providerConfiguration, modelType, true)

			if err != nil {
				return nil, err
			}

			allModels = append(allModels, providerModels...)
		}
	}

	return nil, nil

}

func (mpd *ModelProviderDomain) getCustomProviderModels(modelType []base.ModelType, providerInstance *model_provider.ModelProvider, modelSettingMap map[string]map[string]model_provider.ModelSettings, providerConfiguration model_provider.ProviderConfiguration) ([]*model_provider.ModelWithProviderEntity, error) {

	// var (
	// 	providerModels []*model_provider.ModelWithProviderEntity
	// 	credentials    interface{}
	// )

	// if providerConfiguration.CustomConfiguration.Provider != nil {
	// 	credentials = providerConfiguration.CustomConfiguration.Provider.Credentials
	// }

	return nil, nil

}

func (mpd *ModelProviderDomain) getProviderModels(ctx context.Context, providerConfiguration *model_provider.ProviderConfiguration, modelType string, onlyActive bool) ([]*model_provider.ModelWithProviderEntity, error) {
	providerInstance, err := mpd.ModelProviderRepo.GetProviderInstance(ctx, providerConfiguration.Provider.Provider)

	if err != nil {
		return nil, err
	}

	modelTypes := make([]base.ModelType, 0, 2)

	if modelType != "" {
		modelTypes = append(modelTypes, base.ModelType(modelType))
	} else {
		providerEntity, err := providerInstance.GetProviderSchema()
		if err != nil {
			return nil, err
		}

		modelTypes = append(modelTypes, providerEntity.SupportedModelTypes...)
	}

	modelSettingMap := make(map[string]map[string]model_provider.ModelSettings)

	for _, modelSetting := range providerConfiguration.ModelSettings {
		modelSettingMap[string(modelSetting.Model)][modelSetting.Model] = *modelSetting
	}

	if providerConfiguration.UsingProviderType == model.CUSTOM {
		return nil, nil
	}

	return nil, nil

}

func (mpd *ModelProviderDomain) validateProviderCredentials(ctx context.Context, providerConfiguration *model_provider.ProviderConfiguration, credentials map[string]interface{}) (*model.Provider, map[string]interface{}, error) {

	provider, err := mpd.ModelProviderRepo.GetTenantProvider(ctx, providerConfiguration.TenantId, providerConfiguration.Provider.Provider, string(model.CUSTOM))

	if err != nil {
		return nil, nil, err
	}
	// credentials 对 apikey 进行 validate and encrypt
	return provider, credentials, nil
}

func (mpd *ModelProviderDomain) AddOrUpdateCustomProviderCredentials(ctx context.Context, providerConfiguration *model_provider.ProviderConfiguration, credentialParam map[string]interface{}) error {

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
	providerEntity *model_provider.ProviderEntity,
	providerRecords []*model.Provider,
) *model_provider.CustomConfiguration {

	var (
		custom_provider_record *model.Provider
		// todo 从缓存中取 credentials information
		cache_provider_credentials  map[string]interface{}
		providerCredentials         map[string]interface{}
		customProviderConfiguration *model_provider.CustomProviderConfiguration
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
		customProviderConfiguration = &model_provider.CustomProviderConfiguration{
			Credentials: providerCredentials,
		}
	}

	return &model_provider.CustomConfiguration{
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
