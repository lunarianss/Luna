// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package app_model_config

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/lunarianss/Luna/infrastructure/errors"
	biz_entity_app_config "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/provider_app_config"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"
	biz_entity_provider_config "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_configuration"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
)

type ModelConfigManager struct {
	ProviderDomain *domain_service.ProviderDomain
}

func NewModelConfigManager(providerDomain *domain_service.ProviderDomain) *ModelConfigManager {
	return &ModelConfigManager{
		ProviderDomain: providerDomain,
	}
}

func (m *ModelConfigManager) ValidateAndSetDefaults(ctx context.Context, tenantID string, config *dto.AppModelConfigDto) (*dto.AppModelConfigDto, []string, error) {
	//todo dto层或者这里 使用 validate 库
	var (
		modelIDs        []string
		availableModels []*biz_entity_provider_config.ModelWithProvider
		isOk            bool
		modelMode       interface{}
		modelModeStr    string
	)

	providerConfigurations, err := m.ProviderDomain.GetConfigurations(ctx, tenantID)

	if err != nil {
		return nil, nil, err
	}

	availableModels, err = providerConfigurations.GetModels(ctx, config.Model.Provider, common.LLM, false)

	if err != nil {
		return nil, nil, err
	}

	if len(availableModels) == 0 {
		return nil, nil, errors.WithCode(code.ErrAllModelsEmpty, "models cannot be empty")
	}

	for _, availableModel := range availableModels {
		modelIDs = append(modelIDs, availableModel.Model)
	}

	if !slices.Contains(modelIDs, config.Model.Name) {
		return nil, nil, errors.WithCode(code.ErrRequiredCorrectModel, fmt.Sprintf("model %s not found in %s", config.Model.Name, strings.Join(modelIDs, ",")))
	}

	for _, availableModel := range availableModels {
		if availableModel.Model == config.Model.Name {
			if modelMode, isOk = availableModel.ModelProperties[common.MODE]; isOk {
				modelModeStr, _ = modelMode.(string)
			}
			break
		}	
	}

	if modelModeStr == "" {
		config.Model.Mode = "completion"
	} else {
		config.Model.Mode = modelModeStr
	}

	return config, []string{"model"}, nil
}

func (m *ModelConfigManager) Convert(ctx context.Context, config *dto.AppModelConfigDto) (*biz_entity_app_config.ModelConfigEntity, error) {

	return &biz_entity_app_config.ModelConfigEntity{
		Provider: config.Model.Provider,
		Model:    config.Model.Name,
		Mode:     config.Model.Mode,
	}, nil
}
