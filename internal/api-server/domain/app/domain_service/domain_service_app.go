// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package domain_service

import (
	"context"
	"fmt"

	"github.com/lunarianss/Luna/infrastructure/errors"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/provider_app_config"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/repository"
	web_app_entity "github.com/lunarianss/Luna/internal/api-server/domain/web_app/entity/po_entity"
	web_app_repo "github.com/lunarianss/Luna/internal/api-server/domain/web_app/repository"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"gorm.io/gorm"
)

type AppDomain struct {
	AppRepo    repository.AppRepo
	WebAppRepo web_app_repo.WebAppRepo
	db         *gorm.DB
}

func NewAppDomain(appRepo repository.AppRepo, webAppRepo web_app_repo.WebAppRepo, db *gorm.DB) *AppDomain {
	return &AppDomain{
		AppRepo:    appRepo,
		WebAppRepo: webAppRepo,

		db: db,
	}
}

func (ad *AppDomain) GetTemplate(ctx context.Context, mode string) (*biz_entity.AppTemplate, error) {
	appTemplate, ok := biz_entity.DefaultAppTemplates[biz_entity.AppMode(mode)]

	if !ok {
		return nil, errors.WithCode(code.ErrAppMapMode, fmt.Sprintf("Invalid node template: %v", mode))
	}
	return &appTemplate, nil
}

func (ad *AppDomain) CreateApp(ctx context.Context, app *po_entity.App, appConfig *po_entity.AppModelConfig, provider, modelName, accountDefaultLanguage string) (*po_entity.App, *po_entity.AppModelConfig, error) {
	tx := ad.db.Begin()
	var err error

	if provider != "" && modelName != "" {
		app, err = ad.AppRepo.CreateAppWithConfig(ctx, tx, app, appConfig)
		if err != nil {
			tx.Rollback()
			return nil, nil, err
		}
	} else {
		app, err = ad.AppRepo.CreateApp(ctx, tx, app)
		if err != nil {
			tx.Rollback()
			return nil, nil, err
		}
	}

	installApp := &web_app_entity.InstalledApp{
		TenantID:         app.TenantID,
		AppID:            app.ID,
		AppOwnerTenantID: app.TenantID,
	}

	if _, err := ad.WebAppRepo.CreateInstallApp(ctx, installApp, tx); err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	siteCode, err := ad.WebAppRepo.GenerateUniqueCodeForSite(ctx)

	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	site := &web_app_entity.Site{
		AppID:                  app.ID,
		Title:                  app.Name,
		IconType:               app.IconType,
		Icon:                   app.Icon,
		IconBackground:         app.IconBackground,
		DefaultLanguage:        accountDefaultLanguage,
		CustomizeTokenStrategy: "not_allowed",
		Code:                   siteCode,
		CreatedBy:              app.CreatedBy,
		UpdatedBy:              app.UpdatedBy,
	}

	if _, err := ad.WebAppRepo.CreateSite(ctx, site, tx); err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, nil, err
	}

	return app, appConfig, nil
}
