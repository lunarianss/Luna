// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package domain

import "github.com/lunarianss/Hurricane/internal/api-server/repo"

type BlogDomain struct {
	BlogRepo repo.BlogRepo
}

func NewBlogDomain(blogRepo repo.BlogRepo) *BlogDomain {
	return &BlogDomain{
		BlogRepo: blogRepo,
	}
}
