// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repo_impl

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/domain/account/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/account/repository"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/lunarianss/Luna/infrastructure/errors"
	"gorm.io/gorm"
)

type AccountRepoImpl struct {
	db *gorm.DB
}

func NewAccountRepoImpl(db *gorm.DB) *AccountRepoImpl {
	return &AccountRepoImpl{db: db}
}

var _ repository.AccountRepo = (*AccountRepoImpl)(nil)

func (ad *AccountRepoImpl) GetAccountByEmail(context context.Context, email string) (*po_entity.Account, error) {
	var account po_entity.Account

	if err := ad.db.Limit(1).Find(&account, "email = ?", email).Error; err != nil {
		return nil, errors.WithSCode(code.ErrDatabase, err.Error())
	}
	return &account, nil
}
func (ad *AccountRepoImpl) GetAccountByID(context context.Context, ID string) (*po_entity.Account, error) {
	var account po_entity.Account

	if err := ad.db.First(&account, "id = ?", ID).Error; err != nil {
		return nil, errors.WithSCode(code.ErrDatabase, err.Error())
	}
	return &account, nil
}

func (ad *AccountRepoImpl) CreateAccount(context context.Context, account *po_entity.Account, tx *gorm.DB) (*po_entity.Account, error) {

	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = ad.db
	}

	if err := dbIns.Create(account).Error; err != nil {
		return nil, errors.WithSCode(code.ErrDatabase, err.Error())
	}
	return account, nil
}

func (ad *AccountRepoImpl) UpdateAccountIpAddress(context context.Context, account *po_entity.Account) error {
	if err := ad.db.Model(account).Where("id = ?", account.ID).Select("last_login_at", "last_login_ip").Updates(account).Error; err != nil {
		return errors.WithSCode(code.ErrDatabase, err.Error())
	}
	return nil
}

func (ad *AccountRepoImpl) UpdateAccountStatus(context context.Context, account *po_entity.Account) error {
	if err := ad.db.Model(account).Where("id = ?", account.ID).Update("status", account.Status).Error; err != nil {
		return errors.WithSCode(code.ErrDatabase, err.Error())
	}
	return nil
}

func (ad *AccountRepoImpl) UpdateAccountLastActive(context context.Context, account *po_entity.Account) error {
	if err := ad.db.Model(account).Where("id = ?", account.ID).Update("last_active_at", account.LastActiveAt).Error; err != nil {
		return errors.WithSCode(code.ErrDatabase, err.Error())
	}
	return nil
}
