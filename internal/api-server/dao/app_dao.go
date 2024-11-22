package dao

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/repo"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/internal/pkg/mysql"
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

func (ad *AppDao) CreateApp(ctx context.Context, tx *gorm.DB, app *model.App) (*model.App, error) {
	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = ad.db
	}

	if err := dbIns.Create(app).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return app, nil
}

func (ad *AppDao) FindTenantApps(ctx context.Context, tenant *model.Tenant, page, pageSize int) ([]*model.App, int64, error) {
	var apps []*model.App
	var appCount int64

	if err := ad.db.Table("apps").Count(&appCount).Scopes(mysql.Paginate(page, pageSize)).Find(&apps, "tenant_id = ? AND is_universal = ?", tenant.ID, 0).Error; err != nil {
		return nil, 0, err
	}
	return apps, appCount, nil
}

func (ad *AppDao) CreateConversation(ctx context.Context, conversation *model.Conversation) (*model.Conversation, error) {
	if err := ad.db.Create(conversation).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return conversation, nil
}
func (ad *AppDao) CreateMessage(ctx context.Context, message *model.Message) (*model.Message, error) {
	if err := ad.db.Create(message).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return message, nil
}

func (ad *AppDao) CreateAppWithConfig(ctx context.Context, tx *gorm.DB, app *model.App, appConfig *model.AppModelConfig) (*model.App, error) {
	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = ad.db
	}

	if err := dbIns.Create(app).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	appConfig.AppID = app.ID

	if err := dbIns.Create(appConfig).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	app.AppModelConfigID = appConfig.ID

	if err := dbIns.Model(app).Update("app_model_config_id", appConfig.ID).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return app, nil
}

func (ad *AppDao) GetAppByID(ctx context.Context, appID string) (*model.App, error) {
	var app model.App

	if err := ad.db.First(&app, "id = ?", appID).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &app, nil
}

func (ad *AppDao) GetMessageByID(ctx context.Context, messageID string) (*model.Message, error) {
	var message model.Message

	if err := ad.db.First(&message, "id = ?", messageID).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &message, nil
}
func (ad *AppDao) GetConversationByID(ctx context.Context, conversationID string) (*model.Conversation, error) {
	var conversation model.Conversation

	if err := ad.db.First(&conversation, "id = ?", conversationID).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &conversation, nil
}

func (ad *AppDao) GetAppModelConfigById(ctx context.Context, appConfigID string) (*model.AppModelConfig, error) {
	var appConfig model.AppModelConfig

	if err := ad.db.First(&appConfig, "id = ?", appConfigID).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &appConfig, nil
}

func (ad *AppDao) GetAppModelConfigByAppID(ctx context.Context, appID string) (*model.AppModelConfig, error) {
	var appConfig model.AppModelConfig

	if err := ad.db.First(&appConfig, "app_id = ?", appID).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &appConfig, nil
}
