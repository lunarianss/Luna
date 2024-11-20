package service

import (
	"context"

	domain "github.com/lunarianss/Luna/internal/api-server/domain/account"
	tenantDomain "github.com/lunarianss/Luna/internal/api-server/domain/tenant"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/workspace"
)

type TenantService struct {
	AccountDomain *domain.AccountDomain
	TenantDomain  *tenantDomain.TenantDomain
}

func NewTenantService(accountDomain *domain.AccountDomain, tenantDomain *tenantDomain.TenantDomain) *TenantService {
	return &TenantService{
		AccountDomain: accountDomain,
		TenantDomain:  tenantDomain,
	}
}

func (s *TenantService) GetTenantCurrentWorkspace(ctx context.Context, accountID string) (*dto.CurrentTenantInfo, error) {

	tenantRecord, accountJoinRecord, err := s.AccountDomain.GetCurrentTenantOfAccount(ctx, accountID)

	if err != nil {
		return nil, err
	}

	return &dto.CurrentTenantInfo{
		ID:           tenantRecord.ID,
		Name:         tenantRecord.Name,
		Plan:         tenantRecord.Plan,
		Status:       tenantRecord.Status,
		CreateAt:     tenantRecord.CreatedAt,
		InTrail:      true,
		Role:         accountJoinRecord.Role,
		CustomConfig: tenantRecord.CustomConfig,
	}, nil

}

func (s *TenantService) GetJoinTenants(ctx context.Context, accountID string) ([]*dto.CurrentTenantInfo, error) {

	accountRecord, err := s.AccountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, err
	}

	tenants, err := s.TenantDomain.TenantRepo.FindTenantsJoinByAccount(ctx, accountRecord)

	if err != nil {
		return nil, err
	}

	tenantsInfo := make([]*dto.CurrentTenantInfo, 0)

	for _, tenantJoinResult := range tenants {
		tenantsInfo = append(tenantsInfo, &dto.CurrentTenantInfo{
			ID:       tenantJoinResult.ID,
			Name:     tenantJoinResult.Name,
			Plan:     tenantJoinResult.Plan,
			Status:   tenantJoinResult.Status,
			CreateAt: tenantJoinResult.CreatedAt,
			InTrail:  true,
			Role:     tenantJoinResult.Role,
		})
	}

	return tenantsInfo, nil
}
