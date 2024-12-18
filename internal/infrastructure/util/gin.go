// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package util

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/infrastructure/errors"
	po_account "github.com/lunarianss/Luna/internal/api-server/domain/account/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
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

func GetServiceTokenFromGin(g *gin.Context) (*po_entity.App, *po_account.Tenant, error) {
	app, exist := g.Get("app")
	tenant, tenantExist := g.Get("tenant")

	if !exist || !tenantExist {
		return nil, nil, errors.WithCode(code.ErrGinNotExistServiceTokenInfo, "")
	}

	appRecord, ok := app.(*po_entity.App)
	tenantRecord, tenantOk := tenant.(*po_account.Tenant)

	if !ok || !tenantOk {
		return nil, nil, errors.WithCode(code.ErrGinNotExistServiceTokenInfo, "")
	}

	return appRecord, tenantRecord, nil
}
