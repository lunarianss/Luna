package dao

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/repo"
	"gorm.io/gorm"
)

type AppDao struct {
	db *gorm.DB
}

var _ repo.AppRepo = (*AppDao)(nil)

func NewAppDao(db *gorm.DB) *AppDao {
	return &AppDao{db}
}

func (d *AppDao) CreateApp(ctx context.Context, app *model.App) (*model.App, error) {
	return nil, nil
}
