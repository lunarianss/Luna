// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"fmt"

	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	providerDomain "github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"
	biz_entity_model "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider/model_provider"
	biz_entity_provider_config "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_configuration"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/po_entity"

	"github.com/lunarianss/Luna/internal/api-server/config"

	"github.com/lunarianss/Luna/infrastructure/errors"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/provider"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
)

type ModelService struct {
	providerDomain *providerDomain.ProviderDomain
	accountDomain  *accountDomain.AccountDomain
	config         *config.Config
}

func NewModelService(providerDomain *providerDomain.ProviderDomain, accountDomain *accountDomain.AccountDomain, config *config.Config) *ModelService {
	return &ModelService{providerDomain: providerDomain, accountDomain: accountDomain, config: config}
}

func (ms *ModelService) SaveModelCredentials(ctx context.Context, tenantId, model, modelTpe, provider string, credentials map[string]interface{}) error {

	providerConfigurations, _, err := ms.providerDomain.GetConfigurations(ctx, tenantId)

	if err != nil {
		return err
	}

	providerConfiguration, ok := providerConfigurations.Configurations[provider]

	if !ok {
		return errors.WithCode(code.ErrProviderMapModel, "provider %s not found in map provider configuration", provider)
	}

	err = providerConfiguration.AddOrUpdateCustomModelCredentials(ctx, credentials, modelTpe, model)

	if err != nil {
		return err
	}

	return nil
}

func (ms *ModelService) GetAccountAvailableModels(ctx context.Context, accountID string, modelType common.ModelType) (*dto.DataWrapperResponse[[]*dto.ProviderWithModelsResponse], error) {

	tenantRecord, _, err := ms.accountDomain.GetCurrentTenantOfAccount(ctx, accountID)

	if err != nil {
		return nil, err
	}
	providerConfigurations, orderedProviders, err := ms.providerDomain.GetConfigurations(ctx, tenantRecord.ID)
	if err != nil {
		return nil, err
	}

	activeModels, err := providerConfigurations.GetModels(ctx, orderedProviders, "", common.ModelType(modelType), true)

	if err != nil {
		return nil, err
	}

	providerModelsMap := make(map[string][]*biz_entity_provider_config.ModelWithProvider)

	for _, activeModel := range activeModels {
		if _, ok := providerModelsMap[activeModel.Provider.Provider]; !ok {
			providerModelsMap[activeModel.Provider.Provider] = make([]*biz_entity_provider_config.ModelWithProvider, 0, 10)
		}

		if activeModel.Deprecated {
			continue
		}

		if activeModel.Status != biz_entity_provider_config.ACTIVE {
			continue
		}

		providerModelsMap[activeModel.Provider.Provider] = append(providerModelsMap[activeModel.Provider.Provider], activeModel)
	}

	providerResponses := make([]*dto.ProviderWithModelsResponse, 0, 2)

	for _, orderedProvider := range orderedProviders {
		providerModels := providerModelsMap[orderedProvider]

		if len(providerModels) == 0 {
			continue
		}

		providerModelStatus := make([]*biz_entity_provider_config.ProviderModelWithStatus, 0, 10)
		firstModel := providerModels[0]

		for _, mapModel := range providerModels {
			providerModelStatus = append(providerModelStatus, &biz_entity_provider_config.ProviderModelWithStatus{
				Status: mapModel.Status,
				ProviderModel: &common.ProviderModel{
					Model:           mapModel.Model,
					Label:           mapModel.Label,
					ModelType:       mapModel.ModelType,
					Features:        mapModel.Features,
					FetchFrom:       mapModel.FetchFrom,
					ModelProperties: mapModel.ModelProperties,
				},
			})
		}

		providerResponses = append(providerResponses, &dto.ProviderWithModelsResponse{
			Provider:  orderedProvider,
			Label:     firstModel.Provider.Label,
			IconSmall: firstModel.Provider.IconSmall,
			IconLarge: firstModel.Provider.Label,
			Status:    dto.ACTIVE,
			Models:    providerModelStatus,
		})

		util.PatchI18nObject(providerResponses)

		for _, p := range providerResponses {
			p.PatchIcon(ms.config)
		}
	}

	return &dto.DataWrapperResponse[[]*dto.ProviderWithModelsResponse]{
		Data: providerResponses,
	}, nil
}

