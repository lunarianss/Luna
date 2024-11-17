package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	ID string `json:"id"`
}

func (a *Account) TableName() string {
	return "accounts"
}

func (a *Account) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.NewString()
	return
}
