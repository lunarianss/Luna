package dao

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/repo"
	"github.com/lunarianss/Luna/internal/pkg/code"
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
