package repo

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
)

type MessageRepo interface {
	CreateMessage(ctx context.Context, message *model.Message) (*model.Message, error)
	CreateConversation(ctx context.Context, app *model.Conversation) (*model.Conversation, error)

	UpdateMessage(ctx context.Context, message *model.Message) error

	GetMessageByID(ctx context.Context, messageID string) (*model.Message, error)
	GetConversationByID(ctx context.Context, conversationID string) (*model.Conversation, error)
}
