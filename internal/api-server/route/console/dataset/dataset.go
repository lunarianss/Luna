package route

import (
	"github.com/gin-gonic/gin"
	controller "github.com/lunarianss/Luna/internal/api-server/controller/gin/v1/dataset"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/dataset"
	"github.com/lunarianss/Luna/internal/api-server/middleware"
	"github.com/lunarianss/Luna/internal/api-server/service"
)

type DatasetRoutes struct {
}

func (a *DatasetRoutes) Register(g *gin.Engine) error {
	// gormIns, err := mysql.GetMySQLIns(nil)

	// if err != nil {
	// 	return err
	// }

	// redisIns, err := redis.GetRedisIns(nil)

	// if err != nil {
	// 	return err
	// }

	// email, err := email.GetEmailSMTPIns(nil)

	// if err != nil {
	// 	return err
	// }

	// // config
	// config, err := config.GetLunaRuntimeConfig()

	// if err != nil {
	// 	return err
	// }

	// domain

	datasetDomain := domain.NewDatasetDomain()
	// service
	datasetService := service.NewDatasetService(datasetDomain)
	datasetController := controller.NewDatasetController(datasetService)

	v1 := g.Group("/v1")
	authV1 := v1.Group("/console/api")
	authV1.Use(middleware.TokenAuthMiddleware())
	authV1.GET("/files/upload", datasetController.GetFileUploadConfiguration)
	return nil
}

func (r *DatasetRoutes) GetModule() string {
	return "dataset"
}
