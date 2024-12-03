// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repo_impl

import (
	"context"

	"github.com/google/uuid"
	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/internal/api-server/domain/web_app/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/web_app/repository"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
	"gorm.io/gorm"
)

type WebAppRepoImpl struct {
	db *gorm.DB
}

var _ repository.WebAppRepo = (*WebAppRepoImpl)(nil)

func NewWebAppRepoImpl(db *gorm.DB) *WebAppRepoImpl {
	return &WebAppRepoImpl{
		db: db,
	}
}

func (ad *WebAppRepoImpl) GetSiteByAppID(ctx context.Context, appID string) (*po_entity.Site, error) {
	var site po_entity.Site
	if err := ad.db.First(&site, "app_id = ?", appID).Error; err != nil {
		return nil, err
	}
	return &site, nil
}

func (ad *WebAppRepoImpl) GetSiteByCode(ctx context.Context, code string) (*po_entity.Site, error) {
	var site po_entity.Site
	if err := ad.db.First(&site, "code = ? AND status = ?", code, "normal").Error; err != nil {
		return nil, err
	}
	return &site, nil
}

func (ad *WebAppRepoImpl) GetEndUserByID(ctx context.Context, endUserID string) (*po_entity.EndUser, error) {
	var endUser po_entity.EndUser
	if err := ad.db.First(&endUser, "id = ?", endUserID).Error; err != nil {
		return nil, err
	}
	return &endUser, nil
}

func (ad *WebAppRepoImpl) GetEndUserBySession(ctx context.Context, sessionID string) (*po_entity.EndUser, error) {
	var endUser po_entity.EndUser
	if err := ad.db.First(&endUser, "session_id = ?", sessionID).Error; err != nil {
		return nil, err
	}
	return &endUser, nil
}

func (ad *WebAppRepoImpl) GenerateSessionForEndUser(ctx context.Context) (string, error) {

	sessionID := uuid.NewString()

	for {
		endUserRecord, err := ad.GetEndUserBySession(ctx, sessionID)
		if errors.Is(err, gorm.ErrRecordNotFound) && endUserRecord == nil {
			return sessionID, nil
		} else if err != nil {
			return "", err
		}
	}
}

func (ad *WebAppRepoImpl) GenerateUniqueCodeForSite(ctx context.Context) (string, error) {

	siteCode, err := util.GenerateRandomString(16)

	if err != nil {
		return "", errors.WithCode(code.ErrRunTimeCaller, err.Error())
	}

	for {
		siteRecord, err := ad.GetSiteByCode(ctx, siteCode)
		if errors.Is(err, gorm.ErrRecordNotFound) && siteRecord == nil {
			return siteCode, nil
		} else if err != nil {
			return "", err
		}
	}
}

func (ad *WebAppRepoImpl) CreateSite(ctx context.Context, site *po_entity.Site, tx *gorm.DB) (*po_entity.Site, error) {

	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = ad.db
	}
	if err := dbIns.Create(site).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return site, nil

}
func (ad *WebAppRepoImpl) CreateInstallApp(ctx context.Context, installApp *po_entity.InstalledApp, tx *gorm.DB) (*po_entity.InstalledApp, error) {
	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = ad.db
	}
	if err := dbIns.Create(installApp).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return installApp, nil
}
func (ad *WebAppRepoImpl) CreateEndUser(ctx context.Context, endUser *po_entity.EndUser, tx *gorm.DB) (*po_entity.EndUser, error) {
	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = ad.db
	}
	if err := dbIns.Create(endUser).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return endUser, nil
}
