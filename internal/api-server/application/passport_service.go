// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"time"

	"github.com/gin-gonic/gin"
	jwtV5 "github.com/golang-jwt/jwt/v5"
	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/internal/api-server/config"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	webAppDomain "github.com/lunarianss/Luna/internal/api-server/domain/web_app/domain_service"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/passport"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/lunarianss/Luna/internal/infrastructure/jwt"
)

type PassportService struct {
	webAppDomain *webAppDomain.WebAppDomain
	appDomain    *appDomain.AppDomain
	config       *config.Config
	jwt          *jwt.JWT
}

func NewPassportService(webAppDomain *webAppDomain.WebAppDomain,
	appDomain *appDomain.AppDomain,
	config *config.Config,
	jwt *jwt.JWT) *PassportService {
	return &PassportService{
		webAppDomain: webAppDomain,
		appDomain:    appDomain,
		config:       config,
		jwt:          jwt,
	}
}

func (ps *PassportService) AcquirePassport(c *gin.Context, appCode string) (*dto.AcquirePassportResponse, error) {
	siteRecord, err := ps.webAppDomain.WebAppRepo.GetSiteByCode(c, appCode)

	if err != nil {
		return nil, err
	}

	appRecord, err := ps.appDomain.AppRepo.GetAppByID(c, siteRecord.AppID)

	if err != nil {
		return nil, err
	}

	if appRecord.Status != "normal" || appRecord.EnableSite == 0 {
		return nil, errors.WithCode(code.ErrResourceNotFound, "status %s not normal or site status %v is disabled", appRecord.Status, appRecord.EnableSite)
	}

	endUserRecord, err := ps.webAppDomain.CreateEndUser(c, appRecord)

	if err != nil {
		return nil, err
	}

	jwtClaims := jwt.LunaPassportClaims{
		RegisteredClaims: jwtV5.RegisteredClaims{
			IssuedAt:  jwtV5.NewNumericDate(time.Now()),
			NotBefore: jwtV5.NewNumericDate(time.Now()),
			Issuer:    ps.config.JwtOptions.Realm,
			Subject:   "Web API Token",
			Audience:  []string{"not yet"},
		},
		AppID:     siteRecord.AppID,
		AppCode:   appCode,
		EndUserID: endUserRecord.ID,
	}

	token, err := ps.jwt.GenerateJWT(jwtClaims)

	if err != nil {
		return nil, err
	}

	return &dto.AcquirePassportResponse{AccessToken: token}, nil
}
