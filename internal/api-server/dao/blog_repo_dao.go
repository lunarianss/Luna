// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dao

import (
	"context"

	"gorm.io/gorm"

	"github.com/lunarianss/Hurricane/internal/api-server/model/v1"
	"github.com/lunarianss/Hurricane/internal/api-server/repo"
	"github.com/lunarianss/Hurricane/internal/pkg/code"
	"github.com/lunarianss/Hurricane/internal/pkg/mysql"
	"github.com/lunarianss/Hurricane/pkg/errors"
)

type BlogDao struct {
	db *gorm.DB
}

var _ repo.BlogRepo = (*BlogDao)(nil)

func NewBlogDao(db *gorm.DB) *BlogDao {
	return &BlogDao{db}
}

func (u *BlogDao) Create(ctx context.Context, blog *model.Blog) (*model.Blog, error) {
	if err := u.db.Create(blog).Error; err != nil {
		return nil, err
	}
	return blog, nil
}

func (u *BlogDao) Delete(ctx context.Context, id int64) error {
	blog, err := u.Get(ctx, id)
	if err != nil {
		return err
	}
	blog.IsDeleted = 1
	if err := u.db.Save(blog).Error; err != nil {
		return err
	}
	return nil
}

func (u *BlogDao) Update(ctx context.Context, blog *model.Blog) (*model.Blog, error) {
	if err := u.db.Updates(blog).Error; err != nil {
		return nil, err
	}
	return blog, nil
}

func (u *BlogDao) Get(ctx context.Context, id int64) (*model.Blog, error) {
	var blog model.Blog

	if err := u.db.Scopes(mysql.LogicalObjects()).First(&blog, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrRecordNotFound, err.Error())
		}
		return nil, err
	}
	return &blog, nil
}

func (u *BlogDao) List(ctx context.Context, page, pageSize int) ([]*model.Blog, int64, error) {
	var (
		blogs []*model.Blog
		count int64
	)
	if err := u.db.Table("blog").Count(&count).Scopes(mysql.LogicalObjects(), mysql.Paginate(page, pageSize), mysql.IDDesc()).Find(&blogs).Error; err != nil {
		return nil, 0, err
	}
	return blogs, count, nil
}
