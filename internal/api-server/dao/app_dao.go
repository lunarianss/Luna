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
