package route

import (
	"github.com/gin-gonic/gin"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/_domain/account/domain_service"
	repo_impl "github.com/lunarianss/Luna/internal/api-server/_repo"
	"github.com/lunarianss/Luna/internal/api-server/config"
	controller "github.com/lunarianss/Luna/internal/api-server/controller/gin/v1/workspace"
	"github.com/lunarianss/Luna/internal/api-server/middleware"
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

	// repos
	accountRepo := repo_impl.NewAccountRepoImpl(gormIns)
	tenantRepo := repo_impl.NewTenantRepoImpl(gormIns)

	// domain

	accountDomain := accountDomain.NewAccountDomain(accountRepo, redisIns, config, email, tenantRepo)
	// service
	tenantService := service.NewTenantService(accountDomain)

	workspaceController := controller.NewWorkspaceController(tenantService)
	v1 := g.Group("/v1")
	authV1 := v1.Group("/console/api")
	authV1.Use(middleware.TokenAuthMiddleware())
	authV1.GET("/workspaces", workspaceController.List)
	authV1.GET("/workspaces/current", workspaceController.GetTenantCurrentWorkspace)

	return nil
}

func (r *WorkspaceRoutes) GetModule() string {
	return "setup"
}
