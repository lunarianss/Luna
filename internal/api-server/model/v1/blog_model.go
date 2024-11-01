// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"github.com/lunarianss/Luna/internal/pkg/field"
)

type Blog struct {
	Id          int64                  `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"             json:"id"`
	Title       string                 `gorm:"column:title;NOT NULL;comment:'文章标题'"                      json:"title"`
	BlogPic     string                 `gorm:"column:blog_pic;NOT NULL;comment:'文章首图，用于随机文章展示'"          json:"blog_pic"`
	Content     map[string]interface{} `gorm:"column:content;serializer:json;NOT NULL;comment:'文章正文'"    json:"content"`
	Description string                 `gorm:"column:description;NOT NULL;comment:'描述'"                  json:"description"`
	IsPublished field.BitBool          `gorm:"column:is_published;NOT NULL;comment:'公开或私密'"              json:"is_published"`
	Views       int32                  `gorm:"column:views;NOT NULL;comment:'浏览次数'"                      json:"views"`
	Words       int32                  `gorm:"column:words;NOT NULL;comment:'文章字数'"                      json:"words"`
	ReadTime    int32                  `gorm:"column:read_time;NOT NULL;comment:'阅读时长(分钟)"               json:"read_time"`
	CategoryId  int64                  `gorm:"column:category_id;NOT NULL;comment:'文章分类'"                json:"category"`
	IsTop       field.BitBool          `gorm:"column:is_top;NOT NULL;comment:'是否置顶'"                     json:"is_top"`
	IsDeleted   field.BitBool          `gorm:"column:is_deleted;NOT NULL;comment:'是否删除'"                 json:"-"`
	CreateTime  field.LocalTime        `gorm:"column:create_time;autoCreateTime;NOT NULL;comment:'创建时间'" json:"create_time"`
	UpdateTime  field.LocalTime        `gorm:"column:update_time;autoUpdateTime;NOT NULL;comment:'更新时间'" json:"update_time"`
}

func (b *Blog) TableName() string {
	return "blog"
}
