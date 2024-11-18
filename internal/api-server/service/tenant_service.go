package service

import (
	"context"

	domain "github.com/lunarianss/Luna/internal/api-server/domain/tenant"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
)

type TenantService struct {
	tenantDomain *domain.TenantDomain
}

func (ts *TenantService) CreateOwnerTenantIfNotExists(ctx context.Context, name string, account *model.Account, isSetup bool) error {

	return ts.tenantDomain.CreateOwnerTenantIfNotExists(ctx, name, account, false)
}
