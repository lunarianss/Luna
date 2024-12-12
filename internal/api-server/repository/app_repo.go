// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repo_impl

import (
	"context"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/repository"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/lunarianss/Luna/internal/infrastructure/mysql"
	"gorm.io/gorm"
)

type AppRepoImpl struct {
	db *gorm.DB
}

var _ repository.AppRepo = (*AppRepoImpl)(nil)

func NewAppRepoImpl(db *gorm.DB) *AppRepoImpl {
	return &AppRepoImpl{db}
}

func (ad *AppRepoImpl) CreateApp(ctx context.Context, tx *gorm.DB, app *po_entity.App) (*po_entity.App, error) {
	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = ad.db
	}

	if err := dbIns.Create(app).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return app, nil
}

func (ad *AppRepoImpl) CreateAppConfig(ctx context.Context, tx *gorm.DB, appConfig *po_entity.AppModelConfig) (*po_entity.AppModelConfig, error) {
	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = ad.db
	}

	if err := dbIns.Create(appConfig).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return appConfig, nil
}

func (ad *AppRepoImpl) UpdateAppConfigID(ctx context.Context, app *po_entity.App) error {
	if err := ad.db.Model(app).Where("id = ?", app.ID).Update("app_model_config_id", app.AppModelConfigID).Error; err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil
}

func (ad *AppRepoImpl) FindTenantApps(ctx context.Context, tenantID string, page, pageSize int) ([]*po_entity.App, int64, error) {
	var apps []*po_entity.App
	var appCount int64

	if err := ad.db.Model(&po_entity.App{}).Scopes(mysql.Paginate(page, pageSize)).Find(&apps, "tenant_id = ? AND is_universal = ?", tenantID, 0).Count(&appCount).Error; err != nil {
		return nil, 0, err
	}
	return apps, appCount, nil
}

func (ad *AppRepoImpl) CreateAppWithConfig(ctx context.Context, tx *gorm.DB, app *po_entity.App, appConfig *po_entity.AppModelConfig) (*po_entity.App, error) {
	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = ad.db
	}

	if err := dbIns.Create(app).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	appConfig.AppID = app.ID

	if err := dbIns.Create(appConfig).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	app.AppModelConfigID = appConfig.ID

	if err := dbIns.Model(app).Update("app_model_config_id", appConfig.ID).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return app, nil
}

func (ad *AppRepoImpl) UpdateEnableAppSite(ctx context.Context, app *po_entity.App) (*po_entity.App, error) {
	if err := ad.db.Model(&po_entity.App{}).Select("enable_site", "updated_by", "updated_at").Where("id = ?", app.ID).Updates(app).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return app, nil
}

func (ad *AppRepoImpl) GetAppByID(ctx context.Context, appID string) (*po_entity.App, error) {
	var app po_entity.App

	if err := ad.db.First(&app, "id = ?", appID).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &app, nil
}

func (ad *AppRepoImpl) GetAppModelConfigById(ctx context.Context, appConfigID, appID string) (*po_entity.AppModelConfig, error) {
	var appConfig po_entity.AppModelConfig
	if err := ad.db.First(&appConfig, "id = ? AND app_id = ?", appConfigID, appID).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &appConfig, nil
}

func (ad *AppRepoImpl) GetTenantApp(ctx context.Context, appID, tenantID string) (*po_entity.App, error) {
	var app po_entity.App
	if err := ad.db.First(&app, "id = ? AND tenant_id = ? AND status = ?", appID, tenantID, "normal").Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &app, nil
}
