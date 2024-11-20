package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/api-server/config"
	controller "github.com/lunarianss/Luna/internal/api-server/controller/gin/v1/feature"
	"github.com/lunarianss/Luna/internal/api-server/middleware"
	"github.com/lunarianss/Luna/internal/api-server/service"
)

type FeatureRoutes struct{}

func (a *FeatureRoutes) Register(g *gin.Engine) error {
	// config
	config, err := config.GetLunaRuntimeConfig()

	if err != nil {
		return err
	}

	featureService := service.NewFeatureService(config)

	featureController := controller.NewFeatureController(featureService)
	v1 := g.Group("/v1")
	authV1 := v1.Group("/console/api")
	authV1.GET("/system-features", middleware.TokenAuthMiddleware(), featureController.GetSystemConfigs)
	return nil
}

func (r *FeatureRoutes) GetModule() string {
	return "setup"
}
