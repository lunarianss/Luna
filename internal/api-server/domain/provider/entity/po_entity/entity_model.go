package po_entity

import (
	"github.com/google/uuid"
	"github.com/lunarianss/Luna/internal/pkg/field"
	"gorm.io/gorm"
)

type ProviderModel struct {
	ID              string        `gorm:"column:id"                        json:"id"`
	TenantID        string        `gorm:"column:tenant_id"                 json:"tenant_id"`
	ProviderName    string        `gorm:"column:provider_name"             json:"provider_name"`
	ModelName       string        `gorm:"column:model_name"             json:"mode_name"`
	ModelType       string        `gorm:"column:model_type"            json:"model_type"`
	EncryptedConfig string        `gorm:"column:encrypted_config"          json:"encrypted_config,omitempty"`
	IsValid         field.BitBool `gorm:"column:is_valid"                  json:"is_valid"`
	CreatedAt       int64         `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       int64         `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (*ProviderModel) TableName() string {
	return "provider_models"
}

func (u *ProviderModel) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return
}

type TenantDefaultModel struct {
	ID           string `gorm:"column:id"                        json:"id"`
	TenantID     string `gorm:"column:tenant_id"                 json:"tenant_id"`
	ProviderName string `gorm:"column:provider_name"             json:"provider_name"`
	ModelName    string `gorm:"column:model_name"             json:"mode_name"`
	ModelType    string `gorm:"column:model_type"            json:"model_type"`
	CreatedAt    int64  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    int64  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (*TenantDefaultModel) TableName() string {
	return "tenant_default_models"
}

func (u *TenantDefaultModel) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return
}
