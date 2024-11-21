package service

import (
	"github.com/lunarianss/Luna/internal/api-server/config"

	dto "github.com/lunarianss/Luna/internal/api-server/dto/system"
	"github.com/lunarianss/Luna/internal/pkg/options"
)

type FeatureService struct {
	config *config.Config
}

func NewFeatureService(config *config.Config) *FeatureService {
	return &FeatureService{
		config: config,
	}
}

func (s *FeatureService) GetSystemConfig() (*options.SystemOptions, error) {
	return s.config.SystemOptions, nil
}

func (s *FeatureService) ListFeatures() (*dto.FeatureModel, error) {
	return dto.NewFeatureModel(), nil
}
