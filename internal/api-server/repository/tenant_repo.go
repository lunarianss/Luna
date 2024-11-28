// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repo_impl

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/domain/account/entity/po_entity"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/internal/pkg/mysql"
	"github.com/lunarianss/Luna/pkg/errors"
	"gorm.io/gorm"
)

type TenantRepoImpl struct {
	db *gorm.DB
}

func NewTenantRepoImpl(db *gorm.DB) *TenantRepoImpl {
	return &TenantRepoImpl{db: db}
}

func (td *TenantRepoImpl) CreateOwnerTenant(ctx context.Context, tenant *po_entity.Tenant, account *po_entity.Account, tx *gorm.DB) (*po_entity.Tenant, error) {

	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = td.db
	}

	if err := dbIns.Create(tenant).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return tenant, nil
}

func (td *TenantRepoImpl) FindTenantJoinByAccount(ctx context.Context, account *po_entity.Account, tx *gorm.DB) (*po_entity.TenantAccountJoin, error) {
	var tenantJoin po_entity.TenantAccountJoin

	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = td.db
	}

	if err := dbIns.Scopes(mysql.IDDesc()).Limit(1).Find(&tenantJoin, "account_id = ?", account.ID).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &tenantJoin, nil
}

func (td *TenantRepoImpl) GetTenantByID(ctx context.Context, ID string) (*po_entity.Tenant, error) {
	var tenant po_entity.Tenant

	if err := td.db.First(&tenant, "id = ?", ID).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &tenant, nil
}

func (td *TenantRepoImpl) HasRoles(ctx context.Context, tenant *po_entity.Tenant, roles []po_entity.TenantAccountJoinRole, tx *gorm.DB) (bool, error) {
	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = td.db
	}

	var tenantMember []*po_entity.TenantAccountJoin
	if err := dbIns.Where("role IN ?", roles).Find(&tenantMember).Error; err != nil {
		return false, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return len(tenantMember) != 0, nil
}

func (td *TenantRepoImpl) GetCurrentTenantJoinByAccount(ctx context.Context, account *po_entity.Account) (*po_entity.TenantAccountJoin, error) {
	var tenantAccountJoin po_entity.TenantAccountJoin
	if err := td.db.First(&tenantAccountJoin, "account_id = ? AND current = ?", account.ID, 1).Error; err != nil {
		return nil, err
	}
	return &tenantAccountJoin, nil
}

func (td *TenantRepoImpl) FindTenantsJoinByAccount(ctx context.Context, account *po_entity.Account) ([]*po_entity.TenantJoinResult, error) {
	var tenantJoinResults []*po_entity.TenantJoinResult

	if err := td.db.Table("tenants").Select("tenants.id as tenant_id, tenant_account_joins.role as tenant_join_role, tenants.name as tenant_name, tenants.plan as tenant_plan, tenants.status as tenant_status, tenants.created_at as tenant_created_at, tenants.updated_at as tenant_updated_at, tenants.custom_config as tenants_custom_config").Joins("join tenant_account_joins on tenants.id = tenant_account_joins.tenant_id").Scan(&tenantJoinResults).Error; err != nil {
		return nil, err
	}

	return tenantJoinResults, nil
}

func (td *TenantRepoImpl) GetTenantJoinOfAccount(ctx context.Context, tenant *po_entity.Tenant, account *po_entity.Account, tx *gorm.DB) (*po_entity.TenantAccountJoin, error) {

	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = td.db
	}

	var tenantAccountJoin po_entity.TenantAccountJoin

	if err := dbIns.Find(&tenantAccountJoin, "account_id = ? AND tenant_id = ?", account.ID, tenant.ID).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return &tenantAccountJoin, nil
}

func (td *TenantRepoImpl) UpdateRoleTenantOfAccount(ctx context.Context, ta *po_entity.TenantAccountJoin, tx *gorm.DB) (*po_entity.TenantAccountJoin, error) {
	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = td.db
	}

	if err := dbIns.Model(ta).Where("id = ?", ta.ID).Update("role", ta.Role).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return ta, nil
}

func (td *TenantRepoImpl) UpdateEncryptPublicKey(ctx context.Context, ta *po_entity.Tenant, tx *gorm.DB) (*po_entity.Tenant, error) {

	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = td.db
	}

	if err := dbIns.Model(ta).Where("id = ?", ta.ID).Update("encrypt_public_key", ta.EncryptPublicKey).Error; err != nil {
		return nil, err
	}

	return ta, nil
}

func (td *TenantRepoImpl) CreateCurrentTenantJoinOfAccount(ctx context.Context, tenant *po_entity.Tenant, account *po_entity.Account, role po_entity.TenantAccountJoinRole, tx *gorm.DB) (*po_entity.TenantAccountJoin, error) {
	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = td.db
	}

	var tenantAccountJoin = &po_entity.TenantAccountJoin{
		AccountID: account.ID, TenantID: tenant.ID, Role: string(role), Current: 1,
	}

	if err := dbIns.Create(tenantAccountJoin).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return tenantAccountJoin, nil
}

func (td *TenantRepoImpl) UpdateCurrentTenantAccountJoin(ctx context.Context, ta *po_entity.TenantAccountJoin) (*po_entity.TenantAccountJoin, error) {
	if err := td.db.Model(ta).Update("current", 1).Error; err != nil {
		return nil, err
	}
	return ta, nil
}
