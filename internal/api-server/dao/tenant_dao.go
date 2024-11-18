package dao

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
	"gorm.io/gorm"
)

type TenantDao struct {
	db *gorm.DB
}

func NewTenantDao(db *gorm.DB) *TenantDao {
	return &TenantDao{db: db}
}

func (td *TenantDao) CreateOwnerTenant(ctx context.Context, tenant *model.Tenant, account *model.Account) (*model.Tenant, error) {
	if err := td.db.Create(tenant).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return tenant, nil
}
func (td *TenantDao) FindTenantJoinByAccount(ctx context.Context, account *model.Account) (*model.TenantAccountJoin, error) {
	var tenantJoin model.TenantAccountJoin
	if err := td.db.Find(&tenantJoin, "account_id = ?", account.ID).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &tenantJoin, nil
}

func (td *TenantDao) HasRoles(ctx context.Context, tenant *model.Tenant, roles []model.TenantAccountJoinRole) (bool, error) {
	var tenantMember []*model.TenantAccountJoin
	if err := td.db.Where("role IN ?", roles).Find(&tenantMember).Error; err != nil {
		return false, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return len(tenantMember) != 0, nil
}

func (td *TenantDao) GetTenantOfAccount(ctx context.Context, tenant *model.Tenant, account *model.Account) (*model.TenantAccountJoin, error) {

	var tenantAccountJoin model.TenantAccountJoin

	if err := td.db.Find(&tenantAccountJoin, "account_id = ? AND tenant_id = ?", account.ID, tenant.ID).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return &tenantAccountJoin, nil
}

func (td *TenantDao) UpdateRoleTenantOfAccount(ctx context.Context, ta *model.TenantAccountJoin) (*model.TenantAccountJoin, error) {

	if err := td.db.Model(ta).Where("id = ?", ta.ID).Update("role", ta.Role).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return ta, nil
}

func (td *TenantDao) CreateTenantOfAccount(ctx context.Context, tenant *model.Tenant, account *model.Account, role model.TenantAccountJoinRole) (*model.TenantAccountJoin, error) {
	var tenantAccountJoin = &model.TenantAccountJoin{
		AccountID: account.ID, TenantID: tenant.ID, Role: string(role),
	}

	if err := td.db.Create(tenantAccountJoin).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return tenantAccountJoin, nil
}