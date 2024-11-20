package controller

import "github.com/lunarianss/Luna/internal/api-server/service"

type FeatureController struct {
	FeatureService *service.FeatureService
}

func NewFeatureController(featureService *service.FeatureService) *FeatureController {
	return &FeatureController{
		FeatureService: featureService,
	}
}
