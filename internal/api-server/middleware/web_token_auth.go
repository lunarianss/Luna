// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/lunarianss/Luna/internal/infrastructure/core"
	"github.com/lunarianss/Luna/internal/infrastructure/jwt"
	"github.com/lunarianss/Luna/internal/infrastructure/mysql"
	"gorm.io/gorm"
)

func WebTokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if !strings.HasPrefix(tokenString, "Bearer ") {
			core.WriteResponse(c, errors.WithCode(code.ErrTokenMissBearer, "token %s miss a header of Bearer ", tokenString), nil)
			c.Abort()
			return
		}

		// 截取 Bearer 前缀，获取 token
		tokenString = tokenString[7:]

		jwtIns, err := jwt.GetJWTIns()

		if err != nil {
			core.WriteResponse(c, err, nil)
			c.Abort()
			return
		}

		lunaClaims, err := jwtIns.ParseLunaPassportClaimsJWT(tokenString)

		if err != nil {
			core.WriteResponse(c, err, nil)
			c.Abort()
			return
		}

		if lunaClaims.AppCode == "" || lunaClaims.AppID == "" || lunaClaims.EndUserID == "" {
			core.WriteResponse(c, errors.WithCode(code.ErrTokenInvalid, "there is no web app info after parse web token"), nil)
			c.Abort()
			return
		}

		// 将当前用户信息保存到 Context
		c.Set("appID", lunaClaims.AppID)
		c.Set("appCode", lunaClaims.AppCode)
		c.Set("endUserID", lunaClaims.EndUserID)

		gormIns, err := mysql.GetMySQLIns(nil)

		if err != nil {
			core.WriteResponse(c, errors.WithCode(code.ErrDatabase, err.Error()), nil)
			c.Abort()
		}

		var app po_entity.App

		if err := gormIns.First(&app, "id = ?", lunaClaims.AppID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				core.WriteResponse(c, errors.WithCode(code.ErrResourceNotFound, err.Error()), nil)
				c.Abort()
			} else {
				core.WriteResponse(c, errors.WithCode(code.ErrDatabase, err.Error()), nil)
				c.Abort()
			}
		}

		if app.EnableSite == 0 {
			core.WriteResponse(c, errors.WithCode(code.ErrAppSiteDisabled, ""), nil)
			c.Abort()
		}
		// 继续处理请求
		c.Next()
	}
}
