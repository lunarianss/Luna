// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"time"

	"github.com/lunarianss/Luna/internal/pkg/field"
)

type ProviderType string

type ProviderQuotaType string

const (
	CUSTOM ProviderType = "custom"
	SYSTEM ProviderType = "system"
)

const (
	PAID ProviderQuotaType = "paid"

	FREE ProviderQuotaType = "free"

	TRIAL ProviderQuotaType = "trial"
)

type Provider struct {
	ID              string        `gorm:"column:id"                        json:"id"`
	TenantID        string        `gorm:"column:tenant_id"                 json:"tenant_id"`
	ProviderName    string        `gorm:"column:provider_name"             json:"provider_name"`
	ProviderType    string        `gorm:"column:provider_type"             json:"provider_type"`
	EncryptedConfig string        `gorm:"column:encrypted_config"          json:"encrypted_config,omitempty"`
	IsValid         field.BitBool `gorm:"column:is_valid"                  json:"is_valid"`
	LastUsed        *time.Time    `gorm:"column:last_used"                 json:"last_used,omitempty"`
	QuotaType       string        `gorm:"column:quota_type"                json:"quota_type,omitempty"`
	QuotaLimit      *int64        `gorm:"column:quota_limit"               json:"quota_limit,omitempty"`
	QuotaUsed       int64         `gorm:"column:quota_used"                json:"quota_used"`
	CreatedAt       int64         `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       int64         `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (*Provider) TableName() string {
	return "providers"
}
