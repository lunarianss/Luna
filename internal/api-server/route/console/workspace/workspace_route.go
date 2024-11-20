package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/api-server/config"
	controller "github.com/lunarianss/Luna/internal/api-server/controller/gin/v1/workspace"
	"github.com/lunarianss/Luna/internal/api-server/dao"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/account"
	tenantDomain "github.com/lunarianss/Luna/internal/api-server/domain/tenant"
	"github.com/lunarianss/Luna/internal/api-server/middlewares"
	"github.com/lunarianss/Luna/internal/api-server/service"
	"github.com/lunarianss/Luna/internal/pkg/email"
	"github.com/lunarianss/Luna/internal/pkg/mysql"
	"github.com/lunarianss/Luna/internal/pkg/redis"
)

type WorkspaceRoutes struct{}

func (a *WorkspaceRoutes) Register(g *gin.Engine) error {
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
	accountDomain := domain.NewAccountDomain(accountDao, redisIns, config, email, tenantDao)
	tenantDomain := tenantDomain.NewTenantDomain(tenantDao)

	// service
	tenantService := service.NewTenantService(accountDomain, tenantDomain)

	workspaceController := controller.NewWorkspaceController(tenantService)
	v1 := g.Group("/v1")
	authV1 := v1.Group("/console/api")
	workspaceV1 := authV1.Group("/workspaces")
	workspaceV1.Use(middlewares.TokenAuthMiddleware())

	workspaceV1.GET("/current", workspaceController.GetTenantCurrentWorkspace)

	return nil
}

func (r *WorkspaceRoutes) GetModule() string {
	return "setup"
}
