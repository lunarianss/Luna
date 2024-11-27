package route

import (
	"github.com/gin-gonic/gin"
	domain "github.com/lunarianss/Luna/internal/api-server/_domain/account/domain_service"
	repo_impl "github.com/lunarianss/Luna/internal/api-server/_repo"
	"github.com/lunarianss/Luna/internal/api-server/config"
	controller "github.com/lunarianss/Luna/internal/api-server/controller/gin/v1/auth"
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

	// repo
	accountRepo := repo_impl.NewAccountRepo(gormIns)
	tenantRepo := repo_impl.NewTenantRepo(gormIns)

	// domain
	accountDomain := domain.NewAccountDomain(accountRepo, redisIns, config, email, tenantRepo)
	tenantDomain := domain.NewTenantDomain(tenantRepo)

	// service
	accountService := service.NewAccountService(accountDomain, tenantDomain, gormIns)
	accountController := controller.NewAuthController(accountService)

	v1 := g.Group("/v1")
	authV1 := v1.Group("/console/api")
	authV1.POST("/email-code-login", accountController.SendEmailCode)
	authV1.POST("/email-code-login/validity", accountController.EmailValidity)
	authV1.POST("/refresh-token", accountController.RefreshToken)
	return nil
}

func (r *AuthRoutes) GetModule() string {
	return "auth"
}
