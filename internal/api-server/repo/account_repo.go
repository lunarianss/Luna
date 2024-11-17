package repo

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
)

type AccountRepo interface {
	GetAccountByEmail(context context.Context, email string) (*model.Account, error)
}
