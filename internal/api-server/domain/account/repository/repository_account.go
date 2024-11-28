// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/domain/account/entity/po_entity"
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
