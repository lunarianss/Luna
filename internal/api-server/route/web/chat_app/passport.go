package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/api-server/config"
	controller "github.com/lunarianss/Luna/internal/api-server/controller/gin/v1/web/chat_app/passport"
	"github.com/lunarianss/Luna/internal/api-server/dao"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/app"
	appRunningDomain "github.com/lunarianss/Luna/internal/api-server/domain/app_running"
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

	jwt := jwt.GetJWTIns()

	// dao
	appDao := dao.NewAppDao(gormIns)
	appRunningDao := dao.NewAppRunningDao(gormIns)

	// domain
	appDomain := domain.NewAppDomain(appDao, appRunningDao)
	appRunningDomain := appRunningDomain.NewAppRunningDomain(appRunningDao)

	passportService := service.NewPassportService(appRunningDomain, appDomain, config, jwt)
	passportController := controller.NewPassportController(passportService)
	v1 := g.Group("/v1")
	authV1 := v1.Group("/api")
	authV1.GET("/passport", passportController.Acquire)
	return nil
}

func (r *PassportRoutes) GetModule() string {
	return "passport"
}
