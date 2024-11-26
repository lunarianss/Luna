package dao

import (
	"context"
	"fmt"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/repo"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/internal/pkg/mysql"
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

func (md *MessageDao) GetConversationByUser(ctx context.Context, appID, conversationID string, user model.BaseAccount) (*model.Conversation, error) {
	var conversation model.Conversation

	var db *gorm.DB

	if _, ok := user.(*model.EndUser); ok {
		db = md.db.Where("from_end_user_id = ?", user.GetAccountID())
	}

	if _, ok := user.(*model.Account); ok {
		db = md.db.Where("from_account_id = ?", user.GetAccountID())
	}

	if err := db.Where("id = ? AND status = ? AND app_id = ?", conversationID, "normal", appID).First(&conversation).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &conversation, nil
}

func (md *MessageDao) UpdateConversationUpdateAt(ctx context.Context, appID string, conversation *model.Conversation) error {
	if err := md.db.Model(conversation).Where("id = ? AND status = ? AND app_id = ?", conversation.ID, "normal", appID).Update("updated_at", conversation.UpdatedAt).Error; err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil
}

func (md *MessageDao) FindEndUserConversationsOrderByUpdated(ctx context.Context, appId string, invokeFrom string, user model.BaseAccount, pageSize int, includeIDs []string, excludeIDs []string, lastID string, sortBy string) ([]*model.Conversation, int64, error) {
	var (
		query         *gorm.DB
		fromSource    string
		fromEndUserID string
		fromAccountID string
		count         int64
		conversations []*model.Conversation
	)

	if _, ok := user.(*model.EndUser); ok {
		fromSource = "api"
		fromEndUserID = user.GetAccountID()
	}

	if _, ok := user.(*model.Account); ok {
		fromSource = "console"
		fromAccountID = user.GetAccountID()
	}

	query = md.db

	query = query.Scopes(mysql.LogicalObjects()).Where("app_id = ? AND from_source = ? AND from_end_user_id = ? AND from_account_id = ?", appId, fromSource, fromEndUserID, fromAccountID).Where("invoke_from = '' OR invoke_from = ?", invokeFrom)

	if includeIDs != nil {
		query = query.Where("id IN ?", includeIDs)
	}

	if excludeIDs != nil {
		query = query.Where("id NOT IN ?", excludeIDs)
	}

	sortField, sortDirection := mysql.GetSortParams(sortBy)
	sortOperation := mysql.BuildFilterCondition(sortField, sortDirection)

	if lastID != "" {
		lastConversation, err := md.GetConversationByID(ctx, lastID)
		if err != nil {
			return nil, 0, err
		}
		opStr := fmt.Sprintf("%s %s", sortField, sortOperation)
		query = query.Where(fmt.Sprintf("%s %d", opStr, lastConversation.UpdatedAt))
	}

	if err := query.Model(&model.Conversation{}).Count(&count).Find(&conversations).Order(fmt.Sprintf("%s %s", sortField, sortDirection)).Error; err != nil {
		return nil, 0, err
	}

	return conversations, count, nil
}
