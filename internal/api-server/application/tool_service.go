package service

import (
	"context"

	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/agent/biz_entity"
	agentDomain "github.com/lunarianss/Luna/internal/api-server/domain/agent/domain_service"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
)

type ToolService struct {
	accountDomain *accountDomain.AccountDomain
	appDomain     *appDomain.AppDomain
	agentDomain   *agentDomain.AgentDomain
}

func NewToolService(accountDomain *accountDomain.AccountDomain, appDomain *appDomain.AppDomain, agentDomain *agentDomain.AgentDomain) *ToolService {
	return &ToolService{
		accountDomain: accountDomain,
		appDomain:     appDomain,
		agentDomain:   agentDomain,
	}
}

func (ts *ToolService) GetBuiltInTools(ctx context.Context, accountID string) ([]*biz_entity.UserToolProvider, error) {

	accountRecord, err := ts.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, err
	}

	tenant, _, err := ts.accountDomain.GetCurrentTenantOfAccount(ctx, accountRecord.ID)

	if err != nil {
		return nil, err
	}

	return ts.agentDomain.ListBuiltInTools(ctx, tenant.ID)
}
