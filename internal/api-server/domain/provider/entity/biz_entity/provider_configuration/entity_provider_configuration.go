// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package biz_entity

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"

	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider"
	biz_entity_model "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider/model_provider"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/infrastructure/log"
	apo "github.com/lunarianss/Luna/internal/api-server/domain/account/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/repository"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
)

// ModelConfigWithCredentialsEntity struct
type ModelConfigWithCredentialsEntity struct {
	Provider            string                                       `json:"provider"`
	Model               string                                       `json:"model"`
	ModelSchema         *biz_entity_model.AIModelStaticConfiguration `json:"model_schema"`
	Mode                string                                       `json:"mode"`
	ProviderModelBundle *ProviderModelBundleRuntime                  `json:"provider_model_bundle"`
	Credentials         interface{}                                  `json:"credentials"`
	Parameters          map[string]interface{}                       `json:"parameters"`
	Stop                []string                                     `json:"stop"`
}

type ModelWithProvider struct {
	*ProviderModelWithStatus
	Provider *SimpleModelProvider `json:"provider"`
}

type ProviderReposGetter interface {
	GetProviderRepo() repository.ProviderRepo
	GetModelRepo() repository.ModelRepo
}

type ModelIntegratedInstance struct {
	ProviderModelBundle *ProviderModelBundleRuntime      `json:"provider_model_bundle"`
	Model               string                           `json:"model"`
	Provider            string                           `json:"provider"`
	Credentials         map[string]interface{}           `json:"credentials"`
	ModelTypeInstance   *biz_entity_model.AIModelRuntime `json:"model_type_instance"`
}

type ProviderModelBundleRuntime struct {
	Configuration     *ProviderConfiguration
	ProviderInstance  *biz_entity.ProviderRuntime
	ModelTypeInstance *biz_entity_model.AIModelRuntime
}

type ProviderConfigurations struct {
	TenantId       string `json:"tenant_id"`
	ProviderRepo   repository.ProviderRepo
	ModelRepo      repository.ModelRepo
	Configurations map[string]*ProviderConfiguration `json:"configurations"`
}

func (pcm *ProviderConfigurations) GetModels(ctx context.Context, orderedProviders []string, provider string, modelType common.ModelType, onlyActive bool) ([]*ModelWithProvider, error) {
	var (
		providerModels []*ModelWithProvider
	)

	for _, orderedProvider := range orderedProviders {
		providerConfiguration, ok := pcm.Configurations[orderedProvider]

		if !ok {
			log.Warnf("%s provider is not in the configuration", orderedProvider)
			continue
		}

		if provider != "" && provider != providerConfiguration.Provider.Provider {
			continue
		}
		models, err := providerConfiguration.GetProviderModels(ctx, modelType, onlyActive)
		if err != nil {
			return nil, err
		}
		providerModels = append(providerModels, models...)
	}

	return providerModels, nil
}

func (pcm *ProviderConfigurations) GetConfigurationByProvider(ctx context.Context, provider string) (*ProviderConfiguration, error) {
	for providerName, configuration := range pcm.Configurations {
		if providerName == provider {
			return configuration, nil
		}
	}
	return nil, errors.WithCode(code.ErrRequiredCorrectProvider, fmt.Sprintf("provider %s not found", provider))
}

func (pcm *ProviderConfigurations) GetProviderRepo() repository.ProviderRepo {
	return pcm.ProviderRepo

}
func (pcm *ProviderConfigurations) GetModelRepo() repository.ModelRepo {
	return pcm.ModelRepo
}

type ProviderConfiguration struct {
	ProviderReposGetter
	TenantId              string                                  `json:"tenant_id"`
	Provider              *biz_entity.ProviderStaticConfiguration `json:"provider"`
	PreferredProviderType po_entity.ProviderType                  `json:"preferred_provider_type"`
	UsingProviderType     po_entity.ProviderType                  `json:"using_provider_type"`
	SystemConfiguration   *SystemConfiguration                    `json:"system_configuration"`
	CustomConfiguration   *CustomConfiguration                    `json:"custom_configuration"`
	ModelSettings         []*ModelSettings                        `json:"model_settings"`
}

func (c *ProviderConfiguration) ensureManager() error {
	if c.ProviderReposGetter == nil {
		return errors.WithCode(code.ErrNotSetManagerForProvider, "")
	}
	return nil
}

func (c *ProviderConfiguration) SetManager(manager ProviderReposGetter) {
	c.ProviderReposGetter = manager
}

func (c *ProviderConfiguration) GetCurrentCredentials(modelType common.ModelType, model string) (map[string]interface{}, error) {
	var credentials map[string]interface{}
	if c.CustomConfiguration.Models != nil {
		for _, modelConfiguration := range c.CustomConfiguration.Models {
			if modelConfiguration.ModelType == string(modelType) && modelConfiguration.Model == model {
				credentials = modelConfiguration.Credentials
				break
			}
		}
	}

	if c.CustomConfiguration.Provider != nil && credentials == nil {

		credentials = c.CustomConfiguration.Provider.Credentials
	}
	return credentials, nil

}

