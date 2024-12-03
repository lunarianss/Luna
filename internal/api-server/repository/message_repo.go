// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repo_impl

import (
	"context"
	"fmt"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/infrastructure/log"
	po_entity_account "github.com/lunarianss/Luna/internal/api-server/domain/account/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/repository"
	repo_common "github.com/lunarianss/Luna/internal/api-server/domain/common/repository"
	po_entity_web_app "github.com/lunarianss/Luna/internal/api-server/domain/web_app/entity/po_entity"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/lunarianss/Luna/internal/infrastructure/mysql"
	"gorm.io/gorm"
)

type MessageRepoImpl struct {
	db *gorm.DB
}

func NewMessageRepoImpl(db *gorm.DB) *MessageRepoImpl {
	return &MessageRepoImpl{db: db}
}

var _ repository.MessageRepo = (*MessageRepoImpl)(nil)

func (md *MessageRepoImpl) CreateMessage(ctx context.Context, message *po_entity.Message) (*po_entity.Message, error) {
	if err := md.db.Create(message).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return message, nil
}

func (md *MessageRepoImpl) DeletePinnedConversation(ctx context.Context, pinnedConversationID string) error {
	if err := md.db.Where("id = ?", pinnedConversationID).Delete(&po_entity.PinnedConversation{}).Error; err != nil {
		return err
	}
	return nil
}

func (md *MessageRepoImpl) CreatePinnedConversation(ctx context.Context, pinnedConversation *po_entity.PinnedConversation) (*po_entity.PinnedConversation, error) {
	if err := md.db.Create(pinnedConversation).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return pinnedConversation, nil
}

func (md *MessageRepoImpl) CreateConversation(ctx context.Context, conversation *po_entity.Conversation) (*po_entity.Conversation, error) {
	if err := md.db.Create(conversation).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return conversation, nil
}

func (md *MessageRepoImpl) UpdateMessage(ctx context.Context, message *po_entity.Message) error {
	if err := md.db.Updates(message).Error; err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil
}

func (md *MessageRepoImpl) UpdateConversationUpdateAt(ctx context.Context, appID string, conversation *po_entity.Conversation) error {
	if err := md.db.Model(conversation).Where("id = ? AND status = ? AND app_id = ?", conversation.ID, "normal", appID).Update("updated_at", conversation.UpdatedAt).Error; err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil
}
func (md *MessageRepoImpl) UpdateConversationName(ctx context.Context, conversation *po_entity.Conversation) error {
	if err := md.db.Model(conversation).Where("id = ?", conversation.ID).Select("name", "updated_at").Updates(conversation).Error; err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil
}

func (md *MessageRepoImpl) GetMessageByID(ctx context.Context, messageID string) (*po_entity.Message, error) {
	var message po_entity.Message

	if err := md.db.First(&message, "id = ?", messageID).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &message, nil
}

func (md *MessageRepoImpl) GetMessageByConversation(ctx context.Context, cID string, messageID string) (*po_entity.Message, error) {
	var message po_entity.Message

	if err := md.db.First(&message, "id = ? AND conversation_id = ?", messageID, cID).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &message, nil
}

func (md *MessageRepoImpl) GetConversationByID(ctx context.Context, conversationID string) (*po_entity.Conversation, error) {
	var conversation po_entity.Conversation

	if err := md.db.First(&conversation, "id = ?", conversationID).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &conversation, nil
}

