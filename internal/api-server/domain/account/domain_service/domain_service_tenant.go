// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package domain_service

import (
	"context"
	"fmt"

	"github.com/lunarianss/Luna/internal/api-server/domain/account/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/account/repository"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
	"github.com/lunarianss/Luna/infrastructure/errors"
	"gorm.io/gorm"
)

type TenantDomain struct {
	TenantRepo repository.TenantRepo
}

func NewTenantDomain(tenantRepo repository.TenantRepo) *TenantDomain {
	return &TenantDomain{
		TenantRepo: tenantRepo,
	}
}

func (ts *TenantDomain) CreateOwnerTenantIfNotExists(ctx context.Context, tx *gorm.DB, account *po_entity.Account, isSetup bool) error {
	tenantJoin, err := ts.TenantRepo.FindTenantJoinByAccount(ctx, account, tx)

	if err != nil {
		return err
	}

	if tenantJoin.ID != "" {
		return nil
	}

	name := fmt.Sprintf("%s's Workspace", account.Name)

	tenant, err := ts.CreateTenant(ctx, tx, account, &po_entity.Tenant{Name: name}, false)

	if err != nil {
		return err
	}

	encryptPublicKey, err := util.GenerateKeyPair(&util.FileStorage{}, tenant.ID)

	if err != nil {
		return errors.WithCode(code.ErrRSAGenerate, err.Error())
	}

	tenant.EncryptPublicKey = encryptPublicKey

	tenant, err = ts.TenantRepo.UpdateEncryptPublicKey(ctx, tenant, tx)

	if err != nil {
		return err
	}

	if _, err := ts.CreateTenantMember(ctx, tx, account, tenant, string(po_entity.OWNER)); err != nil {
		return err
	}

	return nil
}

func (ts *TenantDomain) CreateTenant(ctx context.Context, tx *gorm.DB, account *po_entity.Account, tenant *po_entity.Tenant, isSetup bool) (*po_entity.Tenant, error) {
	return ts.TenantRepo.CreateOwnerTenant(ctx, tenant, account, tx)
}

func (ts *TenantDomain) CreateTenantMember(ctx context.Context, tx *gorm.DB, account *po_entity.Account, tenant *po_entity.Tenant, role string) (*po_entity.TenantAccountJoin, error) {
	var (
		tenantAccountJoin *po_entity.TenantAccountJoin
	)
	hasRole, err := ts.TenantRepo.HasRoles(ctx, tenant, []po_entity.TenantAccountJoinRole{po_entity.OWNER}, tx)

	if err != nil {
		return nil, err
	}

	if hasRole {
		return nil, errors.WithCode(code.ErrTenantAlreadyExist, "tenant %s already exists role %s", tenant.Name, role)
	}

	tenantMember, err := ts.TenantRepo.GetTenantJoinOfAccount(ctx, tenant, account, tx)

	if err != nil {
		return nil, err
	}

	if tenantMember.ID != "" {
		tenantMember.Role = role
		tenantMember, err := ts.TenantRepo.UpdateRoleTenantOfAccount(ctx, tenantMember, tx)
		if err != nil {
			return nil, err
		}
		return tenantMember, nil
	}

	if tenantAccountJoin, err = ts.TenantRepo.CreateCurrentTenantJoinOfAccount(ctx, tenant, account, po_entity.TenantAccountJoinRole(role), tx); err != nil {
		return nil, err
	}

	return tenantAccountJoin, nil
}
