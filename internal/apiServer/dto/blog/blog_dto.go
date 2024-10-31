// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

import (
	"github.com/Ryan-eng-del/hurricane/internal/apiServer/model/v1"
	"github.com/Ryan-eng-del/hurricane/internal/pkg/field"
)

// create
type CreateBlogRequest struct {
	Title       string                 `json:"title" validate:"required"`
	BlogPic     string                 `json:"blog_pic" validate:"required"`
	Content     map[string]interface{} `json:"content" validate:"required"`
	Words       int32                  `json:"words" validate:"required"`
	ReadTime    int32                  `json:"read_time" validate:"required"`
	CategoryId  int64                  `json:"category_id" validate:"required"`
	Description string                 `json:"description"`
	IsPublished int                    `json:"is_published" validate:"min=0,max=1"`
	IsTop       int                    `json:"is_top" validate:"min=0,max=1"`
}

// get list
type GetBlogRequest struct {
	Page     int `json:"page" form:"page" validate:"required,min=1"`
	PageSize int `json:"pageSize" form:"page_size" validate:"required,min=1,max=100"`
}

type Blog struct {
	Id          int64                  `json:"id"`
	Title       string                 `json:"title"`
	BlogPic     string                 `json:"blog_pic"`
	Content     map[string]interface{} `json:"content"`
	Description string                 `json:"description"`
	IsPublished field.BitBool          `json:"is_published"`
	CreateTime  field.LocalTime        `json:"create_time"`
	UpdateTime  field.LocalTime        `json:"update_time"`
	Views       int32                  `json:"views"`
	Words       int32                  `json:"words"`
	ReadTime    int32                  `json:"read_time"`
	IsTop       field.BitBool          `json:"is_top"`
	IsDeleted   field.BitBool          `json:"is_deleted"`
}

type GetBlogListResponse struct {
	Count    int64   `json:"count"`
	Items    []*Blog `json:"items"`
	NextPage int     `json:"next_page"`
}

// update
type UpdateBlogRequest struct {
	Title       string                 `json:"title"`
	BlogPic     string                 `json:"blog_pic"`
	Content     map[string]interface{} `json:"content"`
	Words       int32                  `json:"words"`
	ReadTime    int32                  `json:"read_time"`
	CategoryId  int64                  `json:"category_id"`
	Description string                 `json:"description"`
	IsPublished int                    `json:"is_published" validate:"min=0,max=1"`
	IsTop       int                    `json:"is_top" validate:"min=0,max=1"`
}

func Convert(b *model.Blog) *Blog {
	return &Blog{
		Id:          b.Id,
		Title:       b.Title,
		BlogPic:     b.BlogPic,
		Content:     b.Content,
		Words:       b.Words,
		ReadTime:    b.ReadTime,
		Description: b.Description,
		IsPublished: b.IsPublished,
		IsTop:       b.IsTop,
	}
}
