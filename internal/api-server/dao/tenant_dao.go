package dao

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/internal/pkg/mysql"
	"github.com/lunarianss/Luna/pkg/errors"
	"gorm.io/gorm"
)

type TenantDao struct {
	db *gorm.DB
}

func NewTenantDao(db *gorm.DB) *TenantDao {
	return &TenantDao{db: db}
}

func (td *TenantDao) CreateOwnerTenant(ctx context.Context, tenant *model.Tenant, account *model.Account, isTransaction bool, tx *gorm.DB) (*model.Tenant, error) {

	var dbIns *gorm.DB

	if isTransaction && tx != nil {
		dbIns = tx
	} else {
		dbIns = td.db
	}

	if err := dbIns.Create(tenant).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return tenant, nil
}
func (td *TenantDao) FindTenantJoinByAccount(ctx context.Context, account *model.Account, isTransaction bool, tx *gorm.DB) (*model.TenantAccountJoin, error) {
	var tenantJoin model.TenantAccountJoin

	var dbIns *gorm.DB

	if isTransaction && tx != nil {

		dbIns = tx
	} else {
		dbIns = td.db
	}

	if err := dbIns.Find(&tenantJoin, "account_id = ?", account.ID).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &tenantJoin, nil
}

func (td *TenantDao) HasRoles(ctx context.Context, tenant *model.Tenant, roles []model.TenantAccountJoinRole, isTransaction bool, tx *gorm.DB) (bool, error) {
	var dbIns *gorm.DB

	if isTransaction && tx != nil {
		dbIns = tx
	} else {
		dbIns = td.db
	}

	var tenantMember []*model.TenantAccountJoin
	if err := dbIns.Where("role IN ?", roles).Find(&tenantMember).Error; err != nil {
		return false, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return len(tenantMember) != 0, nil
}

func (td *TenantDao) GetTenantOfAccount(ctx context.Context, tenant *model.Tenant, account *model.Account, isTransaction bool, tx *gorm.DB) (*model.TenantAccountJoin, error) {

	var dbIns *gorm.DB

	if isTransaction && tx != nil {
		dbIns = tx
	} else {
		dbIns = td.db
	}
	var tenantAccountJoin model.TenantAccountJoin

	if err := dbIns.Find(&tenantAccountJoin, "account_id = ? AND tenant_id = ?", account.ID, tenant.ID).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return &tenantAccountJoin, nil
}

func (td *TenantDao) UpdateRoleTenantOfAccount(ctx context.Context, ta *model.TenantAccountJoin, isTransaction bool, tx *gorm.DB) (*model.TenantAccountJoin, error) {
	var dbIns *gorm.DB

	if isTransaction && tx != nil {
		dbIns = tx
	} else {
		dbIns = td.db
	}

	if err := dbIns.Model(ta).Where("id = ?", ta.ID).Update("role", ta.Role).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return ta, nil
}

func (td *TenantDao) UpdateEncryptPublicKey(ctx context.Context, ta *model.Tenant, isTransaction bool, tx *gorm.DB) (*model.Tenant, error) {

	var dbIns *gorm.DB

	if isTransaction && tx != nil {
		dbIns = tx
	} else {
		dbIns = td.db
	}
	if err := dbIns.Model(ta).Where("id = ?", ta.ID).Update("encrypt_public_key", ta.EncryptPublicKey).Error; err != nil {
		return nil, err
	}

	return ta, nil
}

func (td *TenantDao) CreateTenantOfAccount(ctx context.Context, tenant *model.Tenant, account *model.Account, role model.TenantAccountJoinRole, isTransaction bool, tx *gorm.DB) (*model.TenantAccountJoin, error) {
	var dbIns *gorm.DB

	if isTransaction && tx != nil {
		dbIns = tx
	} else {
		dbIns = td.db
	}

	var tenantAccountJoin = &model.TenantAccountJoin{
		AccountID: account.ID, TenantID: tenant.ID, Role: string(role),
	}

	if err := dbIns.Create(tenantAccountJoin).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return tenantAccountJoin, nil
}

func (td *TenantDao) FindTenantMemberByAccount(ctx context.Context, account *model.Account) (*model.TenantAccountJoin, error) {

	var tenantAccountJoin model.TenantAccountJoin
	if err := td.db.Scopes(mysql.IDDesc()).Limit(1).Find(&tenantAccountJoin, "account_id = ?", account.ID).Error; err != nil {
		return nil, err
	}
	return &tenantAccountJoin, nil
}

func (td *TenantDao) FindCurrentTenantMemberByAccount(ctx context.Context, account *model.Account) (*model.TenantAccountJoin, error) {
	var tenantAccountJoin model.TenantAccountJoin
	if err := td.db.Scopes(mysql.IDDesc()).Limit(1).Find(&tenantAccountJoin, "account_id = ? AND current = ?", account.ID, 1).Error; err != nil {
		return nil, err
	}

	return &tenantAccountJoin, nil
}

func (td *TenantDao) UpdateCurrentTenantAccountJoin(ctx context.Context, ta *model.TenantAccountJoin) (*model.TenantAccountJoin, error) {
	if err := td.db.Model(ta).Update("current", 1).Error; err != nil {
		return nil, err
	}

	return ta, nil
}
