// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repo

import (
	"context"

	model "github.com/Ryan-eng-del/hurricane/internal/apiServer/model/v1"
)

type BlogRepo interface {
	Create(ctx context.Context, blog *model.Blog) (*model.Blog, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, blog *model.Blog) (*model.Blog, error)
	Get(ctx context.Context, id int64) (*model.Blog, error)
	List(ctx context.Context, page, pageSize int) ([]*model.Blog, int64, error)
}
