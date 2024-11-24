package model

import (
	"github.com/google/uuid"
	"github.com/lunarianss/Luna/internal/pkg/field"
	"gorm.io/gorm"
)

type AccountStatus string

const (
	PENDING       AccountStatus = "pending"
	UNINITIALIZED AccountStatus = "uninitialized"
	ACTIVE        AccountStatus = "active"
	BANNED        AccountStatus = "banned"
	CLOSED        AccountStatus = "closed"
)

type BaseAccount interface {
	GetAccountType() string
	GetAccountID() string
}

func (*Account) GetAccountType() string {
	return "account"
}

func (a *Account) GetAccountID() string {
	return a.ID
}

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

func (a *Account) TableName() string {
	return "accounts"
}

func (a *Account) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.NewString()
	return
}

type TenantAccountJoinRole string

const (
	OWNER            TenantAccountJoinRole = "owner"
	ADMIN            TenantAccountJoinRole = "admin"
	NORMAL           TenantAccountJoinRole = "normal"
	DATASET_OPERATOR TenantAccountJoinRole = "dataset_operator"
)

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

func (a *TenantAccountJoin) TableName() string {
	return "tenant_account_joins"
}

func (a *TenantAccountJoin) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.NewString()

	return
}

type TenantStatus string

const (
	TNORMAL TenantStatus = "normal"
	ARCHIVE TenantStatus = "archive"
)

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

type TenantJoinResult struct {
	ID           string                 `json:"tenant_id" gorm:"column:tenant_id"`
	Name         string                 `json:"tenant_name" gorm:"column:tenant_name"`
	Plan         string                 `json:"tenant_plan" gorm:"column:tenant_plan"`
	Status       string                 `json:"tenant_status" gorm:"column:tenant_status"`
	CreatedAt    int64                  `json:"tenant_created_at" gorm:"column:tenant_created_at"`
	UpdatedAt    int64                  `json:"tenant_updated_at" gorm:"column:tenant_updated_at"`
	CustomConfig map[string]interface{} `json:"tenant_custom_config" gorm:"column:tenant_custom_config;serializer:json"`
	Role         string                 `json:"tenant_join_role" gorm:"column:tenant_join_role"`
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
		a.Status = string(TNORMAL)
	}

	return
}
