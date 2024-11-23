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

func GetWebAppFromGin(g *gin.Context) (string, string, string, error) {
	appID, exist := g.Get("appID")
	appCode, appCodeExist := g.Get("appCode")
	endUser, endUserIDExist := g.Get("endUserID")

	if !exist || !appCodeExist || !endUserIDExist {
		return "", "", "", errors.WithCode(code.ErrGinNotExistAppSiteInfo, "")
	}

	appIDStr, ok := appID.(string)
	appCodeStr, appCodeOk := appCode.(string)
	endUserStr, endUserOk := endUser.(string)

	if !ok || !appCodeOk || !endUserOk {
		return "", "", "", errors.WithCode(code.ErrGinNotExistAppSiteInfo, "")
	}

	return appIDStr, appCodeStr, endUserStr, nil
}
