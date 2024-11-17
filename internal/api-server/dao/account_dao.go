package dao

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/repo"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
	"gorm.io/gorm"
)

type AccountDao struct {
	db *gorm.DB
}

var _ repo.AccountRepo = (*AccountDao)(nil)

func (ad *AccountDao) GetAccountByEmail(context context.Context, email string) (*model.Account, error) {
	var account model.Account

	if err := ad.db.First(&account, "email = ?", email).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &account, nil
}
