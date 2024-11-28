// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"github.com/lunarianss/Luna/internal/api-server/config"

	dto "github.com/lunarianss/Luna/internal/api-server/dto/system"
	"github.com/lunarianss/Luna/internal/infrastructure/options"
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
