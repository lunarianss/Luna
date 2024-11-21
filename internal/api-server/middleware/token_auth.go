package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/internal/pkg/core"
	"github.com/lunarianss/Luna/internal/pkg/jwt"
	"github.com/lunarianss/Luna/pkg/errors"
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

		jwtIns := jwt.GetJWTIns()

		lunaClaims, err := jwtIns.ParseJWT(tokenString)

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
