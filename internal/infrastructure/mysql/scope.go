// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mysql

import (
	"strings"

	"gorm.io/gorm"
)

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func IDDesc() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order("id DESC")
	}
}

func LogicalObjects() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("is_deleted = ?", 0)
	}
}

func GetSortParams(sortBy string) (string, string) {
	sortDirection := "ASC"
	sortField := strings.TrimPrefix(sortBy, "-")

	if strings.HasPrefix(sortBy, "-") {
		sortDirection = "DESC"
	}
	return sortField, sortDirection
}

func BuildFilterCondition(field, direction string) string {
	operator := ">"
	if direction == "DESC" {
		operator = "<"
	}
	return operator
}
