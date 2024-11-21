package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/api-server/config"
	controller "github.com/lunarianss/Luna/internal/api-server/controller/gin/v1/feature"
	"github.com/lunarianss/Luna/internal/api-server/service"
)

type ToolRoutes struct{}

func (a *ToolRoutes) Register(g *gin.Engine) error {
	// config
	config, err := config.GetLunaRuntimeConfig()

	if err != nil {
		return err
	}

	featureService := service.NewFeatureService(config)

	featureController := controller.NewFeatureController(featureService)
	v1 := g.Group("/v1")
	authV1 := v1.Group("/console/api")
	authV1.GET("/system-features", featureController.GetSystemConfigs)
	authV1.GET("/features", featureController.List)
	return nil
}

func (r *ToolRoutes) GetModule() string {
	return "setup"
}
