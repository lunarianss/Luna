package route

import (
	"github.com/gin-gonic/gin"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/_domain/account/domain_service"
	appDomain "github.com/lunarianss/Luna/internal/api-server/_domain/app/domain_service"
	webAppDomain "github.com/lunarianss/Luna/internal/api-server/_domain/web_app/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/config"
	controller "github.com/lunarianss/Luna/internal/api-server/controller/gin/v1/web/chat_app/site"
	"github.com/lunarianss/Luna/internal/api-server/middleware"
	repo_impl "github.com/lunarianss/Luna/internal/api-server/repository"
	"github.com/lunarianss/Luna/internal/api-server/service"
	"github.com/lunarianss/Luna/internal/pkg/email"
	"github.com/lunarianss/Luna/internal/pkg/mysql"
	"github.com/lunarianss/Luna/internal/pkg/redis"
)

type WebSiteRoutes struct{}

func (a *WebSiteRoutes) Register(g *gin.Engine) error {

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
	repo_impl.NewAccountRepoImpl(gormIns)
	// repos
	accountRepo := repo_impl.NewAccountRepoImpl(gormIns)
	tenantRepo := repo_impl.NewTenantRepoImpl(gormIns)
	appRepo := repo_impl.NewAppRepoImpl(gormIns)
	webAppRepo := repo_impl.NewWebAppRepoImpl(gormIns)

	// domain
	appDomain := appDomain.NewAppDomain(appRepo, webAppRepo, gormIns)
	appRunningDomain := webAppDomain.NewWebAppDomain(webAppRepo)

	accountDomain := accountDomain.NewAccountDomain(accountRepo, redisIns, config, email, tenantRepo)

	webSiteService := service.NewWebSiteService(appRunningDomain, accountDomain, appDomain, config)

	webSiteController := controller.NewWebSiteController(webSiteService)
	v1 := g.Group("/v1")
	authV1 := v1.Group("/api")
	authV1.Use(middleware.WebTokenAuthMiddleware())

	authV1.GET("/site", webSiteController.Retrieve)
	return nil
}

func (r *WebSiteRoutes) GetModule() string {
	return "web_site"
}
