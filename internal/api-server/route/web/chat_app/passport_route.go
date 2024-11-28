package route

import (
	"github.com/gin-gonic/gin"
	appDomain "github.com/lunarianss/Luna/internal/api-server/_domain/app/domain_service"
	webAppDomain "github.com/lunarianss/Luna/internal/api-server/_domain/web_app/domain_service"
	repo_impl "github.com/lunarianss/Luna/internal/api-server/_repo"
	"github.com/lunarianss/Luna/internal/api-server/config"
	controller "github.com/lunarianss/Luna/internal/api-server/controller/gin/v1/web/chat_app/passport"

	"github.com/lunarianss/Luna/internal/api-server/service"
	"github.com/lunarianss/Luna/internal/pkg/jwt"
	"github.com/lunarianss/Luna/internal/pkg/mysql"
)

type PassportRoutes struct{}

func (a *PassportRoutes) Register(g *gin.Engine) error {

	gormIns, err := mysql.GetMySQLIns(nil)

	if err != nil {
		return err
	}

	// config
	config, err := config.GetLunaRuntimeConfig()

	if err != nil {
		return err
	}

	jwt, err := jwt.GetJWTIns()

	if err != nil {
		return err
	}
	// repos
	appRepo := repo_impl.NewAppRepoImpl(gormIns)
	webAppRepo := repo_impl.NewWebAppRepoImpl(gormIns)

	// domain
	appDomain := appDomain.NewAppDomain(appRepo, webAppRepo, gormIns)
	webAppDomain := webAppDomain.NewWebAppDomain(webAppRepo)

	passportService := service.NewPassportService(webAppDomain, appDomain, config, jwt)
	passportController := controller.NewPassportController(passportService)
	v1 := g.Group("/v1")
	authV1 := v1.Group("/api")
	authV1.GET("/passport", passportController.Acquire)
	return nil
}

func (r *PassportRoutes) GetModule() string {
	return "passport"
}
