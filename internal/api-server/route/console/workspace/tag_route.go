package route

import (
	"github.com/gin-gonic/gin"
	controller "github.com/lunarianss/Luna/internal/api-server/controller/gin/v1/tag"
	"github.com/lunarianss/Luna/internal/api-server/middleware"
)

type TagRoutes struct{}

func (a *TagRoutes) Register(g *gin.Engine) error {

	accountController := controller.NewTagController()
	v1 := g.Group("/v1")
	authV1 := v1.Group("/console/api/workspaces/current")
	authV1.Use(middleware.TokenAuthMiddleware())
	authV1.GET("/tool-providers", accountController.List)
	return nil
}

func (r *TagRoutes) GetModule() string {
	return "tag"
}
