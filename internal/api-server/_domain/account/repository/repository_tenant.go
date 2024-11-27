package repository

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/_domain/account/entity/po_entity"
	"gorm.io/gorm"
)

type TenantRepo interface {
	CreateOwnerTenant(ctx context.Context, tenant *po_entity.Tenant, account *po_entity.Account, tx *gorm.DB) (*po_entity.Tenant, error)
	CreateCurrentTenantJoinOfAccount(ctx context.Context, tenant *po_entity.Tenant, account *po_entity.Account, role po_entity.TenantAccountJoinRole, tx *gorm.DB) (*po_entity.TenantAccountJoin, error)

	UpdateRoleTenantOfAccount(ctx context.Context, ta *po_entity.TenantAccountJoin, tx *gorm.DB) (*po_entity.TenantAccountJoin, error)
	UpdateCurrentTenantAccountJoin(ctx context.Context, ta *po_entity.TenantAccountJoin) (*po_entity.TenantAccountJoin, error)
	UpdateEncryptPublicKey(ctx context.Context, ta *po_entity.Tenant, tx *gorm.DB) (*po_entity.Tenant, error)

	FindTenantJoinByAccount(ctx context.Context, account *po_entity.Account, tx *gorm.DB) (*po_entity.TenantAccountJoin, error)
	FindTenantsJoinByAccount(ctx context.Context, account *po_entity.Account) ([]*po_entity.TenantJoinResult, error)
	GetCurrentTenantJoinByAccount(ctx context.Context, account *po_entity.Account) (*po_entity.TenantAccountJoin, error)
	GetTenantJoinOfAccount(ctx context.Context, tenant *po_entity.Tenant, account *po_entity.Account, tx *gorm.DB) (*po_entity.TenantAccountJoin, error)
	GetTenantByID(ctx context.Context, ID string) (*po_entity.Tenant, error)
	HasRoles(ctx context.Context, tenant *po_entity.Tenant, roles []po_entity.TenantAccountJoinRole, tx *gorm.DB) (bool, error)
}
