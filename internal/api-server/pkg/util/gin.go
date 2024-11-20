package util

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
)

func GetUserIDFromGin(g *gin.Context) (string, error) {
	userID, exist := g.Get("userID")

	if !exist {
		return "", errors.WithCode(code.ErrGinNotExistAccountInfo, "")
	}

	userIDStr, ok := userID.(string)

	if !ok {
		return "", errors.WithCode(code.ErrGinNotExistAccountInfo, "userID is not a string")
	}

	return userIDStr, nil
}