func (md *MessageRepoImpl) GetMessageCountOfConversation(ctx context.Context, cID string) (int64, error) {
	var count int64
	if err := md.db.Model(&po_entity.Message{}).Where("conversation_id = ?", cID).Count(&count).Error; err != nil {
		return 0, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return count, nil
}

func (md *MessageRepoImpl) GetConversationByApp(ctx context.Context, conversationID string, appID string) (*po_entity.Conversation, error) {
	var conversation po_entity.Conversation

	if err := md.db.First(&conversation, "id = ? AND app_id = ?", conversationID, appID).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &conversation, nil
}

func (md *MessageRepoImpl) GetConversationByUser(ctx context.Context, appID, conversationID string, user repo_common.BaseAccount) (*po_entity.Conversation, error) {
	var conversation po_entity.Conversation

	var db *gorm.DB

	if _, ok := user.(*po_entity_web_app.EndUser); ok {
		db = md.db.Where("from_end_user_id = ?", user.GetAccountID())
	}

	if _, ok := user.(*po_entity_account.Account); ok {
		db = md.db.Where("from_account_id = ?", user.GetAccountID())
	}

	if err := db.Where("id = ? AND status = ? AND app_id = ?", conversationID, "normal", appID).First(&conversation).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &conversation, nil
}

func (md *MessageRepoImpl) GetPinnedConversationByConversation(ctx context.Context, appID, cID string, user repo_common.BaseAccount) (*po_entity.PinnedConversation, error) {
	var (
		conversation po_entity.PinnedConversation
	)
	if err := md.db.First(&conversation, "app_id = ? AND conversation_Id = ? AND created_by_role = ? AND created_by = ?", appID, cID, user.GetAccountType(), user.GetAccountID()).Error; err != nil {
		return nil, err
	}

	return &conversation, nil
}

func (md *MessageRepoImpl) FindPinnedConversationByUser(ctx context.Context, appID string, user repo_common.BaseAccount) ([]*po_entity.PinnedConversation, error) {
	var (
		conversations []*po_entity.PinnedConversation
	)
	if err := md.db.Model(&po_entity.PinnedConversation{}).Order("created_at DESC").Where("app_id = ? AND created_by_role = ? AND created_by = ?", appID, user.GetAccountType(), user.GetAccountID()).Find(&conversations).Error; err != nil {
		return nil, err
	}
	return conversations, nil
}

func (md *MessageRepoImpl) FindEndUserConversationsOrderByUpdated(ctx context.Context, appId string, invokeFrom string, user repo_common.BaseAccount, pageSize int, includeIDs []string, excludeIDs []string, lastID string, sortBy string) ([]*po_entity.Conversation, int64, error) {
	var (
		query         *gorm.DB
		fromSource    string
		fromEndUserID string
		fromAccountID string
		count         int64
		conversations []*po_entity.Conversation
	)

	if _, ok := user.(*po_entity_web_app.EndUser); ok {
		fromSource = "api"
		fromEndUserID = user.GetAccountID()
	}

	if _, ok := user.(*po_entity_account.Account); ok {
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

	if err := query.Model(&po_entity.Conversation{}).Count(&count).Limit(pageSize).Order(fmt.Sprintf("%s %s", sortField, sortDirection)).Find(&conversations).Error; err != nil {
		return nil, 0, err
	}

	return conversations, count, nil
}

func (md *MessageRepoImpl) FindConversationsInConsole(ctx context.Context, page, pageSize int, appID, start, end, sortBy, keyword string) ([]*po_entity.Conversation, int64, error) {
	var (
		conversations []*po_entity.Conversation
		count         int64
	)

	if err := md.db.Exec("SET SESSION sql_mode = REPLACE(@@sql_mode, 'ONLY_FULL_GROUP_BY', '')").Error; err != nil {
		return nil, 0, errors.WithCode(code.ErrDatabase, err.Error())
	}

	subQuery := md.db.Model(&po_entity.Conversation{}).Select("conversations.id AS conversation_id, end_users.session_id AS from_end_user_session_id").Joins("LEFT JOIN end_users ON conversations.from_end_user_id = end_users.id")

	mainQuery := md.db.Model(&po_entity.Conversation{}).Select("conversations.id, conversations.status, conversations.from_source, conversations.from_end_user_id, conversations.from_account_id, conversations.name, conversations.read_at, conversations.created_at, conversations.updated_at").Where("conversations.app_id = ?", appID)

	if keyword != "" {
		keywordFilter := fmt.Sprintf("%%%s%%", keyword)
		mainQuery = mainQuery.Joins("JOIN messages ON messages.conversation_id = conversations.id").Joins("JOIN (?) as subquery ON subquery.conversation_id = conversations.id", subQuery).Where("messages.query LIKE ? OR messages.answer LIKE ? OR conversations.name LIKE ? OR conversations.introduction LIKE ? OR subquery.from_end_user_session_id LIKE ?", keywordFilter, keywordFilter, keywordFilter, keywordFilter, keywordFilter)
		mainQuery = mainQuery.Group("conversations.id")
	}

	if sortBy != "" {
		switch sortBy {
		case "created_at":
			mainQuery = mainQuery.Order("conversations.created_at ASC")
		case "-created_at":
			mainQuery = mainQuery.Order("conversations.created_at DESC")
		case "updated_at":
			mainQuery = mainQuery.Order("conversations.updated_at ASC")
		case "-updated_at":
			mainQuery = mainQuery.Order("conversations.updated_at DESC")
		default:
			mainQuery = mainQuery.Order("conversations.created_at DESC")
		}
	}

	if err := mainQuery.Model(&po_entity.Conversation{}).Count(&count).Scopes(mysql.Paginate(page, pageSize)).Find(&conversations).Error; err != nil {
		return nil, 0, errors.WithCode(code.ErrDatabase, err.Error())
	}

	if err := md.db.Exec("SET SESSION sql_mode = CONCAT(@@sql_mode, ',ONLY_FULL_GROUP_BY')").Error; err != nil {
		log.Errorf(err.Error())
	}

	return conversations, count, nil
}

func (md *MessageRepoImpl) FindConsoleAppMessages(ctx context.Context, conversationID string, pageSize int, firstID string) ([]*po_entity.Message, int64, error) {
	var (
		ret   []*po_entity.Message
		count int64
	)

	if firstID != "" {
		messageRecord, err := md.GetMessageByConversation(ctx, conversationID, firstID)

		if err != nil {
			return nil, 0, err
		}
		if err := md.db.Model(&po_entity.Message{}).Order("created_at DESC").Limit(pageSize).Where("conversation_id = ? AND id != ? AND created_at < ?", conversationID, firstID, messageRecord.CreatedAt).Count(&count).Find(&ret).Error; err != nil {
			return nil, 0, err
		}
	} else {
		if err := md.db.Model(&po_entity.Message{}).Order("created_at DESC").Limit(pageSize).Where("conversation_id = ?", conversationID).Count(&count).Find(&ret).Error; err != nil {
			return nil, 0, errors.WithCode(code.ErrDatabase, err.Error())
		}
	}

	return ret, count, nil
}

func (md *MessageRepoImpl) FindEndUserMessages(ctx context.Context, appID string, user repo_common.BaseAccount, conversationId string, firstID string, pageSize int, order string) ([]*po_entity.Message, int64, error) {
	conversation, err := md.GetConversationByUser(ctx, appID, conversationId, user)

	var (
		count           int64
		historyMessages []*po_entity.Message
	)

	if err != nil {
		return nil, 0, err
	}

	if firstID != "" {
		messageRecord, err := md.GetMessageByConversation(ctx, conversationId, firstID)

		if err != nil {
			return nil, 0, err
		}
		if err := md.db.Model(&po_entity.Message{}).Order("created_at DESC").Limit(pageSize).Where("conversation_id = ? AND id != ? AND created_at < ?", conversation.ID, firstID, messageRecord.CreatedAt).Count(&count).Find(&historyMessages).Error; err != nil {
			return nil, 0, err
		}
	} else {
		if err := md.db.Model(&po_entity.Message{}).Order("created_at DESC").Limit(pageSize).Where("conversation_id = ?", conversation.ID).Count(&count).Find(&historyMessages).Error; err != nil {
			return nil, 0, err
		}
	}

	return historyMessages, count, nil

}

func (md *MessageRepoImpl) LogicalDeleteConversation(ctx context.Context, conversation *po_entity.Conversation) error {

	if err := md.db.Model(conversation).Where("id = ?", conversation.ID).Update("is_deleted", 1).Error; err != nil {
		return err
	}

	return nil
}
