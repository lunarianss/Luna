// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/config"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/model"
	providerDomain "github.com/lunarianss/Luna/internal/api-server/domain/provider"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/provider"
	"github.com/lunarianss/Luna/internal/api-server/entities/base"
	"github.com/lunarianss/Luna/internal/api-server/entities/model_provider"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/internal/pkg/util"
	"github.com/lunarianss/Luna/pkg/errors"
)

type ModelService struct {
	ModelDomain         *domain.ModelDomain
	ModelProviderDomain *providerDomain.ModelProviderDomain
	AccountDomain       *accountDomain.AccountDomain
	config              *config.Config
}

func NewModelService(modelDomain *domain.ModelDomain, modelProviderDomain *providerDomain.ModelProviderDomain, accountDomain *accountDomain.AccountDomain, config *config.Config) *ModelService {
	return &ModelService{ModelDomain: modelDomain, ModelProviderDomain: modelProviderDomain, AccountDomain: accountDomain, config: config}
}

func (ms *ModelService) SaveModelCredentials(ctx context.Context, tenantId, model, modelTpe, provider string, credentials map[string]interface{}) error {

	providerConfigurations, err := ms.ModelProviderDomain.GetConfigurations(ctx, tenantId)

	if err != nil {
		return err
	}

	providerConfiguration, ok := providerConfigurations.Configurations[provider]

	if !ok {
		return errors.WithCode(code.ErrProviderMapModel, "provider %s not found in map provider configuration", provider)
	}

	err = ms.ModelDomain.AddOrUpdateCustomModelCredentials(ctx, providerConfiguration, credentials, modelTpe, model)

	if err != nil {
		return err
	}

	return nil
}

func (ms *ModelService) GetAccountAvailableModels(ctx context.Context, accountID string, modelType base.ModelType) ([]*dto.ProviderWithModelsResponse, error) {

	tenantRecord, _, err := ms.AccountDomain.GetCurrentTenantOfAccount(ctx, accountID)

	if err != nil {
		return nil, err
	}
	providerConfigurations, err := ms.ModelProviderDomain.GetConfigurations(ctx, tenantRecord.ID)
	if err != nil {
		return nil, err
	}

	activeModels, err := ms.ModelProviderDomain.GetModels(ctx, providerConfigurations, base.ModelType(modelType), true)

	if err != nil {
		return nil, err
	}

	providerModelsMap := make(map[string][]*model_provider.ModelWithProviderEntity)

	for _, activeModel := range activeModels {
		if _, ok := providerModelsMap[activeModel.Provider.Provider]; !ok {
			providerModelsMap[activeModel.Provider.Provider] = make([]*model_provider.ModelWithProviderEntity, 0, 10)
		}

		if activeModel.Deprecated {
			continue
		}

		if activeModel.Status != model_provider.ACTIVE {
			continue
		}

		providerModelsMap[activeModel.Provider.Provider] = append(providerModelsMap[activeModel.Provider.Provider], activeModel)
	}

	providerResponses := make([]*dto.ProviderWithModelsResponse, 0, 2)

	for providerName, providerModels := range providerModelsMap {
		if len(providerModels) == 0 {
			continue
		}

		providerModelStatus := make([]*model_provider.ProviderModelWithStatusEntity, 0, 10)
		firstModel := providerModels[0]

		for _, mapModel := range providerModels {
			providerModelStatus = append(providerModelStatus, &model_provider.ProviderModelWithStatusEntity{
				Status: mapModel.Status,
				ProviderModel: &model_provider.ProviderModel{
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
			Provider:  providerName,
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

	return providerResponses, nil
}

func (ms *ModelService) GetModelParameterRules(ctx context.Context, accountID string, provider string, model string) ([]*model_provider.ParameterRule, error) {
	tenantRecord, _, err := ms.AccountDomain.GetCurrentTenantOfAccount(ctx, accountID)

	if err != nil {
		return nil, err
	}

	providerConfigurations, err := ms.ModelProviderDomain.GetConfigurations(ctx, tenantRecord.ID)

	if err != nil {
		return nil, err
	}

	providerConfiguration, err := ms.ModelProviderDomain.GetConfigurationByProvider(ctx, providerConfigurations, provider)

	if err != nil {
		return nil, err
	}

	credentials, err := providerConfiguration.GetCurrentCredentials(base.LLM, model)

	if err != nil {
		return nil, err
	}

	model_provider, err := ms.ModelProviderDomain.ModelProviderRepo.GetProviderInstance(ctx, provider)

	if err != nil {
		return nil, err
	}

	modelInstance := model_provider.GetModelInstance(base.LLM)

	AIModelEntity, err := ms.ModelProviderDomain.GetModelSchema(ctx, model, credentials, modelInstance)

	if err != nil {
		return nil, err
	}

	return AIModelEntity.ParameterRules, nil

}

func (ms *ModelService) GetDefaultModelByType(ctx context.Context, accountID string, modelType string) (*dto.DefaultModelResponse, error) {

	tenantRecord, _, err := ms.AccountDomain.GetCurrentTenantOfAccount(ctx, accountID)

	if err != nil {
		return nil, err
	}
	defaultModelEntity, err := ms.ModelProviderDomain.GetDefaultModel(ctx, tenantRecord.ID, base.ModelType(modelType))

	if errors.IsCode(err, code.ErrDefaultModelNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &dto.DefaultModelResponse{
		Model:     defaultModelEntity.Model,
		ModelType: defaultModelEntity.ModelType,
		Provider: &model_provider.SimpleProviderEntity{
			Provider:            defaultModelEntity.Provider.Provider,
			Label:               defaultModelEntity.Provider.Label,
			IconSmall:           defaultModelEntity.Provider.IconSmall,
			IconLarge:           defaultModelEntity.Provider.IconLarge,
			SupportedModelTypes: defaultModelEntity.Provider.SupportedModelTypes,
		},
	}, nil

}
