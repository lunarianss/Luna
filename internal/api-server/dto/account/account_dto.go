// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

import (
	"github.com/lunarianss/Luna/internal/api-server/domain/account/entity/po_entity"
)

type GetAccountProfileResp struct {
	ID                string `json:"id" gorm:"column:id"`
	Name              string `json:"name" gorm:"column:name"`
	Email             string `json:"email" gorm:"column:email"`
	Avatar            string `json:"avatar" gorm:"column:avatar"`
	InterfaceLanguage string `json:"interface_language" gorm:"column:interface_language"`
	InterfaceTheme    string `json:"interface_theme" gorm:"column:interface_theme"`
	Timezone          string `json:"timezone" gorm:"column:timezone"`
	LastLoginIP       string `json:"last_login_ip" gorm:"column:last_login_ip"`
	LastLoginAt       *int64 `json:"last_login_at" gorm:"column:last_login_at"`
	CreatedAt         int64  `json:"created_at" gorm:"column:created_at"`
	IsPasswordSet     bool   `json:"is_password_set"`
}

func AccountConvertToProfile(a *po_entity.Account) (s *GetAccountProfileResp) {
	return &GetAccountProfileResp{
		ID:                a.ID,
		Name:              a.Name,
		Email:             a.Email,
		Avatar:            a.Avatar,
		InterfaceLanguage: a.InterfaceLanguage,
		InterfaceTheme:    a.InterfaceTheme,
		Timezone:          a.Timezone,
		LastLoginIP:       a.LastLoginIP,
		LastLoginAt:       a.LastLoginAt,
		CreatedAt:         a.CreatedAt,
	}
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	Data   *TokenPair `json:"data"`
	Result string     `json:"result"`
}
