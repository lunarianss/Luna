package repo

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
)

type AppRepo interface {
	CreateApp(ctx context.Context, app *model.App) (*model.App, error)
	CreateAppWithConfig(ctx context.Context, app *model.App, appConfig *model.AppModelConfig) (*model.App, error)
}
