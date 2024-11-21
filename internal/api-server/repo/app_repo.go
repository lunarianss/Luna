package repo

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
)

type AppRepo interface {
	CreateApp(ctx context.Context, app *model.App) (*model.App, error)
	CreateAppWithConfig(ctx context.Context, app *model.App, appConfig *model.AppModelConfig) (*model.App, error)
	CreateConversation(ctx context.Context, app *model.Conversation) (*model.Conversation, error)
	CreateMessage(ctx context.Context, message *model.Message) (*model.Message, error)

	FindTenantApps(ctx context.Context, tenant *model.Tenant, page, pageSize int) ([]*model.App, int64, error)
	GetAppByID(ctx context.Context, appID string) (*model.App, error)
	GetAppModelConfigById(ctx context.Context, appConfigID string) (*model.AppModelConfig, error)
	GetMessageByID(ctx context.Context, messageID string) (*model.Message, error)
	GetConversationByID(ctx context.Context, conversationID string) (*model.Conversation, error)
}
