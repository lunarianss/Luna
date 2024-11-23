package controller

import "github.com/lunarianss/Luna/internal/api-server/service"

type FeatureController struct {
	featureService *service.FeatureService
}

func NewFeatureController(featureService *service.FeatureService) *FeatureController {
	return &FeatureController{
	  featureService: featureService,
	}
}