func (ms *ModelService) GetModelParameterRules(ctx context.Context, accountID string, provider string, model string) (*dto.DataWrapperResponse[[]*biz_entity_model.ParameterRule], error) {
	tenantRecord, _, err := ms.accountDomain.GetCurrentTenantOfAccount(ctx, accountID)

	if err != nil {
		return nil, err
	}

	providerConfigurations, _, err := ms.providerDomain.GetConfigurations(ctx, tenantRecord.ID)

	if err != nil {
		return nil, err
	}

	providerConfiguration, err := providerConfigurations.GetConfigurationByProvider(ctx, provider)

	if err != nil {
		return nil, err
	}

	credentials, err := providerConfiguration.GetCurrentCredentials(common.LLM, model)

	if err != nil {
		return nil, err
	}

	model_provider, err := ms.providerDomain.ProviderRepo.GetProviderInstance(ctx, provider)

	if err != nil {
		return nil, err
	}

	modelInstance := model_provider.GetModelInstance(common.LLM)

	AIModelEntity, err := ms.providerDomain.GetModelSchema(ctx, model, credentials, modelInstance)

	if err != nil {
		return nil, err
	}

	return &dto.DataWrapperResponse[[]*biz_entity_model.ParameterRule]{
		Data: AIModelEntity.ParameterRules,
	}, nil

}

func (ms *ModelService) GetDefaultModelByType(ctx context.Context, accountID string, modelType string) (*dto.DataWrapperResponse[*dto.DefaultModelResponse], error) {

	tenantRecord, _, err := ms.accountDomain.GetCurrentTenantOfAccount(ctx, accountID)

	if err != nil {
		return nil, err
	}
	defaultModelEntity, err := ms.providerDomain.GetDefaultModel(ctx, tenantRecord.ID, common.ModelType(modelType))

	if errors.IsCode(err, code.ErrDefaultModelNotFound) {
		return dto.NewEmptyDataWrapperResponse[*dto.DefaultModelResponse](nil), nil
	}

	if err != nil {
		return nil, err
	}

	return &dto.DataWrapperResponse[*dto.DefaultModelResponse]{
		Data: &dto.DefaultModelResponse{
			Model:     defaultModelEntity.Model,
			ModelType: defaultModelEntity.ModelType,
			Provider: &biz_entity_provider_config.SimpleModelProvider{
				Provider:            defaultModelEntity.Provider.Provider,
				Label:               defaultModelEntity.Provider.Label,
				IconSmall:           defaultModelEntity.Provider.IconSmall,
				IconLarge:           defaultModelEntity.Provider.IconLarge,
				SupportedModelTypes: defaultModelEntity.Provider.SupportedModelTypes,
				Models:              make([]any, 0),
			},
		},
	}, nil
}

func (ms *ModelService) UpdateDefaultModel(ctx context.Context, accountID string, args []*dto.ModelSetting) error {
	tenantRecord, tenantJoin, err := ms.accountDomain.GetCurrentTenantOfAccount(ctx, accountID)

	if err != nil {
		return err
	}

	if !tenantJoin.IsPrivilegedRole() {
		return errors.WithCode(code.ErrForbidden, fmt.Sprintf("tenant %s don't have the permission", tenantRecord.Name))
	}

	providerConfigurations, _, err := ms.providerDomain.GetConfigurations(ctx, tenantRecord.ID)

	if err != nil {
		return err
	}

	for _, modelSetting := range args {

		if modelSetting.Provider == "" || modelSetting.Model == "" {
			continue
		}

		providerConfiguration, ok := providerConfigurations.Configurations[modelSetting.Provider]
		if !ok {
			return errors.WithCode(code.ErrRequiredCorrectProvider, fmt.Sprintf("provider %s is not exist", modelSetting.Provider))
		}

		providerModels, err := providerConfiguration.GetProviderModels(ctx, common.ModelType(modelSetting.ModelType), true)

		if err != nil {
			return err
		}

		findModel := util.SliceFind(providerModels, func(a *biz_entity_provider_config.ModelWithProvider) bool {
			return a.ProviderModel.Model == modelSetting.Model
		})

		if findModel == nil {
			return errors.WithCode(code.ErrRequiredCorrectModel, fmt.Sprintf("model %s is not exist", modelSetting.Model))
		}

		originModelType, err := common.ModelType(modelSetting.ModelType).ToOriginModelType()

		if err != nil {
			return err
		}
		defaultModel, err := ms.providerDomain.ModelRepo.GetTenantDefaultModel(ctx, tenantRecord.ID, originModelType)

		if err != nil {
			return err
		}

		if defaultModel != nil {
			defaultModel.ProviderName = modelSetting.Provider
			defaultModel.ModelName = modelSetting.Model
			if err := ms.providerDomain.ModelRepo.UpdateTenantDefaultModel(ctx, defaultModel); err != nil {
				return err
			}
		} else {
			defaultModel := &po_entity.TenantDefaultModel{
				TenantID:     tenantRecord.ID,
				ModelName:    modelSetting.Model,
				ProviderName: modelSetting.Provider,
				ModelType:    modelSetting.ModelType,
			}

			if _, err := ms.providerDomain.ModelRepo.CreateTenantDefaultModel(ctx, defaultModel); err != nil {
				return err
			}
		}
	}

	return nil
}
