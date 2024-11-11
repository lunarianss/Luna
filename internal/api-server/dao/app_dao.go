package dao

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/repo"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
	"gorm.io/gorm"
)

type AppDao struct {
	db *gorm.DB
}

var _ repo.AppRepo = (*AppDao)(nil)

func NewAppDao(db *gorm.DB) *AppDao {
	return &AppDao{db}
}

func (ad *AppDao) CreateApp(ctx context.Context, app *model.App) (*model.App, error) {
	if err := ad.db.Create(app).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return app, nil
}

func (ad *AppDao) CreateAppWithConfig(ctx context.Context, app *model.App, appConfig *model.AppModelConfig) (*model.App, error) {
	tx := ad.db.Begin()
	if err := tx.Create(app).Error; err != nil {
		tx.Rollback()
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	appConfig.AppID = app.ID

	if err := tx.Create(appConfig).Error; err != nil {
		tx.Rollback()
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return app, tx.Commit().Error
}

func (ad *AppDao) GetAppByID(ctx context.Context, appID string) (*model.App, error) {
	var app model.App

	if err := ad.db.First(&app, "id = ?", appID).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &app, nil
}
