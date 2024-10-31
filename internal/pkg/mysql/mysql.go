// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mysql

import (
	"fmt"
	"sync"

	"github.com/Ryan-eng-del/hurricane/internal/pkg/options"
	"github.com/Ryan-eng-del/hurricane/pkg/db"
	"github.com/Ryan-eng-del/hurricane/pkg/log"
	"gorm.io/gorm"
)

var (
	once    sync.Once
	GormIns *gorm.DB
)

// GetMySQLIns create mysql factory with the given config.
func GetMySQLIns(opts *options.MySQLOptions) (*gorm.DB, error) {
	if opts == nil && GormIns == nil {
		return nil, fmt.Errorf("failed to get mysql store factory")
	}
	var err error

	var dbIns *gorm.DB

	once.Do(func() {
		options := &db.Options{
			Host:                  opts.Host,
			Username:              opts.Username,
			Password:              opts.Password,
			Database:              opts.Database,
			MaxIdleConnections:    opts.MaxIdleConnections,
			MaxOpenConnections:    opts.MaxOpenConnections,
			MaxConnectionLifeTime: opts.MaxConnectionLifeTime,
			LogLevel:              opts.LogLevel,
			Logger:                db.NewGormLogger(opts.LogLevel),
		}
		if dbIns, err = db.New(options); err != nil {
			log.Error("new gorm db instance error: %v", err)
			return
		}

		GormIns = dbIns
	})

	if GormIns == nil || err != nil {
		return nil, fmt.Errorf("failed to get mysql store factory, mysqlFactory: %+v, error: %w", dbIns, err)
	}

	return GormIns, nil
}
