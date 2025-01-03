// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package po_entity

import "github.com/lunarianss/Luna/internal/infrastructure/field"

type TenantJoinResult struct {
	ID           string                 `json:"tenant_id" gorm:"column:tenant_id"`
	Name         string                 `json:"tenant_name" gorm:"column:tenant_name"`
	Plan         string                 `json:"tenant_plan" gorm:"column:tenant_plan"`
	Status       string                 `json:"tenant_status" gorm:"column:tenant_status"`
	CreatedAt    int64                  `json:"tenant_created_at" gorm:"column:tenant_created_at"`
	UpdatedAt    int64                  `json:"tenant_updated_at" gorm:"column:tenant_updated_at"`
	CustomConfig map[string]interface{} `json:"tenant_custom_config" gorm:"column:tenant_custom_config;serializer:json"`
	Role         string                 `json:"tenant_join_role" gorm:"column:tenant_join_role"`
	Current      field.BitBool          `json:"current" gorm:"column:tenant_join_current"`
}
