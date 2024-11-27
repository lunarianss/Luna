package repository

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/_domain/account/entity/po_entity"
	"gorm.io/gorm"
)

type AccountRepo interface {
	CreateAccount(context context.Context, account *po_entity.Account, tx *gorm.DB) (*po_entity.Account, error)

	UpdateAccountIpAddress(context context.Context, account *po_entity.Account) error
	UpdateAccountStatus(context context.Context, account *po_entity.Account) error
	UpdateAccountLastActive(context context.Context, account *po_entity.Account) error

	GetAccountByEmail(context context.Context, email string) (*po_entity.Account, error)
	GetAccountByID(context context.Context, ID string) (*po_entity.Account, error)
}
