package service

import (
	"time"

	"github.com/gin-gonic/gin"
	jwtV5 "github.com/golang-jwt/jwt/v5"
	"github.com/lunarianss/Luna/internal/api-server/config"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/app_running"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/passport"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/internal/pkg/jwt"
	"github.com/lunarianss/Luna/pkg/errors"
)

type PassportService struct {
	appRunningDomain *domain.AppRunningDomain
	appDomain        *appDomain.AppDomain
	config           *config.Config
	jwt              *jwt.JWT
}

func NewPassportService(appRunningDomain *domain.AppRunningDomain,
	appDomain *appDomain.AppDomain,
	config *config.Config,
	jwt *jwt.JWT) *PassportService {
	return &PassportService{
		appRunningDomain: appRunningDomain,
		appDomain:        appDomain,
		config:           config,
		jwt:              jwt,
	}
}

func (ps *PassportService) AcquirePassport(c *gin.Context, appCode string) (*dto.AcquirePassportResponse, error) {
	siteRecord, err := ps.appRunningDomain.AppRunningRepo.GetSiteByCode(c, appCode)

	if err != nil {
		return nil, err
	}

	appRecord, err := ps.appDomain.AppRepo.GetAppByID(c, siteRecord.AppID)

	if err != nil {
		return nil, err
	}

	if appRecord.Status != "normal" {
		return nil, errors.WithCode(code.ErrAppStatusNotNormal, "status %s not normal", appRecord.Status)
	}

	endUserRecord, err := ps.appRunningDomain.CreateEndUser(c, appRecord)

	if err != nil {
		return nil, err
	}

	jwtClaims := jwt.LunaPassportClaims{
		RegisteredClaims: jwtV5.RegisteredClaims{
			ExpiresAt: jwtV5.NewNumericDate(time.Now().Add(ps.config.JwtOptions.Timeout)),
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