func (pc *ProviderConfiguration) GetProviderModels(ctx context.Context, modelType common.ModelType, onlyActive bool) ([]*ModelWithProvider, error) {

	if err := pc.ensureManager(); err != nil {
		return nil, err
	}

	var (
		providerModels []*ModelWithProvider
	)

	providerInstance, err := pc.GetProviderRepo().GetProviderInstance(ctx, pc.Provider.Provider)

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

	modelSettingMap := make(map[string]map[string]ModelSettings)

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
		providerModels = util.SliceFilter(providerModels, func(data *ModelWithProvider) bool {
			if data.Status == ACTIVE {
				return true
			} else {
				return false
			}
		})
	}
	return providerModels, nil

}

func (pc *ProviderConfiguration) getCustomProviderModels(modelTypes []common.ModelType, providerInstance *biz_entity.ProviderRuntime, modelSettingMap map[string]map[string]ModelSettings) ([]*ModelWithProvider, error) {

	var (
		providerModels []*ModelWithProvider
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
			var status ModelStatus
			if credentials != nil {
				status = ACTIVE
			} else {
				status = NO_CONFIGURE
			}

			if _, ok := modelSettingMap[string(modelType)]; ok {
				if modelSetting, ok := modelSettingMap[string(modelType)][AIModelEntity.Model]; ok {
					if !modelSetting.Enabled {
						status = DISABLED
					}
				}
			}

			modelWithProviderEntity := &ModelWithProvider{
				ProviderModelWithStatus: &ProviderModelWithStatus{
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
				Provider: &SimpleModelProvider{
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

func (pc *ProviderConfiguration) validateProviderCredentials(ctx context.Context, credentials map[string]interface{}, tenant *apo.Tenant) (*po_entity.Provider, map[string]interface{}, error) {
	var secretVariables []string
	if err := pc.ensureManager(); err != nil {
		return nil, nil, err
	}
	provider, err := pc.GetProviderRepo().GetTenantProvider(ctx, pc.TenantId, pc.Provider.Provider, string(po_entity.CUSTOM))

	if err != nil {
		return nil, nil, err
	}

	if pc.Provider.ProviderCredentialSchema != nil {
		secretVariables = pc.extractSecretVariables(pc.Provider.ProviderCredentialSchema.CredentialFormSchemas)
	}

	for key, value := range credentials {
		if slices.Contains(secretVariables, key) {
			encryptedData, err := util.Encrypt(value.(string), tenant.EncryptPublicKey)
			if err != nil {
				return nil, nil, errors.WithCode(code.ErrRunTimeCaller, err.Error())
			}
			credentials[key] = encryptedData
		}
	}

	return provider, credentials, nil
}

func (pc *ProviderConfiguration) extractSecretVariables(credentials []*biz_entity.CredentialFormSchema) []string {
	var secretInputVariables []string

	for _, credential := range credentials {
		if credential.Type == biz_entity.SECRET_INPUT {
			secretInputVariables = append(secretInputVariables, credential.Variable)
		}
	}
	return secretInputVariables
}

func (pc *ProviderConfiguration) AddOrUpdateCustomProviderCredentials(ctx context.Context, credentialParam map[string]interface{}, tenant *apo.Tenant) error {
	if err := pc.ensureManager(); err != nil {
		return err
	}
	providerRecord, credentials, err := pc.validateProviderCredentials(ctx, credentialParam, tenant)

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

		if err := pc.GetProviderRepo().UpdateProvider(ctx, providerRecord); err != nil {
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

		if err := pc.GetProviderRepo().CreateProvider(ctx, provider); err != nil {
			return err
		}
	}
	return nil
}

func (pc *ProviderConfiguration) AddOrUpdateCustomModelCredentials(ctx context.Context, credentialParam map[string]interface{}, modelType, modelName string) error {
	if err := pc.ensureManager(); err != nil {
		return err
	}
	modelRecord, credentials, err := pc.validateModelCredentials(ctx, credentialParam, modelType, modelName)

	if err != nil {
		return err
	}

	byteCredentials, err := json.Marshal(credentials)

	if err != nil {
		return errors.WithCode(code.ErrEncodingJSON, err.Error())
	}

	if modelRecord != nil {
		modelRecord.EncryptedConfig = string(byteCredentials)
		modelRecord.IsValid = 1

		if err := pc.GetModelRepo().UpdateModel(ctx, modelRecord); err != nil {
			return err
		}

	} else {
		model := &po_entity.ProviderModel{
			ProviderName:    pc.Provider.Provider,
			ModelName:       modelName,
			ModelType:       modelType,
			EncryptedConfig: string(byteCredentials),
			IsValid:         1,
			TenantID:        pc.TenantId,
		}

		if err := pc.GetModelRepo().CreateModel(ctx, model); err != nil {
			return err
		}
	}
	return nil
}

func (pc *ProviderConfiguration) validateModelCredentials(ctx context.Context, credentials map[string]interface{}, modelType, modeName string) (*po_entity.ProviderModel, map[string]interface{}, error) {
	if err := pc.ensureManager(); err != nil {
		return nil, nil, err
	}
	model, err := pc.GetModelRepo().GetTenantModel(ctx, pc.TenantId, pc.Provider.Provider, modeName, modelType)

	if err != nil {
		return nil, nil, err
	}

	// credentials 对 apikey 进行 validate and encrypt
	return model, credentials, nil
}
