package service

import (
	"context"

	domain "github.com/lunarianss/Luna/internal/api-server/_domain/account/domain_service"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/workspace"
)

type TenantService struct {
	accountDomain *domain.AccountDomain
}

func NewTenantService(accountDomain *domain.AccountDomain) *TenantService {
	return &TenantService{
		accountDomain: accountDomain,
	}
}

func (s *TenantService) GetTenantCurrentWorkspace(ctx context.Context, accountID string) (*dto.CurrentTenantInfo, error) {

	tenantRecord, accountJoinRecord, err := s.accountDomain.GetCurrentTenantOfAccount(ctx, accountID)

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

	accountRecord, err := s.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, err
	}

	tenants, err := s.accountDomain.TenantRepo.FindTenantsJoinByAccount(ctx, accountRecord)

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
