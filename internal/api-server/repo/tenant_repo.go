package repo

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"gorm.io/gorm"
)

type TenantRepo interface {
	CreateOwnerTenant(ctx context.Context, tenant *model.Tenant, account *model.Account, isTransaction bool, tx *gorm.DB) (*model.Tenant, error)
	FindTenantJoinByAccount(ctx context.Context, account *model.Account, isTransaction bool, tx *gorm.DB) (*model.TenantAccountJoin, error)
	HasRoles(ctx context.Context, tenant *model.Tenant, roles []model.TenantAccountJoinRole, isTransaction bool, tx *gorm.DB) (bool, error)
	GetTenantOfAccount(ctx context.Context, tenant *model.Tenant, account *model.Account, isTransaction bool, tx *gorm.DB) (*model.TenantAccountJoin, error)
	CreateTenantOfAccount(ctx context.Context, tenant *model.Tenant, account *model.Account, role model.TenantAccountJoinRole, isTransaction bool, tx *gorm.DB) (*model.TenantAccountJoin, error)
	UpdateRoleTenantOfAccount(ctx context.Context, ta *model.TenantAccountJoin, isTransaction bool, tx *gorm.DB) (*model.TenantAccountJoin, error)
}
