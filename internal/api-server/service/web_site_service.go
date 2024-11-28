// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/config"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	webAppDomain "github.com/lunarianss/Luna/internal/api-server/domain/web_app/domain_service"
	siteDto "github.com/lunarianss/Luna/internal/api-server/dto/app"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/web_app"
)

type WebSiteService struct {
	webAppDomain  *webAppDomain.WebAppDomain
	accountDomain *accountDomain.AccountDomain
	appDomain     *domain_service.AppDomain
	config        *config.Config
}

func NewWebSiteService(webAppDomain *webAppDomain.WebAppDomain, accountDomain *accountDomain.AccountDomain, appDomain *domain_service.AppDomain, config *config.Config) *WebSiteService {
	return &WebSiteService{
		webAppDomain:  webAppDomain,
		accountDomain: accountDomain,
		appDomain:     appDomain,
		config:        config,
	}
}

func (s *WebSiteService) GetSiteByWebToken(ctx context.Context, appID string, endUserID string, appCode string) (*dto.GetWebSiteResponse, error) {
	siteRecord, err := s.webAppDomain.WebAppRepo.GetSiteByCode(ctx, appCode)

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
