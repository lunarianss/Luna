// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package domain_service

import (
	"context"
	"encoding/json"
	"slices"
	"strings"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/infrastructure/log"
	model_providers "github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_providers"
	ac "github.com/lunarianss/Luna/internal/api-server/domain/account/repository"
	common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider"
	biz_entity_model "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider/model_provider"
	biz_entity_provider_config "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_configuration"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/repository"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
)

type ProviderDomain struct {
	ProviderRepo                  repository.ProviderRepo
	TenantRepo                    ac.TenantRepo
	ModelRepo                     repository.ModelRepo
	providerConfigurationsManager *providerConfigurationsManager
}

func NewProviderDomain(providerRepo repository.ProviderRepo, modelRepo repository.ModelRepo, tenantRepo ac.TenantRepo, providerConfigurationsManager *providerConfigurationsManager) *ProviderDomain {
	return &ProviderDomain{
		ProviderRepo:                  providerRepo,
		ModelRepo:                     modelRepo,
		providerConfigurationsManager: providerConfigurationsManager,
		TenantRepo:                    tenantRepo,
	}
}

// GetConfigurations Get all providers, models config for tenant
func (mpd *ProviderDomain) GetConfigurations(ctx context.Context, tenantId string) (*biz_entity_provider_config.ProviderConfigurations, []string, error) {

	providerNameMapRecords, err := mpd.ProviderRepo.GetMapTenantModelProviders(ctx, tenantId)

	if err != nil {
		return nil, nil, err
	}

	providerNameMapEntities, orderedProvider, err := mpd.ProviderRepo.GetSystemProviders(ctx)

	if err != nil {
		return nil, nil, err
	}

	providerConfigurations := NewProviderConfigurationsManager(mpd.ProviderRepo, mpd.ModelRepo, tenantId, make(map[string]*biz_entity_provider_config.ProviderConfiguration, model_providers.PROVIDER_COUNT))

	for _, providerEntity := range providerNameMapEntities {
		providerName := providerEntity.Provider
		providerRecords := providerNameMapRecords[providerName]
		customConfiguration, err := mpd.toCustomConfiguration(tenantId, providerEntity, providerRecords)

		if err != nil {
			return nil, nil, err
		}

		providerConfiguration := &biz_entity_provider_config.ProviderConfiguration{
			TenantId:              tenantId,
			Provider:              providerEntity,
			UsingProviderType:     po_entity.CUSTOM,
			PreferredProviderType: po_entity.SYSTEM,
			CustomConfiguration:   customConfiguration,
		}

		providerConfiguration.SetManager(providerConfigurations)
		providerConfigurations.Configurations[providerName] = providerConfiguration
	}

	return providerConfigurations.ProviderConfigurations, orderedProvider, nil
}

func (mpd *ProviderDomain) GetModelSchema(ctx context.Context, model string, credentials interface{}, AIModel *biz_entity_model.AIModelRuntime) (*biz_entity_model.AIModelStaticConfiguration, error) {

	AIModelEntities, err := AIModel.PredefinedModels()
	if err != nil {
		return nil, err
	}
	for _, modelEntity := range AIModelEntities {
		if modelEntity.Model == model {
			return modelEntity, nil
		}
	}
	return nil, errors.WithCode(code.ErrModelSchemaNotFound, "schema of model %s not found", model)
}

func (mpd *ProviderDomain) GetProviderModelBundle(ctx context.Context, tenantId, provider string, modelType common.ModelType) (*biz_entity_provider_config.ProviderModelBundleRuntime, error) {
	providerConfigurations, _, err := mpd.GetConfigurations(ctx, tenantId)

	if err != nil {
		return nil, err
	}

	providerConfiguration, ok := providerConfigurations.Configurations[provider]

	if !ok {
		return nil, errors.WithCode(code.ErrProviderMapModel, "provider %s not found", provider)
	}

	providerRuntime, err := mpd.ProviderRepo.GetProviderInstance(ctx, provider)

	if err != nil {
		return nil, err
	}

	AIModelInstance := providerRuntime.GetModelInstance(modelType)

	return &biz_entity_provider_config.ProviderModelBundleRuntime{
		Configuration:     providerConfiguration,
		ProviderInstance:  providerRuntime,
		ModelTypeInstance: AIModelInstance,
	}, nil
}

func (mpd *ProviderDomain) GetFirstProviderFirstModel(ctx context.Context, tenantID, modelType string) (string, string, error) {

	providerConfigurations, orderedProviders, err := mpd.GetConfigurations(ctx, tenantID)

	if err != nil {
		return "", "", err
	}

	firstProviderModels, err := providerConfigurations.Configurations[orderedProviders[0]].GetProviderModels(ctx, common.ModelType(modelType), false)

	if err != nil {
		return "", "", err
	}

	if len(firstProviderModels) == 0 {
		return "", "", errors.WithCode(code.ErrAllModelsEmpty, "tenant %s does not have any type of %s models", tenantID, modelType)
	}

	return firstProviderModels[0].Provider.Provider, firstProviderModels[0].Model, nil
}

func (mpd *ProviderDomain) GetModelInstance(ctx context.Context, tenantId, provider, model string, modelType common.ModelType) (*biz_entity_provider_config.ModelIntegratedInstance, error) {
	providerModelBundle, err := mpd.GetProviderModelBundle(ctx, tenantId, provider, modelType)

	if err != nil {
		return nil, err
	}

	return &biz_entity_provider_config.ModelIntegratedInstance{
		ProviderModelBundle: providerModelBundle,
		Model:               model,
		ModelTypeInstance:   providerModelBundle.ModelTypeInstance,
		Provider:            providerModelBundle.Configuration.Provider.Provider,
		Credentials:         providerModelBundle.Configuration.CustomConfiguration.Provider.Credentials,
	}, nil
}

