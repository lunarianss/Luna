package app_agent_config

import (
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/provider_app_config"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
)

type AgentConfigManager struct {
}

func NewAgentConfigManager() *AgentConfigManager {
	return &AgentConfigManager{}
}

func (acm *AgentConfigManager) Convert(config *dto.AppModelConfigDto) (*biz_entity.AgentEntity, error) {
	var (
		strategy   string
		agentTools []*biz_entity.AgentToolEntity
	)

	agentStrategy := config.AgentMode.Strategy

	if agentStrategy == "function_call" {
		strategy = string(biz_entity.FUNCTION_CALLING)
	} else {
		strategy = string(biz_entity.CHAIN_OF_THOUGHT)
	}

	for _, agentTool := range config.AgentMode.Tools {
		agentTools = append(agentTools, &biz_entity.AgentToolEntity{
			ProviderType:   biz_entity.AgentToolProviderType(agentTool.ProviderType),
			ProviderID:     agentTool.ProviderID,
			ToolName:       agentTool.ToolName,
			ToolParameters: agentTool.ToolParameters,
		})
	}

	return &biz_entity.AgentEntity{
		Provider:     config.Model.Provider,
		Model:        config.Model.Name,
		Strategy:     biz_entity.Strategy(strategy),
		Tools:        agentTools,
		MaxIteration: 5,
		Prompt:       &biz_entity.AgentPromptEntity{},
	}, nil

}
