package domain_service

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	common "github.com/lunarianss/Luna/internal/api-server/_domain/provider/entity/biz_entity/common_relation"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/_domain/provider/entity/biz_entity/provider"
	biz_entity_model "github.com/lunarianss/Luna/internal/api-server/_domain/provider/entity/biz_entity/provider/model_provider"
	biz_entity_provider_config "github.com/lunarianss/Luna/internal/api-server/_domain/provider/entity/biz_entity/provider_configuration"
	"github.com/lunarianss/Luna/internal/api-server/_domain/provider/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/_domain/provider/repository"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/internal/pkg/util"
	"github.com/lunarianss/Luna/pkg/errors"
)

type ModelIntegratedInstance struct {
	ProviderModelBundle *ProviderModelBundleRuntime      `json:"provider_model_bundle"`
	Model               string                           `json:"model"`
	Provider            string                           `json:"provider"`
	Credentials         interface{}                      `json:"credentials"`
	ModelTypeInstance   *biz_entity_model.AIModelRuntime `json:"model_type_instance"`
}

type ProviderModelBundleRuntime struct {
	Configuration     *ProviderConfigurationManager
	ProviderInstance  *biz_entity.ProviderRuntime
	ModelTypeInstance *biz_entity_model.AIModelRuntime
}

type ProviderConfigurationsManager struct {
	TenantId       string                                   `json:"tenant_id"`
	Configurations map[string]*ProviderConfigurationManager `json:"configurations"`
}

func (pcm *ProviderConfigurationsManager) GetModels(ctx context.Context, modelType common.ModelType, onlyActive bool) ([]*biz_entity_provider_config.ModelWithProvider, error) {
	var (
		providerModels []*biz_entity_provider_config.ModelWithProvider
	)

	for _, providerConfiguration := range pcm.Configurations {
		models, err := providerConfiguration.GetProviderModels(ctx, modelType, onlyActive)

		if err != nil {
			return nil, err
		}

		providerModels = append(providerModels, models...)
	}
	return providerModels, nil
}

func (pcm *ProviderConfigurationsManager) GetConfigurationByProvider(ctx context.Context, provider string) (*ProviderConfigurationManager, error) {
	for providerName, configuration := range pcm.Configurations {
		if providerName == provider {
			return configuration, nil
		}
	}
	return nil, errors.WithCode(code.ErrRequiredCorrectProvider, fmt.Sprintf("provider %s not found", provider))
}

type ProviderConfigurationManager struct {
	ProviderRepo          repository.ProviderRepo
	TenantId              string                                          `json:"tenant_id"`
	Provider              *biz_entity.ProviderStaticConfiguration         `json:"provider"`
	PreferredProviderType po_entity.ProviderType                          `json:"preferred_provider_type"`
	UsingProviderType     po_entity.ProviderType                          `json:"using_provider_type"`
	SystemConfiguration   *biz_entity_provider_config.SystemConfiguration `json:"system_configuration"`
	CustomConfiguration   *biz_entity_provider_config.CustomConfiguration `json:"custom_configuration"`
	ModelSettings         []*biz_entity_provider_config.ModelSettings     `json:"model_settings"`
}

func (pc *ProviderConfigurationManager) GetProviderModels(ctx context.Context, modelType common.ModelType, onlyActive bool) ([]*biz_entity_provider_config.ModelWithProvider, error) {

	var (
		providerModels []*biz_entity_provider_config.ModelWithProvider
	)

	providerInstance, err := pc.ProviderRepo.GetProviderInstance(ctx, pc.Provider.Provider)

	if err != nil {
		return nil, err
	}

	modelTypes := make([]common.ModelType, 0, 2)

	if modelType != "" {
		modelTypes = append(modelTypes, modelType)
	} else {
		providerEntity, err := providerInstance.GetProviderSchema()
		if err != nil {
			return nil, err
		}
		modelTypes = append(modelTypes, providerEntity.SupportedModelTypes...)
	}

	modelSettingMap := make(map[string]map[string]biz_entity_provider_config.ModelSettings)

	for _, modelSetting := range pc.ModelSettings {
		modelSettingMap[string(modelSetting.Model)][modelSetting.Model] = *modelSetting
	}

	if pc.UsingProviderType == po_entity.CUSTOM {
		providerModels, err = pc.getCustomProviderModels(modelTypes, providerInstance, modelSettingMap)
		if err != nil {
			return nil, err
		}
	}

	if onlyActive {
		providerModels = util.SliceFilter(providerModels, func(data *biz_entity_provider_config.ModelWithProvider) bool {
			if data.Status == biz_entity_provider_config.ACTIVE {
				return true
			} else {
				return false
			}
		})
	}
	return providerModels, nil

}

