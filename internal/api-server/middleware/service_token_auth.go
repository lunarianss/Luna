// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/infrastructure/errors"
	po_account "github.com/lunarianss/Luna/internal/api-server/domain/account/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/lunarianss/Luna/internal/infrastructure/core"
	"github.com/lunarianss/Luna/internal/infrastructure/mysql"
	"gorm.io/gorm"
)

func ServiceTokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if !strings.HasPrefix(tokenString, "Bearer ") {
			core.WriteResponse(c, errors.WithCode(code.ErrTokenMissBearer, "token %s miss a header of Bearer ", tokenString), nil)
			c.Abort()
			return
		}

		// 截取 Bearer 前缀，获取 token
		tokenString = tokenString[7:]

		gormIns, err := mysql.GetMySQLIns(nil)

		if err != nil {
			core.WriteResponse(c, errors.WithSCode(code.ErrDatabase, err.Error()), nil)
			c.Abort()
			return
		}

		var apiToken po_entity.ApiToken

		if err := gormIns.First(&apiToken, "token = ? AND type = ?", tokenString, "app").Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				core.WriteResponse(c, errors.WithSCode(code.ErrResourceNotFound, err.Error()), nil)
				c.Abort()
				return
			} else {
				core.WriteResponse(c, errors.WithSCode(code.ErrDatabase, err.Error()), nil)
				c.Abort()
				return
			}
		}

		if err := gormIns.Model(&po_entity.ApiToken{}).Where("id = ?", apiToken.ID).Update("last_used_at", time.Now().UTC().Unix()).Error; err != nil {
			core.WriteResponse(c, errors.WithSCode(code.ErrDatabase, err.Error()), nil)
			c.Abort()
			return
		}

		appID := apiToken.AppID

		var app po_entity.App

		if err := gormIns.First(&app, "id = ?", appID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				core.WriteResponse(c, errors.WithSCode(code.ErrResourceNotFound, err.Error()), nil)
				c.Abort()
				return
			} else {
				core.WriteResponse(c, errors.WithSCode(code.ErrDatabase, err.Error()), nil)
				c.Abort()
				return
			}
		}

		if app.Status != "normal" {
			core.WriteResponse(c, errors.WithSCode(code.ErrAppStatusNotNormal, ""), nil)
			c.Abort()
			return
		}

		if app.EnableAPI == 0 {
			core.WriteResponse(c, errors.WithSCode(code.ErrAppApiDisabled, ""), nil)
			c.Abort()
			return
		}

		var tenant po_account.Tenant

		if err := gormIns.First(&tenant, "id = ?", app.TenantID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				core.WriteResponse(c, errors.WithSCode(code.ErrResourceNotFound, err.Error()), nil)
				c.Abort()
				return
			} else {
				core.WriteResponse(c, errors.WithSCode(code.ErrDatabase, err.Error()), nil)
				c.Abort()
				return
			}
		}

		if tenant.Status == "archive" {
			core.WriteResponse(c, errors.WithSCode(code.ErrTenantStatusArchive, ""), nil)
			c.Abort()
			return
		}

		c.Set("app", &app)
		c.Set("tenant", &tenant)
		c.Next()
	}
}
