package repo

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"gorm.io/gorm"
)

type AccountRepo interface {
	CreateAccount(context context.Context, account *model.Account, tx *gorm.DB) (*model.Account, error)

	GetAccountByEmail(context context.Context, email string) (*model.Account, error)
	GetAccountByID(context context.Context, ID string) (*model.Account, error)

	UpdateAccountIpAddress(context context.Context, account *model.Account) error
	UpdateAccountStatus(context context.Context, account *model.Account) error
	UpdateAccountLastActive(context context.Context, account *model.Account) error
}
