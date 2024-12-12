// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package po_entity

import (
	"slices"

	"github.com/google/uuid"

	"github.com/lunarianss/Luna/internal/infrastructure/field"
	"gorm.io/gorm"
)

type Account struct {
	ID                string `json:"id" gorm:"column:id"`
	Name              string `json:"name" gorm:"column:name"`
	Email             string `json:"email" gorm:"column:email"`
	Password          string `json:"password" gorm:"column:password"`
	PasswordSalt      string `json:"password_salt" gorm:"column:password_salt"`
	Avatar            string `json:"avatar" gorm:"column:avatar"`
	InterfaceLanguage string `json:"interface_language" gorm:"column:interface_language"`
	InterfaceTheme    string `json:"interface_theme" gorm:"column:interface_theme"`
	Timezone          string `json:"timezone" gorm:"column:timezone"`
	LastLoginIP       string `json:"last_login_ip" gorm:"column:last_login_ip"`
	Status            string `json:"status" gorm:"column:status"`
	LastLoginAt       *int64 `json:"last_login_at" gorm:"column:last_login_at"`
	InitializedAt     *int64 `json:"initialized_at" gorm:"column:initialized_at"`
	CreatedAt         int64  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt         int64  `json:"updated_at" gorm:"column:updated_at"`
	LastActiveAt      int64  `json:"last_active_at" gorm:"column:last_active_at;autoCreateTime"`
}

func (a *Account) GetAccountType() string {
	return "account"
}

func (a *Account) GetAccountID() string {
	return a.ID
}
func (a *Account) TableName() string {
	return "accounts"
}

func (a *Account) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.NewString()
	return
}

type Tenant struct {
	ID               string                 `json:"id" gorm:"column:id"`
	Name             string                 `json:"name" gorm:"column:name"`
	EncryptPublicKey string                 `json:"encrypt_public_key" gorm:"column:encrypt_public_key"`
	Plan             string                 `json:"plan" gorm:"column:plan"`
	Status           string                 `json:"status" gorm:"column:status"`
	CreatedAt        int64                  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt        int64                  `json:"updated_at" gorm:"column:updated_at"`
	CustomConfig     map[string]interface{} `json:"custom_config" gorm:"column:custom_config;serializer:json"`
}

func (a *Tenant) TableName() string {
	return "tenants"
}

func (a *Tenant) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.NewString()

	if a.Plan == "" {
		a.Plan = "basic"
	}

	if a.Status == "" {
		a.Status = string(TENANT_NORMAL)
	}

	return
}

type TenantAccountJoin struct {
	ID        string        `json:"id" gorm:"column:id"`
	TenantID  string        `json:"tenant_id" gorm:"column:tenant_id"`
	AccountID string        `json:"account_id" gorm:"column:account_id"`
	Role      string        `json:"role" gorm:"column:role"`
	InvitedBy string        `json:"invited_by" gorm:"column:invited_by"`
	CreatedAt int64         `json:"created_at" gorm:"column:created_at"`
	UpdatedAt int64         `json:"updated_at" gorm:"column:updated_at"`
	Current   field.BitBool `json:"current" gorm:"column:current"`
}

func (s *TenantAccountJoin) IsEditor() bool {
	editorRoles := []TenantAccountJoinRole{OWNER, ADMIN, EDITOR}
	return s.Role != "" && slices.Contains(editorRoles, TenantAccountJoinRole(s.Role))
}

func (s *TenantAccountJoin) IsPrivilegedRole() bool {
	privilegedRoles := []TenantAccountJoinRole{OWNER, ADMIN}
	return s.Role != "" && slices.Contains(privilegedRoles, TenantAccountJoinRole(s.Role))
}

func (s *TenantAccountJoin) IsNonOwnerRole() bool {
	editorRoles := []TenantAccountJoinRole{ADMIN, EDITOR, DATASET_OPERATOR, NORMAL}
	return s.Role != "" && slices.Contains(editorRoles, TenantAccountJoinRole(s.Role))
}

func (s *TenantAccountJoin) IsDatasetEditRole() bool {
	datasetEditorRoles := []TenantAccountJoinRole{OWNER, ADMIN, EDITOR, DATASET_OPERATOR}
	return s.Role != "" && slices.Contains(datasetEditorRoles, TenantAccountJoinRole(s.Role))
}

func (a *TenantAccountJoin) TableName() string {
	return "tenant_account_joins"
}

func (a *TenantAccountJoin) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.NewString()

	return
}
