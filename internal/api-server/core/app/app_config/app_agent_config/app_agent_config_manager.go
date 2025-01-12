// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package app_agent_config

import (
	"context"
	"encoding/json"

	"github.com/lunarianss/Luna/infrastructure/errors"
	assembler "github.com/lunarianss/Luna/internal/api-server/assembler/app"
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config/app_model_config"
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config/app_prompt_template"
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config/app_variable_config"
	biz_entity_app_config "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/provider_app_config"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	po_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
)

type ChatAgentAppConfigManager struct {
	ProviderDomain *domain_service.ProviderDomain
}

func NewAgentChatAppConfigManager(providerDomain *domain_service.ProviderDomain) *ChatAgentAppConfigManager {
	return &ChatAgentAppConfigManager{
		ProviderDomain: providerDomain,
	}
}

func (m *ChatAgentAppConfigManager) ConfigValidate(ctx context.Context, tenantID string, config *dto.AppModelConfigDto) (*dto.AppModelConfigDto, error) {

	// model
	modelConfigManager := app_model_config.NewModelConfigManager(m.ProviderDomain)

	config, _, err := modelConfigManager.ValidateAndSetDefaults(ctx, tenantID, config)

	if err != nil {
		return nil, err
	}
	return config, nil
}

func (m *ChatAgentAppConfigManager) GetAppConfig(ctx context.Context, appModel *po_entity.App, appModelConfig *po_entity.AppModelConfig, conversation *po_entity_chat.Conversation, overrideConfigDict *dto.AppModelConfigDto) (*biz_entity_app_config.AgentChatAppConfig, error) {

	var (
		configFrom biz_entity_app_config.EasyUIBasedAppModelConfigFrom
		configDict *dto.AppModelConfigDto
	)

	if overrideConfigDict != nil {
		configFrom = biz_entity_app_config.Args
	} else if conversation != nil {
		configFrom = biz_entity_app_config.ConversationSpecificConfig
	} else {
		configFrom = biz_entity_app_config.AppLatestConfig
	}

	if configFrom != biz_entity_app_config.Args {
		appModelByte, err := json.Marshal(appModelConfig)

		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(appModelByte, &configDict); err != nil {
			return nil, err
		}
	} else {
		if overrideConfigDict == nil {
			return nil, errors.WithCode(code.ErrRequiredOverrideConfig, "override_config_dict is required when config_from is ARGS")
		}

		configDict = overrideConfigDict
	}

	modelConfigManager := app_model_config.NewModelConfigManager(m.ProviderDomain)
	modelConfigEntity, err := modelConfigManager.Convert(ctx, configDict)

	if err != nil {
		return nil, err
	}

	promptTemplateConfigManager := app_prompt_template.NewPromptTemplateConfigManager()

	promptTemplate, err := promptTemplateConfigManager.Convert(configDict)

	if err != nil {
		return nil, err
	}

	variables := app_variable_config.NewBasicVariablesConfigManager()

	agentConfig, err := NewAgentConfigManager().Convert(configDict)

	if err != nil {
		return nil, err
	}
	return &biz_entity_app_config.AgentChatAppConfig{
		EasyUIBasedAppConfig: &biz_entity_app_config.EasyUIBasedAppConfig{
			AppConfig: &biz_entity_app_config.AppConfig{
				TenantID:               appModel.TenantID,
				AppID:                  appModel.ID,
				AppMode:                appModel.Mode,
				SensitiveWordAvoidance: nil,
				AdditionalFeatures:     nil,
				Variables:              variables.Convert(configDict),
			},
			AppModelConfig:     assembler.ConvertToConfigEntity(configDict),
			AppModelConfigFrom: configFrom,
			AppModelConfigID:   appModelConfig.ID,
			Model:              modelConfigEntity,
			PromptTemplate:     promptTemplate,
		},
		AgentEntity: agentConfig,
	}, nil
}
