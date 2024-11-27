package domain_service

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	common "github.com/lunarianss/Luna/internal/api-server/_domain/provider/entity/biz_entity/common_relation"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/_domain/provider/entity/biz_entity/provider"

	biz_entity_model "github.com/lunarianss/Luna/internal/api-server/_domain/provider/entity/biz_entity/provider/model_provider"
	biz_entity_provider_config "github.com/lunarianss/Luna/internal/api-server/_domain/provider/entity/biz_entity/provider_configuration"
	"github.com/lunarianss/Luna/internal/api-server/_domain/provider/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/_domain/provider/repository"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	model_providers "github.com/lunarianss/Luna/internal/api-server/model_runtime/model_providers"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/internal/pkg/util"
	"github.com/lunarianss/Luna/pkg/errors"
	"github.com/lunarianss/Luna/pkg/log"
)

type ProviderDomain struct {
	ProviderRepo repository.ProviderRepo
	ModelRepo    repository.ModelRepo

	providerConfigurationsManager *ProviderConfigurationsManager
}

func NewProviderDomain(providerRepo repository.ProviderRepo, modelRepo repository.ModelRepo, providerConfigurationsManager *ProviderConfigurationsManager) *ProviderDomain {
	return &ProviderDomain{
		ProviderRepo:                  providerRepo,
		ModelRepo:                     modelRepo,
		providerConfigurationsManager: providerConfigurationsManager,
	}
}

// GetConfigurations Get all providers, models config for tenant
func (mpd *ProviderDomain) GetSortedListConfigurations(ctx context.Context, tenantId string) ([]*ProviderConfigurationManager, error) {
	var (
		providerListConfigurations []*ProviderConfigurationManager
	)
	providerNameMapRecords, err := mpd.ProviderRepo.GetMapTenantModelProviders(ctx, tenantId)

	if err != nil {
		return nil, err
	}

	providerNameMapEntities, err := mpd.ProviderRepo.GetSystemProviders(ctx)

	if err != nil {
		return nil, err
	}

	for _, providerEntity := range providerNameMapEntities {
		providerName := providerEntity.Provider
		providerRecords := providerNameMapRecords[providerName]
		customConfiguration := mpd.toCustomConfiguration(tenantId, providerEntity, providerRecords)

		providerConfiguration := &ProviderConfigurationManager{
			TenantId:              tenantId,
			Provider:              providerEntity,
			UsingProviderType:     po_entity.CUSTOM,
			PreferredProviderType: po_entity.CUSTOM,
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
func (mpd *ProviderDomain) GetConfigurations(ctx context.Context, tenantId string) (*ProviderConfigurationsManager, error) {
	providerNameMapRecords, err := mpd.ProviderRepo.GetMapTenantModelProviders(ctx, tenantId)

	if err != nil {
		return nil, err
	}

	providerNameMapEntities, err := mpd.ProviderRepo.GetSystemProviders(ctx)

	if err != nil {
		return nil, err
	}

	providerConfigurations := &ProviderConfigurationsManager{
		TenantId:       tenantId,
		Configurations: make(map[string]*ProviderConfigurationManager, model_providers.PROVIDER_COUNT),
	}

	for _, providerEntity := range providerNameMapEntities {
		providerName := providerEntity.Provider
		providerRecords := providerNameMapRecords[providerName]
		customConfiguration := mpd.toCustomConfiguration(tenantId, providerEntity, providerRecords)

		providerConfiguration := &ProviderConfigurationManager{
			TenantId:              tenantId,
			Provider:              providerEntity,
			UsingProviderType:     po_entity.CUSTOM,
			PreferredProviderType: po_entity.SYSTEM,
			CustomConfiguration:   customConfiguration,
		}

		providerConfigurations.Configurations[providerName] = providerConfiguration
	}

	return providerConfigurations, nil
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
	return nil, errors.WithCode(code.ErrModelSchemaNotFound, fmt.Sprintf("schema of model %s not found", model))
}

func (mpd *ProviderDomain) GetProviderModelBundle(ctx context.Context, tenantId, provider string, modelType common.ModelType) (*ProviderModelBundleRuntime, error) {
	providerConfigurations, err := mpd.GetConfigurations(ctx, tenantId)

	if err != nil {
		return nil, err
	}

	providerConfiguration, ok := providerConfigurations.Configurations[provider]

	if !ok {
		return nil, errors.WithCode(code.ErrProviderMapModel, fmt.Sprintf("provider %s not found", provider))
	}

	providerRuntime, err := mpd.ProviderRepo.GetProviderInstance(ctx, provider)

	if err != nil {
		return nil, err
	}

	AIModelInstance := providerRuntime.GetModelInstance(modelType)

	return &ProviderModelBundleRuntime{
		Configuration:     providerConfiguration,
		ProviderInstance:  providerRuntime,
		ModelTypeInstance: AIModelInstance,
	}, nil

}

func (mpd *ProviderDomain) GetFirstProviderFirstModel(ctx context.Context, tenantID, modelType string) (string, string, error) {

	var allModels []*biz_entity_provider_config.ModelWithProvider

	providerConfigurations, err := mpd.GetSortedListConfigurations(ctx, tenantID)

	if err != nil {
		return "", "", err
	}

	for _, providerConfiguration := range providerConfigurations {
		model, err := providerConfiguration.GetProviderModels(ctx, common.ModelType(modelType), false)

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

func (mpd *ProviderDomain) GetModelInstance(ctx context.Context, tenantId, provider, model string, modelType common.ModelType) (*ModelIntegratedInstance, error) {
	providerModelBundle, err := mpd.GetProviderModelBundle(ctx, tenantId, provider, modelType)

	if err != nil {
		return nil, err
	}

	return &ModelIntegratedInstance{
		ProviderModelBundle: providerModelBundle,
		Model:               model,
		ModelTypeInstance:   providerModelBundle.ModelTypeInstance,
		Provider:            providerModelBundle.Configuration.Provider.Provider,
		Credentials:         providerModelBundle.Configuration.CustomConfiguration.Provider.Credentials,
	}, nil
}

func (mpd *ProviderDomain) GetDefaultModelInstance(ctx context.Context, tenantId string, modelType common.ModelType) (*ModelIntegratedInstance, error) {
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

				originType, err := modelType.ToOriginModelType()
				if err != nil {
					return nil, err
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
			}
		}
	}

	if defaultModel == nil {
		return nil, errors.WithCode(code.ErrDefaultModelNotFound, "default model not found")
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

func (mpd *ProviderDomain) toCustomConfiguration(
	_ string,
	providerEntity *biz_entity.ProviderStaticConfiguration,
	providerRecords []*po_entity.Provider,
) *biz_entity_provider_config.CustomConfiguration {

	var (
		custom_provider_record *po_entity.Provider
		// todo 从缓存中取 credentials information
		cache_provider_credentials  map[string]interface{}
		providerCredentials         map[string]interface{}
		customProviderConfiguration *biz_entity_provider_config.CustomProviderConfiguration
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
		customProviderConfiguration = &biz_entity_provider_config.CustomProviderConfiguration{
			Credentials: providerCredentials,
		}
	}

	return &biz_entity_provider_config.CustomConfiguration{
		Provider: customProviderConfiguration,
	}
}
