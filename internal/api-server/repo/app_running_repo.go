package repo

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"gorm.io/gorm"
)

type AppRunningRepo interface {
	CreateSite(ctx context.Context, site *model.Site, tx *gorm.DB) (*model.Site, error)
	CreateInstallApp(ctx context.Context, site *model.InstalledApp, tx *gorm.DB) (*model.Site, error)
	CreateEndUser(ctx context.Context, site *model.EndUser, tx *gorm.DB) (*model.Site, error)
}
