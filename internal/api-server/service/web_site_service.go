package service

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/config"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/app_running"
	siteDto "github.com/lunarianss/Luna/internal/api-server/dto/app"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/web_app"
)

type WebSiteService struct {
	appRunningDomain *domain.AppRunningDomain
	accountDomain    *accountDomain.AccountDomain
	appDomain        *appDomain.AppDomain
	config           *config.Config
}

func NewWebSiteService(appRunningDomain *domain.AppRunningDomain, accountDomain *accountDomain.AccountDomain, appDomain *appDomain.AppDomain, config *config.Config) *WebSiteService {
	return &WebSiteService{
		appRunningDomain: appRunningDomain,
		accountDomain:    accountDomain,
		appDomain:        appDomain,
		config:           config,
	}
}

func (s *WebSiteService) GetSiteByWebToken(ctx context.Context, appID string, endUserID string, appCode string) (*dto.GetWebSiteResponse, error) {
	siteRecord, err := s.appRunningDomain.AppRunningRepo.GetSiteByCode(ctx, appCode)

	if err != nil {
		return nil, err
	}

	appRecord, err := s.appDomain.AppRepo.GetAppByID(ctx, appID)

	if err != nil {
		return nil, err
	}

	tenantRecord, err := s.accountDomain.TenantRepo.GetTenantByID(ctx, appRecord.TenantID)

	if err != nil {
		return nil, err
	}

	return &dto.GetWebSiteResponse{
		AppID:      appID,
		EndUserID:  endUserID,
		EnableSite: int(appRecord.EnableSite),
		Plan:       tenantRecord.Plan,
		Site:       siteDto.SiteRecordToSiteDetail(siteRecord, s.config),
	}, nil
}
