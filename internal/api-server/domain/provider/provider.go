// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model_provider

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/lunarianss/Luna/internal/api-server/entities/base"
	"github.com/lunarianss/Luna/internal/api-server/entities/model_provider"
	model_providers "github.com/lunarianss/Luna/internal/api-server/model_runtime/model-providers"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/repo"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/internal/pkg/util"
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
func (mpd *ModelProviderDomain) GetSortedListConfigurations(ctx context.Context, tenantId string) ([]*model_provider.ProviderConfiguration, error) {
	var (
		providerListConfigurations []*model_provider.ProviderConfiguration
	)
	providerNameMapRecords, err := mpd.ModelProviderRepo.GetMapTenantModelProviders(ctx, tenantId)

	if err != nil {
		return nil, err
	}

	providerNameMapEntities, err := mpd.ModelProviderRepo.GetSystemProviders(ctx)

	if err != nil {
		return nil, err
	}

	for _, providerEntity := range providerNameMapEntities {
		providerName := providerEntity.Provider
		providerRecords := providerNameMapRecords[providerName]
		customConfiguration := mpd.toCustomConfiguration(tenantId, providerEntity, providerRecords)

		providerConfiguration := &model_provider.ProviderConfiguration{
			TenantId:              tenantId,
			Provider:              providerEntity,
			UsingProviderType:     model.CUSTOM,
			PreferredProviderType: model.CUSTOM,
			CustomConfiguration:   customConfiguration,
		}

		providerListConfigurations = append(providerListConfigurations, providerConfiguration)
	}

	sort.Slice(providerListConfigurations, func(i, j int) bool {
		return providerListConfigurations[i].Provider.Position < providerListConfigurations[j].Provider.Position
	})

	return providerListConfigurations, nil
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
			UsingProviderType:     model.CUSTOM,
			PreferredProviderType: model.CUSTOM,
			CustomConfiguration:   customConfiguration,
		}

		providerConfigurations.Configurations[providerName] = providerConfiguration
	}

	return providerConfigurations, nil
}

func (mpd *ModelProviderDomain) GetModelSchema(ctx context.Context, model string, credentials interface{}, AIModel *model_provider.AIModel) (*model_provider.AIModelEntity, error) {

	AIModelEntities, err := AIModel.PredefinedModels()
	if err != nil {
		return nil, err
	}
	for _, modelEntity := range AIModelEntities {
		if modelEntity.Model == model {
			return modelEntity, nil
		}
	}
	return nil, errors.WithCode(code.ErrModelSchemaNotFound, fmt.Sprintf("schema of model %s not found", model))
}

func (mpd *ModelProviderDomain) getProviderModelBundle(ctx context.Context, tenantId, provider string, modelType base.ModelType) (*model_provider.ProviderModelBundle, error) {
	providerConfigurations, err := mpd.GetConfigurations(ctx, tenantId)

	if err != nil {
		return nil, err
	}

	providerConfiguration, ok := providerConfigurations.Configurations[provider]

	if !ok {
		return nil, errors.WithCode(code.ErrProviderMapModel, fmt.Sprintf("provider %s not found", provider))
	}

	providerInstance, err := mpd.ModelProviderRepo.GetProviderInstance(ctx, provider)

	if err != nil {
		return nil, err
	}

	AIModelInstance := providerInstance.GetModelInstance(modelType)

	return &model_provider.ProviderModelBundle{
		Configuration:     providerConfiguration,
		ProviderInstance:  providerInstance,
		ModelTypeInstance: AIModelInstance,
	}, nil

}

func (mpd *ModelProviderDomain) GetFirstProviderFirstModel(ctx context.Context, tenantID, modelType string) (string, string, error) {

	var allModels []*model_provider.ModelWithProviderEntity

	providerConfigurations, err := mpd.GetSortedListConfigurations(ctx, tenantID)

	if err != nil {
		return "", "", err
	}

	for _, providerConfiguration := range providerConfigurations {
		model, err := mpd.getProviderModels(ctx, providerConfiguration, base.ModelType(modelType), false)

		if err != nil {
			return "", "", err
		}
		allModels = append(allModels, model...)
	}

	if len(allModels) == 0 {
		return "", "", errors.WithCode(code.ErrAllModelsEmpty, fmt.Sprintf("tenant %s does not have any type of %s models", tenantID, modelType))
	}

	return allModels[0].Provider.Provider, allModels[0].Model, nil
}

func (mpd *ModelProviderDomain) GetModelInstance(ctx context.Context, tenantId, provider, model string, modelType base.ModelType) (*model_provider.ModelInstance, error) {
	providerModelBundle, err := mpd.getProviderModelBundle(ctx, tenantId, provider, modelType)

	if err != nil {
		return nil, err
	}

	return &model_provider.ModelInstance{
		ProviderModelBundle: providerModelBundle,
		Model:               model,
		ModelTypeInstance:   providerModelBundle.ModelTypeInstance,
		Provider:            providerModelBundle.Configuration.Provider.Provider,
		Credentials:         providerModelBundle.Configuration.CustomConfiguration.Provider.Credentials,
	}, nil
}

func (mpd *ModelProviderDomain) GetDefaultModelInstance(ctx context.Context, tenantId string, modelType base.ModelType) (*model_provider.ModelInstance, error) {
	defaultModelEntity, err := mpd.GetDefaultModel(ctx, tenantId, modelType)

	if err != nil {
		return nil, err
	}

	return mpd.GetModelInstance(ctx, tenantId, defaultModelEntity.Provider.Provider, defaultModelEntity.Model, modelType)
}

