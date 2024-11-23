package dao

import (
	"context"

	"github.com/google/uuid"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/repo"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/internal/pkg/util"
	"github.com/lunarianss/Luna/pkg/errors"
	"gorm.io/gorm"
)

type AppRunningDao struct {
	db *gorm.DB
}

var _ repo.AppRunningRepo = (*AppRunningDao)(nil)

func NewAppRunningDao(db *gorm.DB) *AppRunningDao {
	return &AppRunningDao{
		db: db,
	}
}

func (ad *AppRunningDao) GetSiteByAppID(ctx context.Context, appID string) (*model.Site, error) {
	var site model.Site
	if err := ad.db.First(&site, "app_id = ?", appID).Error; err != nil {
		return nil, err
	}
	return &site, nil
}

func (ad *AppRunningDao) GetSiteByCode(ctx context.Context, code string) (*model.Site, error) {
	var site model.Site
	if err := ad.db.First(&site, "code = ? AND status = ?", code, "normal").Error; err != nil {
		return nil, err
	}
	return &site, nil
}

func (ad *AppRunningDao) GetEndUserByID(ctx context.Context, endUserID string) (*model.EndUser, error) {
	var endUser model.EndUser
	if err := ad.db.First(&endUser, "id = ?", endUserID).Error; err != nil {
		return nil, err
	}
	return &endUser, nil
}

func (ad *AppRunningDao) GetEndUserBySession(ctx context.Context, sessionID string) (*model.EndUser, error) {
	var endUser model.EndUser
	if err := ad.db.First(&endUser, "session_id = ?", sessionID).Error; err != nil {
		return nil, err
	}
	return &endUser, nil
}

func (ad *AppRunningDao) GenerateSessionForEndUser(ctx context.Context) (string, error) {

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

func (ad *AppRunningDao) GenerateUniqueCodeForSite(ctx context.Context) (string, error) {

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

func (ad *AppRunningDao) CreateSite(ctx context.Context, site *model.Site, tx *gorm.DB) (*model.Site, error) {

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
func (ad *AppRunningDao) CreateInstallApp(ctx context.Context, installApp *model.InstalledApp, tx *gorm.DB) (*model.InstalledApp, error) {
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
func (ad *AppRunningDao) CreateEndUser(ctx context.Context, endUser *model.EndUser, tx *gorm.DB) (*model.EndUser, error) {
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
