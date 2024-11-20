package route

import (
	"github.com/gin-gonic/gin"
	controller "github.com/lunarianss/Luna/internal/api-server/controller/gin/v1/model-provider/model"
	"github.com/lunarianss/Luna/internal/api-server/dao"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/model"
	providerDomain "github.com/lunarianss/Luna/internal/api-server/domain/provider"
	"github.com/lunarianss/Luna/internal/api-server/service"
	"github.com/lunarianss/Luna/internal/pkg/mysql"
)

type ModelRoutes struct{}

func (r *ModelRoutes) Register(g *gin.Engine) error {
	gormIns, err := mysql.GetMySQLIns(nil)

	if err != nil {
		return err
	}

	// dao
	modelDao := dao.NewModelDao(gormIns)
	modelProviderDao := dao.NewModelProvider(gormIns)

	// domain
	modelDomain := domain.NewModelDomain(modelDao)
	modelProviderDomain := providerDomain.NewModelProviderDomain(modelProviderDao, modelDao)

	// service
	modelService := service.NewModelService(modelDomain, modelProviderDomain)
	modelController := controller.NewModelController(modelService)

	v1 := g.Group("/v1")
	modelProviderV1 := v1.Group("/console/api/workspaces/current")

	modelProviderV1.POST("/model-providers/:provider/models", modelController.SaveModelCredential)

	modelProviderV1.POST("/model-types/:modelType", modelController.GetAccountAvailableModels)

	return nil
}

func (r *ModelRoutes) GetModule() string {
	return "model"
}
