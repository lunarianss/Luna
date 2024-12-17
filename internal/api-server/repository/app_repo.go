// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repo_impl

import (
	"context"
	"fmt"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/repository"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/lunarianss/Luna/internal/infrastructure/mysql"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
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

func (ad *AppRepoImpl) CreateServiceToken(ctx context.Context, token *po_entity.ApiToken) (*po_entity.ApiToken, error) {
	if err := ad.db.Create(token).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return token, nil
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

func (ad *AppRepoImpl) GetServiceTokenCount(ctx context.Context, appID string) (int64, error) {
	var count int64
	if err := ad.db.Model(&po_entity.ApiToken{}).Count(&count).Where("type = ? AND app_id = ?", "app", appID).Error; err != nil {
		return count, err
	}

	return count, nil
}

func (ad *AppRepoImpl) GetServiceTokenByCode(ctx context.Context, token string) (*po_entity.ApiToken, error) {

	var appToken po_entity.ApiToken

	if err := ad.db.First(&appToken, "token = ?", token).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("token %s record not found", token))
	}

	return &appToken, nil
}

func (ad *AppRepoImpl) GenerateServiceToken(ctx context.Context, num int) (string, error) {

	token, err := util.GenerateRandomString(16)

	if err != nil {
		return "", errors.WithCode(code.ErrRunTimeCaller, err.Error())
	}

	for {
		serverTokenRecord, err := ad.GetServiceTokenByCode(ctx, token)
		if errors.Is(err, gorm.ErrRecordNotFound) && serverTokenRecord == nil {
			return fmt.Sprintf("app-%s", token), nil
		} else if err != nil {
			return "", err
		}
	}
}
