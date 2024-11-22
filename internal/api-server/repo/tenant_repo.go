package repo

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"gorm.io/gorm"
)

type TenantRepo interface {
	CreateOwnerTenant(ctx context.Context, tenant *model.Tenant, account *model.Account, tx *gorm.DB) (*model.Tenant, error)
	CreateCurrentTenantJoinOfAccount(ctx context.Context, tenant *model.Tenant, account *model.Account, role model.TenantAccountJoinRole, tx *gorm.DB) (*model.TenantAccountJoin, error)

	UpdateRoleTenantOfAccount(ctx context.Context, ta *model.TenantAccountJoin, tx *gorm.DB) (*model.TenantAccountJoin, error)
	UpdateCurrentTenantAccountJoin(ctx context.Context, ta *model.TenantAccountJoin) (*model.TenantAccountJoin, error)
	UpdateEncryptPublicKey(ctx context.Context, ta *model.Tenant, tx *gorm.DB) (*model.Tenant, error)

	FindTenantJoinByAccount(ctx context.Context, account *model.Account, tx *gorm.DB) (*model.TenantAccountJoin, error)
	FindTenantsJoinByAccount(ctx context.Context, account *model.Account) ([]*model.TenantJoinResult, error)
	GetCurrentTenantJoinByAccount(ctx context.Context, account *model.Account) (*model.TenantAccountJoin, error)
	GetTenantJoinOfAccount(ctx context.Context, tenant *model.Tenant, account *model.Account, tx *gorm.DB) (*model.TenantAccountJoin, error)
	GetTenantByID(ctx context.Context, ID string) (*model.Tenant, error)
	HasRoles(ctx context.Context, tenant *model.Tenant, roles []model.TenantAccountJoinRole, tx *gorm.DB) (bool, error)
}
