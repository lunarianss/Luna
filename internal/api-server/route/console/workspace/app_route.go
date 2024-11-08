package route

import (
	"github.com/gin-gonic/gin"
	controller "github.com/lunarianss/Luna/internal/api-server/controller/gin/v1/app"
	"github.com/lunarianss/Luna/internal/api-server/dao"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/app"
	modelDomain "github.com/lunarianss/Luna/internal/api-server/domain/model"
	providerDomain "github.com/lunarianss/Luna/internal/api-server/domain/provider"
	"github.com/lunarianss/Luna/internal/api-server/service"
	"github.com/lunarianss/Luna/internal/pkg/mysql"
)

type AppRoutes struct{}

func (a *AppRoutes) Register(g *gin.Engine) error {
	gormIns, err := mysql.GetMySQLIns(nil)

	if err != nil {
		return err
	}

	// dao
	appDao := dao.NewAppDao(gormIns)
	modelDao := dao.NewModelDao(gormIns)
	providerDao := dao.NewModelProvider(gormIns)

	// domain
	appDomain := domain.NewAppDomain(appDao)
	modelDomain := modelDomain.NewModelDomain(modelDao)
	providerDomain := providerDomain.NewModelProviderDomain(providerDao, modelDao)

	// service
	appService := service.NewAppService(appDomain, modelDomain, providerDomain)

	appController := controller.NewAppController(appService)

	v1 := g.Group("/v1")
	modelProviderV1 := v1.Group("/console/api")
	modelProviderV1.POST("/apps", appController.Create)
	modelProviderV1.POST("/apps/appID/chat-message", appController.ChatMessage)
	return nil
}

func (r *AppRoutes) GetModule() string {
	return "app"
}
