package tenant

import (
	"context"
	"fmt"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/repo"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
	"gorm.io/gorm"
)

type TenantDomain struct {
	TenantRepo repo.TenantRepo
}

func NewTenantDomain(tenantRepo repo.TenantRepo) *TenantDomain {
	return &TenantDomain{
		TenantRepo: tenantRepo,
	}
}

func (ts *TenantDomain) CreateOwnerTenantIfNotExistsTx(ctx context.Context, tx *gorm.DB, name string, account *model.Account, isSetup bool) error {
	tenantJoin, err := ts.TenantRepo.FindTenantJoinByAccount(ctx, account, true, tx)

	if err != nil {
		return err
	}

	if tenantJoin.ID != "" {
		return nil
	}

	if name == "" {
		name = fmt.Sprintf("%s's Workspace", account.Name)
	}

	// 补充 encrypt public key
	encryptPublicKey := ""
	tenant, err := ts.CreateTenantTx(ctx, tx, account, &model.Tenant{Name: name, EncryptPublicKey: encryptPublicKey}, false)

	if err != nil {
		return err
	}

	if _, err := ts.CreateTenantMemberTx(ctx, tx, account, tenant, string(model.OWNER)); err != nil {
		return err
	}
	return nil
}

func (ts *TenantDomain) CreateOwnerTenantIfNotExists(ctx context.Context, name string, account *model.Account, isSetup bool) error {
	tenantJoin, err := ts.TenantRepo.FindTenantJoinByAccount(ctx, account, false, nil)

	if err != nil {

		return err
	}

	if tenantJoin.ID != "" {
		return nil
	}

	if name == "" {
		name = fmt.Sprintf("%s's Workspace", account.Name)
	}

	// 补充 encrypt public key
	encryptPublicKey := ""
	tenant, err := ts.CreateTenant(ctx, account, &model.Tenant{Name: name, EncryptPublicKey: encryptPublicKey}, false)

	if err != nil {
		return err
	}

	if _, err := ts.CreateTenantMember(ctx, account, tenant, string(model.OWNER)); err != nil {
		return err
	}
	return nil
}

func (ts *TenantDomain) CreateTenant(ctx context.Context, account *model.Account, tenant *model.Tenant, isSetup bool) (*model.Tenant, error) {
	return ts.TenantRepo.CreateOwnerTenant(ctx, tenant, account, false, nil)
}

func (ts *TenantDomain) CreateTenantTx(ctx context.Context, tx *gorm.DB, account *model.Account, tenant *model.Tenant, isSetup bool) (*model.Tenant, error) {
	return ts.TenantRepo.CreateOwnerTenant(ctx, tenant, account, true, tx)
}

func (ts *TenantDomain) CreateTenantMemberTx(ctx context.Context, tx *gorm.DB, account *model.Account, tenant *model.Tenant, role string) (*model.TenantAccountJoin, error) {
	var (
		tenantAccountJoin *model.TenantAccountJoin
	)
	hasRole, err := ts.TenantRepo.HasRoles(ctx, tenant, []model.TenantAccountJoinRole{model.OWNER}, true, tx)

	if err != nil {
		return nil, err
	}

	if hasRole {
		return nil, errors.WithCode(code.ErrTenantAlreadyExist, "tenant %s already exists role %s", tenant.Name, role)
	}

	tenantMember, err := ts.TenantRepo.GetTenantOfAccount(ctx, tenant, account, true, tx)

	if err != nil {
		return nil, err
	}

	if tenantMember.ID != "" {
		tenantMember.Role = role
		tenantMember, err := ts.TenantRepo.UpdateRoleTenantOfAccount(ctx, tenantMember, true, tx)
		if err != nil {
			return nil, err
		}
		return tenantMember, nil
	}

	if tenantAccountJoin, err = ts.TenantRepo.CreateTenantOfAccount(ctx, tenant, account, model.TenantAccountJoinRole(role), true, tx); err != nil {
		return nil, err
	}

	return tenantAccountJoin, nil
}

func (ts *TenantDomain) CreateTenantMember(ctx context.Context, account *model.Account, tenant *model.Tenant, role string) (*model.TenantAccountJoin, error) {
	var (
		tenantAccountJoin *model.TenantAccountJoin
	)
	hasRole, err := ts.TenantRepo.HasRoles(ctx, tenant, []model.TenantAccountJoinRole{model.OWNER}, false, nil)

	if err != nil {
		return nil, err
	}

	if hasRole {
		return nil, errors.WithCode(code.ErrTenantAlreadyExist, "tenant %s already exists role %s", tenant.Name, role)
	}

	tenantMember, err := ts.TenantRepo.GetTenantOfAccount(ctx, tenant, account, false, nil)

	if err != nil {
		return nil, err
	}

	if tenantMember.ID != "" {
		tenantMember.Role = role
		tenantMember, err := ts.TenantRepo.UpdateRoleTenantOfAccount(ctx, tenantMember, false, nil)
		if err != nil {
			return nil, err
		}
		return tenantMember, nil
	}

	if tenantAccountJoin, err = ts.TenantRepo.CreateTenantOfAccount(ctx, tenant, account, model.TenantAccountJoinRole(role), false, nil); err != nil {
		return nil, err
	}

	return tenantAccountJoin, nil
}
