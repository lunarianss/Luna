package domain_service

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/domain/agent/biz_entity"
)

type AgentDomain struct {
}

func (*AgentDomain) ListBuiltInTools(ctx context.Context) ([]*biz_entity.UserToolProvider, error) {

	return nil, nil
}
