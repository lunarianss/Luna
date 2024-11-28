package repository

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/_domain/app/entity/po_entity"
	"gorm.io/gorm"
)

type AppRepo interface {
	CreateApp(ctx context.Context, tx *gorm.DB, app *po_entity.App) (*po_entity.App, error)
	CreateAppWithConfig(ctx context.Context, tx *gorm.DB, app *po_entity.App, appConfig *po_entity.AppModelConfig) (*po_entity.App, error)
	CreateAppConfig(ctx context.Context, tx *gorm.DB, appConfig *po_entity.AppModelConfig) (*po_entity.AppModelConfig, error)

	UpdateAppConfigID(ctx context.Context, app *po_entity.App) error
	FindTenantApps(ctx context.Context, tenantID string, page, pageSize int) ([]*po_entity.App, int64, error)
	GetAppByID(ctx context.Context, appID string) (*po_entity.App, error)
	GetAppModelConfigById(ctx context.Context, appConfigID string) (*po_entity.AppModelConfig, error)
	GetAppModelConfigByAppID(ctx context.Context, appID string) (*po_entity.AppModelConfig, error)
}
