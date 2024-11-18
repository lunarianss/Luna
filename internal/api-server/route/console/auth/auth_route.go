package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/api-server/config"
	controller "github.com/lunarianss/Luna/internal/api-server/controller/gin/v1/auth"
	"github.com/lunarianss/Luna/internal/api-server/dao"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/account"
	tenantDomain "github.com/lunarianss/Luna/internal/api-server/domain/tenant"
	"github.com/lunarianss/Luna/internal/api-server/service"
	"github.com/lunarianss/Luna/internal/pkg/email"
	"github.com/lunarianss/Luna/internal/pkg/mysql"
	"github.com/lunarianss/Luna/internal/pkg/redis"
)

type AuthRoutes struct{}

func (a *AuthRoutes) Register(g *gin.Engine) error {
	gormIns, err := mysql.GetMySQLIns(nil)

	if err != nil {
		return err
	}

	redisIns, err := redis.GetRedisIns(nil)

	if err != nil {
		return err
	}

	email, err := email.GetEmailSMTPIns(nil)

	if err != nil {
		return err
	}

	// config
	config, err := config.GetLunaRuntimeConfig()

	if err != nil {
		return err
	}

	// dao
	accountDao := dao.NewAccountDao(gormIns)
	tenantDao := dao.NewTenantDao(gormIns)

	// domain
	accountDomain := domain.NewAccountDomain(accountDao, redisIns, config, email)
	tenantDomain := tenantDomain.NewTenantDomain(tenantDao)

	// service
	accountService := service.NewAccountService(accountDomain, tenantDomain)
	accountController := controller.NewAuthController(accountService)

	v1 := g.Group("/v1")
	authV1 := v1.Group("/console/api")
	authV1.POST("/email-code-login", accountController.SendEmailCode)
	authV1.POST("/email-code-login/validity", accountController.EmailValidity)
	return nil
}

func (r *AuthRoutes) GetModule() string {
	return "auth"
}
