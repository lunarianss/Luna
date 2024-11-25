package dao

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/repo"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
	"gorm.io/gorm"
)

type MessageDao struct {
	db *gorm.DB
}

func NewMessageDao(db *gorm.DB) *MessageDao {
	return &MessageDao{db: db}
}

var _ repo.MessageRepo = (*MessageDao)(nil)

func (md *MessageDao) CreateMessage(ctx context.Context, message *model.Message) (*model.Message, error) {
	if err := md.db.Create(message).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return message, nil
}

func (md *MessageDao) CreateConversation(ctx context.Context, conversation *model.Conversation) (*model.Conversation, error) {
	if err := md.db.Create(conversation).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return conversation, nil
}

func (md *MessageDao) UpdateMessage(ctx context.Context, message *model.Message) error {
	if err := md.db.Updates(message).Error; err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil
}

func (md *MessageDao) GetMessageByID(ctx context.Context, messageID string) (*model.Message, error) {
	var message model.Message

	if err := md.db.First(&message, "id = ?", messageID).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &message, nil
}

func (md *MessageDao) GetConversationByID(ctx context.Context, conversationID string) (*model.Conversation, error) {
	var conversation model.Conversation

	if err := md.db.First(&conversation, "id = ?", conversationID).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &conversation, nil
}
