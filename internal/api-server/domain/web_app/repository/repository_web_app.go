package repository

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/domain/web_app/entity/po_entity"
	"gorm.io/gorm"
)

type WebAppRepo interface {
	CreateSite(ctx context.Context, site *po_entity.Site, tx *gorm.DB) (*po_entity.Site, error)
	CreateInstallApp(ctx context.Context, site *po_entity.InstalledApp, tx *gorm.DB) (*po_entity.InstalledApp, error)
	CreateEndUser(ctx context.Context, endUser *po_entity.EndUser, tx *gorm.DB) (*po_entity.EndUser, error)
	GenerateUniqueCodeForSite(ctx context.Context) (string, error)
	GenerateSessionForEndUser(ctx context.Context) (string, error)

	GetSiteByAppID(ctx context.Context, appID string) (*po_entity.Site, error)
	GetSiteByCode(ctx context.Context, code string) (*po_entity.Site, error)
	GetEndUserBySession(ctx context.Context, sessionID string) (*po_entity.EndUser, error)
	GetEndUserByID(ctx context.Context, endUserID string) (*po_entity.EndUser, error)
}
