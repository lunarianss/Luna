package model_config

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"
	biz_entity_provider_config "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_configuration"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime/model_providers"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
)

type ModelConfigManager struct {
	ProviderDomain *domain_service.ProviderDomain
}

func NewModelConfigManager(providerDomain *domain_service.ProviderDomain) *ModelConfigManager {
	return &ModelConfigManager{
		ProviderDomain: providerDomain,
	}
}

func (m *ModelConfigManager) ValidateAndSetDefaults(ctx context.Context, tenantID string, config map[string]any) (map[string]interface{}, []string, error) {
	//todo  以下代码使用 validate 库进行简化
	var (
		providerNames   []string
		modelIDs        []string
		modelConfig     map[string]interface{}
		modelNameStr    string
		modelName       interface{}
		availableModels []*biz_entity_provider_config.ModelWithProvider
		isOk            bool
		providerName    interface{}
		providerNameStr string
		modelMode       interface{}
		modelModeStr    string
	)

	if _, isOk = config["model"]; !isOk {
		return nil, nil, errors.WithCode(code.ErrModelEmptyInConfig, fmt.Sprintf("model field not found in config %v", config))
	}

	if modelConfig, isOk = config["model"].(map[string]any); !isOk {
		return nil, nil, errors.WithCode(code.ErrModelEmptyInConfig, fmt.Sprintf("model field is empty json in config %v", config))
	}

	providerEntities, err := model_providers.Factory.GetProvidersFromDir()

	if err != nil {
		return nil, nil, err
	}

	for _, providerEntity := range providerEntities {
		providerNames = append(providerNames, providerEntity.Provider)
	}

	if providerName, isOk = modelConfig["provider"]; !isOk {
		return nil, nil, errors.WithCode(code.ErrRequiredCorrectProvider, "provider is required")
	}

	if providerNameStr, isOk = providerName.(string); isOk {
		if !slices.Contains(providerNames, providerNameStr) {
			return nil, nil, errors.WithCode(code.ErrRequiredCorrectProvider, fmt.Sprintf("provider %s must be include in %s", providerNameStr, strings.Join(providerNames, ",")))
		}
	}

	if modelName, isOk = modelConfig["name"]; !isOk {
		return nil, nil, errors.WithCode(code.ErrRequiredModelName, "model name is required")
	}

	if modelNameStr, isOk = modelName.(string); !isOk {
		return nil, nil, errors.WithCode(code.ErrRequiredModelName, "model name is not string")
	}

	providerConfigurations, err := m.ProviderDomain.GetConfigurations(ctx, tenantID)

	if err != nil {
		return nil, nil, err
	}

	for _, providerConfiguration := range providerConfigurations.Configurations {

		if providerConfiguration.Provider.Provider != providerNameStr {
			continue
		}
		availableModels, err = providerConfiguration.GetProviderModels(ctx, common.LLM, false)
		if err != nil {
			return nil, nil, err
		}
	}

	if len(availableModels) == 0 {
		return nil, nil, errors.WithCode(code.ErrAllModelsEmpty, "models cannot be empty")
	}

	for _, availableModel := range availableModels {
		modelIDs = append(modelIDs, availableModel.Model)
	}

	if !slices.Contains(modelIDs, modelNameStr) {
		return nil, nil, errors.WithCode(code.ErrRequiredCorrectModel, fmt.Sprintf("model %s not found in %s", modelNameStr, strings.Join(modelIDs, ",")))
	}

	for _, availableModel := range availableModels {
		if availableModel.Model == modelNameStr {
			if modelMode, isOk = availableModel.ModelProperties[common.MODE]; isOk {
				modelModeStr, _ = modelMode.(string)
			}
			break
		}
	}

	if modelModeStr == "" {
		modelConfig["mode"] = "completion"
	} else {
		modelConfig["mode"] = modelModeStr
	}

	// todo validate and default to completion params

	// override
	config["model"] = modelConfig

	return config, []string{"model"}, nil
}

func (m *ModelConfigManager) Convert(ctx context.Context, config map[string]interface{}) (*app_config.ModelConfigEntity, error) {
	var (
		modelConfig map[string]interface{}
	)

	modelConfig, ok := config["model"].(map[string]interface{})

	if !ok {
		return nil, errors.WithCode(code.ErrModelEmptyInConfig, "empty model configuration in config")
	}

	return &app_config.ModelConfigEntity{
		Provider: modelConfig["provider"].(string),
		Model:    modelConfig["name"].(string),
		Mode:     modelConfig["mode"].(string),
	}, nil
}
