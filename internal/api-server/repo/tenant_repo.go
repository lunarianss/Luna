package repo

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
)

type TenantRepo interface {
	CreateOwnerTenant(ctx context.Context, tenant *model.Tenant, account *model.Account) (*model.Tenant, error)
	FindTenantJoinByAccount(ctx context.Context, account *model.Account) (*model.TenantAccountJoin, error)
	HasRoles(ctx context.Context, tenant *model.Tenant, roles []model.TenantAccountJoinRole) (bool, error)
	GetTenantOfAccount(ctx context.Context, tenant *model.Tenant, account *model.Account) (*model.TenantAccountJoin, error)
	CreateTenantOfAccount(ctx context.Context, tenant *model.Tenant, account *model.Account, role model.TenantAccountJoinRole) (*model.TenantAccountJoin, error)
	UpdateRoleTenantOfAccount(ctx context.Context, ta *model.TenantAccountJoin) (*model.TenantAccountJoin, error)
}