func (mpd *ModelProviderDomain) GetDefaultModel(ctx context.Context, tenantId string, modelType base.ModelType) (*model_provider.DefaultModelEntity, error) {

	var (
		defaultModel *model.TenantDefaultModel
		err          error
	)

	defaultModel, err = mpd.ModelRepo.GetTenantDefaultModel(ctx, tenantId, string(modelType))

	if err != nil {
		return nil, err
	}

	if defaultModel == nil {
		providerConfigurations, err := mpd.GetConfigurations(ctx, tenantId)

		if err != nil {
			return nil, err
		}

		for _, providerConfiguration := range providerConfigurations.Configurations {
			availableModels, err := mpd.getProviderModels(ctx, providerConfiguration, modelType, true)

			if err != nil {
				return nil, err
			}

			if availableModels != nil {

				availableModel := util.SliceFind(availableModels, func(t *model_provider.ModelWithProviderEntity) bool {
					return t.Model == "gpt-4"
				})

				if availableModel == nil {
					availableModel = availableModels[0]
				}

				originType, err := modelType.ToOriginModelType()
				if err != nil {
					return nil, err
				}

				defaultModel, err = mpd.ModelRepo.CreateTenantDefaultModel(ctx, &model.TenantDefaultModel{
					TenantID:     tenantId,
					ModelType:    originType,
					ProviderName: providerConfiguration.Provider.Provider,
					ModelName:    availableModel.Model,
				})

				if err != nil {
					return nil, err
				}
			}
		}
	}

	if defaultModel == nil {
		return nil, errors.WithCode(code.ErrDefaultModelNotFound, "default model not found")
	}

	providerInstance, err := mpd.ModelProviderRepo.GetProviderInstance(ctx, defaultModel.ProviderName)

	if err != nil {
		return nil, err
	}

	providerSchema, err := providerInstance.GetProviderSchema()

	if err != nil {
		return nil, err
	}

	return &model_provider.DefaultModelEntity{
		Model:     defaultModel.ModelName,
		ModelType: string(modelType),
		Provider: &model_provider.DefaultModelProviderEntity{
			Provider:  providerSchema.Provider,
			Label:     providerSchema.Label,
			IconSmall: providerSchema.IconSmall,
			IconLarge: providerSchema.IconLarge,
		},
	}, nil

}

func (mpd *ModelProviderDomain) getCustomProviderModels(modelTypes []base.ModelType, providerInstance *model_provider.ModelProvider, modelSettingMap map[string]map[string]model_provider.ModelSettings, providerConfiguration *model_provider.ProviderConfiguration) ([]*model_provider.ModelWithProviderEntity, error) {

	var (
		providerModels []*model_provider.ModelWithProviderEntity
		credentials    interface{}
	)

	if providerConfiguration.CustomConfiguration.Provider != nil {
		credentials = providerConfiguration.CustomConfiguration.Provider.Credentials
	}

	for _, modelType := range modelTypes {
		if !slices.Contains(providerConfiguration.Provider.SupportedModelTypes, modelType) {
			continue
		}
		AIModelEntities, err := providerInstance.Models(modelType)

		if err != nil {
			return nil, err
		}

		for _, AIModelEntity := range AIModelEntities {
			var status model_provider.ModelStatus
			if credentials != nil {
				status = model_provider.ACTIVE
			} else {
				status = model_provider.NO_CONFIGURE
			}

			if _, ok := modelSettingMap[string(modelType)]; ok {
				if modelSetting, ok := modelSettingMap[string(modelType)][AIModelEntity.Model]; ok {
					if !modelSetting.Enabled {
						status = model_provider.DISABLED
					}
				}
			}

			modelWithProviderEntity := &model_provider.ModelWithProviderEntity{
				ProviderModelWithStatusEntity: &model_provider.ProviderModelWithStatusEntity{
					ProviderModel: &model_provider.ProviderModel{
						Model:           AIModelEntity.Model,
						Label:           AIModelEntity.Label,
						ModelType:       AIModelEntity.ModelType,
						Features:        AIModelEntity.Features,
						FetchFrom:       AIModelEntity.FetchFrom,
						ModelProperties: AIModelEntity.ModelProperties,
						Deprecated:      AIModelEntity.Deprecated,
					},
					Status: status,
				},
				Provider: &model_provider.SimpleModelProviderEntity{
					Provider:            providerConfiguration.Provider.Provider,
					Label:               providerConfiguration.Provider.Label,
					IconSmall:           providerConfiguration.Provider.IconSmall,
					IconLarge:           providerConfiguration.Provider.IconLarge,
					SupportedModelTypes: providerConfiguration.Provider.SupportedModelTypes,
				},
			}

			providerModels = append(providerModels, modelWithProviderEntity)
		}
	}
	return providerModels, nil
}

func (mpd *ModelProviderDomain) getProviderModels(ctx context.Context, providerConfiguration *model_provider.ProviderConfiguration, modelType base.ModelType, onlyActive bool) ([]*model_provider.ModelWithProviderEntity, error) {

	var (
		providerModels []*model_provider.ModelWithProviderEntity
	)

	providerInstance, err := mpd.ModelProviderRepo.GetProviderInstance(ctx, providerConfiguration.Provider.Provider)

	if err != nil {
		return nil, err
	}

	modelTypes := make([]base.ModelType, 0, 2)

	if modelType != "" {
		modelTypes = append(modelTypes, modelType)
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
		providerModels, err = mpd.getCustomProviderModels(modelTypes, providerInstance, modelSettingMap, providerConfiguration)

		if err != nil {
			return nil, err
		}
	}

	if onlyActive {
		providerModels = util.SliceFilter(providerModels, func(data *model_provider.ModelWithProviderEntity) bool {
			if data.Status == model_provider.ACTIVE {
				return true
			} else {
				return false
			}
		})
	}
	return providerModels, nil

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