func (pc *ProviderConfigurationManager) getCustomProviderModels(modelTypes []common.ModelType, providerInstance *biz_entity.ProviderRuntime, modelSettingMap map[string]map[string]biz_entity_provider_config.ModelSettings) ([]*biz_entity_provider_config.ModelWithProvider, error) {

	var (
		providerModels []*biz_entity_provider_config.ModelWithProvider
		credentials    interface{}
	)

	if pc.CustomConfiguration.Provider != nil {
		credentials = pc.CustomConfiguration.Provider.Credentials
	}

	for _, modelType := range modelTypes {
		if !slices.Contains(pc.Provider.SupportedModelTypes, modelType) {
			continue
		}
		AIModelEntities, err := providerInstance.Models(modelType)

		if err != nil {
			return nil, err
		}

		for _, AIModelEntity := range AIModelEntities {
			var status biz_entity_provider_config.ModelStatus
			if credentials != nil {
				status = biz_entity_provider_config.ACTIVE
			} else {
				status = biz_entity_provider_config.NO_CONFIGURE
			}

			if _, ok := modelSettingMap[string(modelType)]; ok {
				if modelSetting, ok := modelSettingMap[string(modelType)][AIModelEntity.Model]; ok {
					if !modelSetting.Enabled {
						status = biz_entity_provider_config.DISABLED
					}
				}
			}

			modelWithProviderEntity := &biz_entity_provider_config.ModelWithProvider{
				ProviderModelWithStatus: &biz_entity_provider_config.ProviderModelWithStatus{
					ProviderModel: &common.ProviderModel{
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
				Provider: &biz_entity_provider_config.SimpleModelProvider{
					Provider:            pc.Provider.Provider,
					Label:               pc.Provider.Label,
					IconSmall:           pc.Provider.IconSmall,
					IconLarge:           pc.Provider.IconLarge,
					SupportedModelTypes: pc.Provider.SupportedModelTypes,
				},
			}

			providerModels = append(providerModels, modelWithProviderEntity)
		}
	}
	return providerModels, nil
}

func (pc *ProviderConfigurationManager) validateProviderCredentials(ctx context.Context, credentials map[string]interface{}) (*po_entity.Provider, map[string]interface{}, error) {

	provider, err := pc.ProviderRepo.GetTenantProvider(ctx, pc.TenantId, pc.Provider.Provider, string(po_entity.CUSTOM))

	if err != nil {
		return nil, nil, err
	}
	//todo credentials 对 apikey 进行 validate and encrypt
	return provider, credentials, nil
}

func (pc *ProviderConfigurationManager) AddOrUpdateCustomProviderCredentials(ctx context.Context, credentialParam map[string]interface{}) error {

	providerRecord, credentials, err := pc.validateProviderCredentials(ctx, credentialParam)

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

		if err := pc.ProviderRepo.UpdateProvider(ctx, providerRecord); err != nil {
			return err
		}

	} else {
		provider := &po_entity.Provider{
			ProviderName:    pc.Provider.Provider,
			ProviderType:    string(po_entity.CUSTOM),
			EncryptedConfig: string(byteCredentials),
			IsValid:         1,
			TenantID:        pc.TenantId,
		}

		if err := pc.ProviderRepo.CreateProvider(ctx, provider); err != nil {
			return err
		}
	}
	return nil
}
