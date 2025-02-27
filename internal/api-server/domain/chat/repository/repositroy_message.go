// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"context"

	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/console_app_statistic"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/common/repository"
)

type MessageRepo interface {
	// Create
	CreateMessage(ctx context.Context, message *po_entity.Message) (*po_entity.Message, error)
	CreateConversation(ctx context.Context, conversation *po_entity.Conversation) (*po_entity.Conversation, error)
	CreatePinnedConversation(ctx context.Context, pinnedConversation *po_entity.PinnedConversation) (*po_entity.PinnedConversation, error)

	// Delete
	DeletePinnedConversation(ctx context.Context, pinnedConversationID string) error
	LogicalDeleteConversation(ctx context.Context, conversation *po_entity.Conversation) error

	// Update
	UpdateMessage(ctx context.Context, message *po_entity.Message) error
	UpdateMessageMetadata(ctx context.Context, message *po_entity.Message) error
	UpdateConversationUpdateAt(ctx context.Context, appID string, conversation *po_entity.Conversation) error
	UpdateConversationName(ctx context.Context, conversation *po_entity.Conversation) error

	// Get
	GetMessageByID(ctx context.Context, messageID string) (*po_entity.Message, error)
	GetMessageByApp(ctx context.Context, messageID string, appID string) (*po_entity.Message, error)
	GetConversationByID(ctx context.Context, conversationID string) (*po_entity.Conversation, error)
	GetConversationByApp(ctx context.Context, conversationID string, appID string) (*po_entity.Conversation, error)
	GetConversationByUser(ctx context.Context, appId string, conversationID string, user repository.BaseAccount) (*po_entity.Conversation, error)
	GetPinnedConversationByConversation(ctx context.Context, appID, cID string, user repository.BaseAccount) (*po_entity.PinnedConversation, error)
	GetMessageCountOfConversation(ctx context.Context, cID string) (int64, error)
	GetMessageByConversation(ctx context.Context, cID string, messageID string) (*po_entity.Message, error)
	StatisticDailyConversations(ctx context.Context, appID, start, end, location string) ([]*biz_entity.StatisticDailyConversationsItem, error)
	StatisticTokenCosts(ctx context.Context, appID, start, end, location string) ([]*biz_entity.StatisticTokenCostsItem, error)
	StatisticDailyMessages(ctx context.Context, appID, start, end, location string) ([]*biz_entity.StatisticDailyConversationsItem, error)
	StatisticDailyUsers(ctx context.Context, appID, start, end, location string) ([]*biz_entity.StatisticDailyUsersItem, error)
	StatisticAverageSessionInteraction(ctx context.Context, appID, start, end, location string) ([]*biz_entity.StatisticAverageInteractionItem, error)
	// Find
	FindConversationsInConsole(ctx context.Context, page, pageSize int, appID, start, end, sortBy, keyword, timezone string) ([]*po_entity.Conversation, int64, error)
	FindEndUserConversationsOrderByUpdated(ctx context.Context, appId string, invokeFrom string, user repository.BaseAccount, pageSize int, includeIDs []string, excludeIDs []string, lastID string, sortBy string) ([]*po_entity.Conversation, int64, error)
	FindEndUserMessages(ctx context.Context, appID string, user repository.BaseAccount, conversationId string, firstID string, pageSize int, order string) ([]*po_entity.Message, int64, error)
	FindConsoleAppMessages(ctx context.Context, conversationID string, pageSize int, firstID string) ([]*po_entity.Message, int64, error)
	FindPinnedConversationByUser(ctx context.Context, appID string, user repository.BaseAccount) ([]*po_entity.PinnedConversation, error)
	FindHistoryPromptMessage(ctx context.Context, conversationID string, limit int) ([]*po_entity.Message, error)
}
