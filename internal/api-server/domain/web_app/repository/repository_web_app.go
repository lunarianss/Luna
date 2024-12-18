// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

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
	GetEndUserByInfo(ctx context.Context, sessionID string, tenantID string, appID string, endUserType string) (*po_entity.EndUser, error)
}
