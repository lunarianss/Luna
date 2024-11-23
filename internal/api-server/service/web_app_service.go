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

type WebAppService struct {
	appRunningDomain *domain.AppRunningDomain
	accountDomain    *accountDomain.AccountDomain
	appDomain        *appDomain.AppDomain
	config           *config.Config
}

func NewWebAppService(appRunningDomain *domain.AppRunningDomain, accountDomain *accountDomain.AccountDomain, appDomain *appDomain.AppDomain, config *config.Config) *WebAppService {
	return &WebAppService{
		appRunningDomain: appRunningDomain,
		accountDomain:    accountDomain,
		appDomain:        appDomain,
		config:           config,
	}
}

func (s *WebAppService) GetWebAppParameters(ctx context.Context, appID string) (*dto.WebAppParameterResponse, error) {
	appRecord, err := s.appDomain.AppRepo.GetAppByID(ctx, appID)

	if err != nil {
		return nil, err
	}

	appConfigRecord, err := s.appDomain.AppRepo.GetAppModelConfigByAppID(ctx, appRecord.ID)

	if err != nil {
		return nil, err
	}

	return dto.AppConfigRecordToParameter(appConfigRecord), nil
}

func (s *WebAppService) GetWebAppMeta(ctx context.Context, appID string, endUserID string, appCode string) (*dto.GetWebSiteResponse, error) {
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
