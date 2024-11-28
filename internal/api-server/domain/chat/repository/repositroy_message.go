package repository

import (
	"context"

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
	UpdateConversationUpdateAt(ctx context.Context, appID string, conversation *po_entity.Conversation) error
	UpdateConversationName(ctx context.Context, conversation *po_entity.Conversation) error
	// Get
	GetMessageByID(ctx context.Context, messageID string) (*po_entity.Message, error)
	GetConversationByID(ctx context.Context, conversationID string) (*po_entity.Conversation, error)
	GetConversationByUser(ctx context.Context, appId string, conversationID string, user repository.BaseAccount) (*po_entity.Conversation, error)
	GetPinnedConversationByConversation(ctx context.Context, appID, cID string, user repository.BaseAccount) (*po_entity.PinnedConversation, error)

	// Find
	FindEndUserConversationsOrderByUpdated(ctx context.Context, appId string, invokeFrom string, user repository.BaseAccount, pageSize int, includeIDs []string, excludeIDs []string, lastID string, sortBy string) ([]*po_entity.Conversation, int64, error)
	FindEndUserMessages(ctx context.Context, appID string, user repository.BaseAccount, conversationId string, firstID string, pageSize int, order string) ([]*po_entity.Message, int64, error)
	FindPinnedConversationByUser(ctx context.Context, appID string, user repository.BaseAccount) ([]*po_entity.PinnedConversation, error)
}
