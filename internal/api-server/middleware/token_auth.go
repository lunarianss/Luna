// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/lunarianss/Luna/internal/infrastructure/core"
	"github.com/lunarianss/Luna/internal/infrastructure/jwt"
	"github.com/lunarianss/Luna/infrastructure/errors"
)

func TokenAuthMiddleware() gin.HandlerFunc {
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

		lunaClaims, err := jwtIns.ParseLunaClaimsJWT(tokenString)

		if err != nil {
			core.WriteResponse(c, err, nil)
			c.Abort()
			return
		}

		if lunaClaims.AccountId == "" {
			core.WriteResponse(c, errors.WithCode(code.ErrTokenInvalid, "there is no account id after parse token"), nil)
			c.Abort()
			return
		}

		// 将当前用户信息保存到 Context
		c.Set("userID", lunaClaims.AccountId)

		// 继续处理请求
		c.Next()
	}
}
