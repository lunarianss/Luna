package domain_service

import (
	"context"
	"sync"

	"github.com/lunarianss/Luna/internal/api-server/core/tools"
	"github.com/lunarianss/Luna/internal/api-server/domain/agent/biz_entity"
)

type AgentDomain struct {
	*ToolTransformService
	*tools.ToolManager
	*sync.RWMutex
}

func NewAgentDomain(ts *ToolTransformService, tm *tools.ToolManager) *AgentDomain {
	return &AgentDomain{
		ToolTransformService: ts,
		ToolManager:          tm,
		RWMutex:              &sync.RWMutex{},
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

func (ad *AgentDomain) ListBuiltInLabels(ctx context.Context) ([]*biz_entity.ToolLabel, error) {
	ad.RLock()
	defer ad.RUnlock()

	return biz_entity.GetDefaultTools(), nil
}
