package domain_service

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/core/tools"
	"github.com/lunarianss/Luna/internal/api-server/domain/agent/biz_entity"
)

type AgentDomain struct {
	*ToolTransformService
	*tools.ToolManager
}

func NewAgentDomain(ts *ToolTransformService, tm *tools.ToolManager) *AgentDomain {
	return &AgentDomain{
		ToolTransformService: ts,
		ToolManager:          tm,
	}
}

func (ad *AgentDomain) ListBuiltInTools(ctx context.Context, tenantID string) ([]*biz_entity.UserToolProvider, error) {

	var result []*biz_entity.UserToolProvider

	systemToolProviders, err := ad.ListBuiltInProviders()

	if err != nil {
		return nil, err
	}

	for _, systemToolProvider := range systemToolProviders {

		userProvider := ad.BuiltInProviderToUserProvider(systemToolProvider, nil, true)

		for _, tool := range systemToolProvider.Tools {
			userProvider.Tools = append(userProvider.Tools, ad.BuiltInToolToUserTool(tool, nil, tenantID, userProvider.Labels))
		}

		result = append(result, userProvider)
	}

	return result, nil
}
