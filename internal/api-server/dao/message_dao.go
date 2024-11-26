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

func (md *MessageDao) DeletePinnedConversation(ctx context.Context, pinnedConversationID string) error {
	if err := md.db.Where("id = ?", pinnedConversationID).Delete(&model.PinnedConversation{}).Error; err != nil {
		return err
	}
	return nil
}

func (md *MessageDao) CreatePinnedConversation(ctx context.Context, pinnedConversation *model.PinnedConversation) (*model.PinnedConversation, error) {
	if err := md.db.Create(pinnedConversation).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return pinnedConversation, nil
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

func (md *MessageDao) UpdateConversationUpdateAt(ctx context.Context, appID string, conversation *model.Conversation) error {
	if err := md.db.Model(conversation).Where("id = ? AND status = ? AND app_id = ?", conversation.ID, "normal", appID).Update("updated_at", conversation.UpdatedAt).Error; err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil
}
func (md *MessageDao) UpdateConversationName(ctx context.Context, conversation *model.Conversation) error {
	if err := md.db.Model(conversation).Where("id = ?", conversation.ID).Select("name", "updated_at").Updates(conversation).Error; err != nil {
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

func (md *MessageDao) GetPinnedConversationByConversation(ctx context.Context, appID, cID string, user model.BaseAccount) (*model.PinnedConversation, error) {
	var (
		conversation model.PinnedConversation
	)
	if err := md.db.First(&conversation, "app_id = ? AND conversation_Id = ? AND created_by_role = ? AND created_by = ?", appID, cID, user.GetAccountType(), user.GetAccountID()).Error; err != nil {
		return nil, err
	}

	return &conversation, nil
}

func (md *MessageDao) FindPinnedConversationByUser(ctx context.Context, appID string, user model.BaseAccount) ([]*model.PinnedConversation, error) {
	var (
		conversations []*model.PinnedConversation
	)
	if err := md.db.Model(&model.PinnedConversation{}).Order("created_at DESC").Where("app_id = ? AND created_by_role = ? AND created_by = ?", appID, user.GetAccountType(), user.GetAccountID()).Find(&conversations).Error; err != nil {
		return nil, err
	}
	return conversations, nil
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

	if err := query.Model(&model.Conversation{}).Count(&count).Limit(pageSize).Order(fmt.Sprintf("%s %s", sortField, sortDirection)).Find(&conversations).Error; err != nil {
		return nil, 0, err
	}

	return conversations, count, nil
}

func (md *MessageDao) FindEndUserMessages(ctx context.Context, appID string, user model.BaseAccount, conversationId string, firstID string, pageSize int, order string) ([]*model.Message, int64, error) {
	conversation, err := md.GetConversationByUser(ctx, appID, conversationId, user)

	var (
		count           int64
		historyMessages []*model.Message
	)

	if err != nil {
		return nil, 0, err
	}

	if err := md.db.Model(&model.Message{}).Count(&count).Order("created_at DESC").Limit(pageSize).Where("conversation_id = ?", conversation.ID).Find(&historyMessages).Error; err != nil {
		return nil, 0, err
	}
	return historyMessages, count, nil

}

func (md *MessageDao) LogicalDeleteConversation(ctx context.Context, conversation *model.Conversation) error {

	if err := md.db.Model(conversation).Where("id = ?", conversation.ID).Update("is_deleted", 1).Error; err != nil {
		return err
	}

	return nil
}
