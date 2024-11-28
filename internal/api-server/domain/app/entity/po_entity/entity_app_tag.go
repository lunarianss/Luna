// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package po_entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tag struct
type Tag struct {
	ID        string `gorm:"column:id" json:"id"`
	TenantID  string `gorm:"column:tenant_id" json:"tenant_id"`
	Type      string `gorm:"column:type" json:"type"`
	Name      string `gorm:"column:name" json:"name"`
	CreatedBy string `gorm:"column:created_by" json:"created_by"`
	CreatedAt string `gorm:"column:created_at" json:"created_at"`
}

func (*Tag) TableName() string {
	return "tags"
}

func (u *Tag) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return
}

// TagBinding struct
type TagBinding struct {
	ID        string `gorm:"column:id" json:"id"`
	TenantID  string `gorm:"column:tenant_id" json:"tenant_id"`
	TagID     string `gorm:"column:tag_id" json:"tag_id"`
	TargetID  string `gorm:"column:target_id" json:"target_id"`
	CreatedBy string `gorm:"column:created_by" json:"created_by"`
	CreatedAt string `gorm:"column:created_at" json:"created_at"`
}

func (*TagBinding) TableName() string {
	return "tag_bindings"
}

func (u *TagBinding) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return
}
