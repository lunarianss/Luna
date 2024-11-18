package repo

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
)

type AccountRepo interface {
	GetAccountByEmail(context context.Context, email string) (*model.Account, error)
	CreateAccount(context context.Context, account *model.Account) (*model.Account, error)
	UpdateAccountIpAddress(context context.Context, account *model.Account) error
	UpdateAccountStatus(context context.Context, account *model.Account) error
}