func (mpd *ProviderDomain) GetDefaultModelInstance(ctx context.Context, tenantId string, modelType common.ModelType) (*biz_entity_provider_config.ModelIntegratedInstance, error) {
	defaultModelEntity, err := mpd.GetDefaultModel(ctx, tenantId, modelType)

	if err != nil {
		return nil, err
	}

	return mpd.GetModelInstance(ctx, tenantId, defaultModelEntity.Provider.Provider, defaultModelEntity.Model, modelType)
}

func (mpd *ProviderDomain) GetDefaultModel(ctx context.Context, tenantId string, modelType common.ModelType) (*biz_entity_model.DefaultModel, error) {

	var (
		defaultModel *po_entity.TenantDefaultModel
		err          error
	)

	originType, err := modelType.ToOriginModelType()

	if err != nil {
		return nil, err
	}

	defaultModel, err = mpd.ModelRepo.GetTenantDefaultModel(ctx, tenantId, originType)

	if err != nil {
		return nil, err
	}

	if defaultModel == nil {
		providerConfigurations, orderedProviders, err := mpd.GetConfigurations(ctx, tenantId)

		if err != nil {
			return nil, err
		}

		for _, orderedProvider := range orderedProviders {
			providerConfiguration, ok := providerConfigurations.Configurations[orderedProvider]

			if !ok {
				log.Warnf("%s provider is not in the configuration", orderedProvider)
				continue
			}

			availableModels, err := providerConfiguration.GetProviderModels(ctx, modelType, true)

			if err != nil {
				return nil, err
			}

			if availableModels != nil {
				availableModel := util.SliceFind(availableModels, func(t *biz_entity_provider_config.ModelWithProvider) bool {
					return t.Model == "gpt-4"
				})

				if availableModel == nil {
					availableModel = availableModels[0]
				}

				defaultModel, err = mpd.ModelRepo.CreateTenantDefaultModel(ctx, &po_entity.TenantDefaultModel{
					TenantID:     tenantId,
					ModelType:    originType,
					ProviderName: providerConfiguration.Provider.Provider,
					ModelName:    availableModel.Model,
				})

				if err != nil {
					return nil, err
				}
				break
			}
		}

	}

	if defaultModel == nil {
		return nil, errors.WithCode(code.ErrDefaultModelNotFound, "default %s model not found", modelType)
	}

	providerInstance, err := mpd.ProviderRepo.GetProviderInstance(ctx, defaultModel.ProviderName)

	if err != nil {
		return nil, err
	}

	providerSchema, err := providerInstance.GetProviderSchema()

	if err != nil {
		return nil, err
	}

	return &biz_entity_model.DefaultModel{
		Model:     defaultModel.ModelName,
		ModelType: string(modelType),
		Provider: &biz_entity_model.DefaultModelProvider{
			Provider:  providerSchema.Provider,
			Label:     providerSchema.Label,
			IconSmall: providerSchema.IconSmall,
			IconLarge: providerSchema.IconLarge,
		},
	}, nil

}

func (mpd *ProviderDomain) SaveProviderCredentials(ctx context.Context, tenantID string, provider string, credentials map[string]interface{}) error {
	providerConfigurations, _, err := mpd.GetConfigurations(ctx, tenantID)
	if err != nil {
		return err
	}

	providerConfiguration, ok := providerConfigurations.Configurations[provider]

	if !ok {
		return errors.WithCode(code.ErrProviderMapModel, "when create %s provider credential for provider", provider)
	}

	tenantRecord, err := mpd.TenantRepo.GetTenantByID(ctx, tenantID)
	if err != nil {
		return err
	}

	if err := providerConfiguration.AddOrUpdateCustomProviderCredentials(ctx, credentials, tenantRecord); err != nil {
		return err
	}
	return nil
}

func (mpd *ProviderDomain) toCustomConfiguration(
	_ string,
	providerEntity *biz_entity.ProviderStaticConfiguration,
	providerRecords []*po_entity.Provider,
) (*biz_entity_provider_config.CustomConfiguration, error) {

	var (
		custom_provider_record *po_entity.Provider
		// todo 从缓存中取 credentials information
		cache_provider_credentials        map[string]interface{}
		providerCredentials               map[string]interface{}
		customProviderConfiguration       *biz_entity_provider_config.CustomProviderConfiguration
		providerCredentialSecretVariables []string
	)

	if providerEntity.ProviderCredentialSchema != nil {
		providerCredentialSecretVariables = mpd.extractSecretVariables(
			providerEntity.ProviderCredentialSchema.CredentialFormSchemas,
		)
	}
	for _, providerRecord := range providerRecords {
		if providerRecord.ProviderType == string(po_entity.SYSTEM) {
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

	if custom_provider_record != nil {

		for k, v := range providerCredentials {

			if slices.Contains(providerCredentialSecretVariables, k) {
				decryptedData, err := util.Decrypt(v.(string), custom_provider_record.TenantID, &util.FileStorage{})

				if err != nil {
					return nil, err
				}
				providerCredentials[k] = decryptedData
			}

		}
		customProviderConfiguration = &biz_entity_provider_config.CustomProviderConfiguration{
			Credentials: providerCredentials,
		}
	}

	return &biz_entity_provider_config.CustomConfiguration{
		Provider: customProviderConfiguration,
	}, nil
}

func (mpd *ProviderDomain) extractSecretVariables(credentials []*biz_entity.CredentialFormSchema) []string {
	var secretInputVariables []string

	for _, credential := range credentials {
		if credential.Type == biz_entity.SECRET_INPUT {
			secretInputVariables = append(secretInputVariables, credential.Variable)
		}
	}
	return secretInputVariables
}
