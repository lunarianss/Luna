// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	"gorm.io/gorm"
)

type AppRepo interface {
	CreateApp(ctx context.Context, tx *gorm.DB, app *po_entity.App) (*po_entity.App, error)
	CreateAppWithConfig(ctx context.Context, tx *gorm.DB, app *po_entity.App, appConfig *po_entity.AppModelConfig) (*po_entity.App, error)
	CreateAppConfig(ctx context.Context, tx *gorm.DB, appConfig *po_entity.AppModelConfig) (*po_entity.AppModelConfig, error)
	UpdateAppConfigID(ctx context.Context, app *po_entity.App) error
	FindTenantApps(ctx context.Context, tenantID string, page, pageSize int) ([]*po_entity.App, int64, error)
	GetAppByID(ctx context.Context, appID string) (*po_entity.App, error)
	GetAppModelConfigById(ctx context.Context, appConfigID string, appID string) (*po_entity.AppModelConfig, error)
	GetTenantApp(ctx context.Context, appID, tenantID string) (*po_entity.App, error)
	UpdateEnableAppSite(ctx context.Context, app *po_entity.App) (*po_entity.App, error)
}
