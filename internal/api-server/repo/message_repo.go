package repo

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
)

type MessageRepo interface {
	CreateMessage(ctx context.Context, message *model.Message) (*model.Message, error)
	CreateConversation(ctx context.Context, conversation *model.Conversation) (*model.Conversation, error)
	CreatePinnedConversation(ctx context.Context, pinnedConversation *model.PinnedConversation) (*model.PinnedConversation, error)

	DeletePinnedConversation(ctx context.Context, pinnedConversationID string) error

	UpdateMessage(ctx context.Context, message *model.Message) error
	UpdateConversationUpdateAt(ctx context.Context, appID string, conversation *model.Conversation) error

	GetMessageByID(ctx context.Context, messageID string) (*model.Message, error)
	GetConversationByID(ctx context.Context, conversationID string) (*model.Conversation, error)
	GetConversationByUser(ctx context.Context, appId string, conversationID string, user model.BaseAccount) (*model.Conversation, error)
	GetPinnedConversationByConversation(ctx context.Context, appID, cID string, user model.BaseAccount) (*model.PinnedConversation, error)
	// FindEndUserMessages(ctx context.Context, tenant *model.Tenant, page, pageSize int) ([]*model.Message, int64, error)
	FindEndUserConversationsOrderByUpdated(ctx context.Context, appId string, invokeFrom string, user model.BaseAccount, pageSize int, includeIDs []string, excludeIDs []string, lastID string, sortBy string) ([]*model.Conversation, int64, error)
	FindEndUserMessages(ctx context.Context, appID string, user model.BaseAccount, conversationId string, firstID string, pageSize int, order string) ([]*model.Message, int64, error)
	FindPinnedConversationByUser(ctx context.Context, appID string, user model.BaseAccount) ([]*model.PinnedConversation, error)
}
