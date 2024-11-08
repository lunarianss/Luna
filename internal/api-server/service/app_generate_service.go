package service

import (
	modelDomain "github.com/lunarianss/Luna/internal/api-server/domain/model"
	providerDomain "github.com/lunarianss/Luna/internal/api-server/domain/provider"
)

type AppGenerateService struct {
	ModelDomain    *modelDomain.ModelDomain
	ProviderDomain *providerDomain.ModelProviderDomain
}

func NewAppGenerateService(modelDomain *modelDomain.ModelDomain, providerDomain *providerDomain.ModelProviderDomain) *AppService {
	return &AppService{ModelDomain: modelDomain, ProviderDomain: providerDomain}
}

func (ags *AppGenerateService) Generate() error {
	return nil
}
