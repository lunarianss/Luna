// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"

	domain "github.com/lunarianss/Luna/internal/api-server/domain/model"
	providerDomain "github.com/lunarianss/Luna/internal/api-server/domain/provider"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
)

type ModelService struct {
	ModelDomain         *domain.ModelDomain
	ModelProviderDomain *providerDomain.ModelProviderDomain
}

func NewModelService(modelDomain *domain.ModelDomain, modelProviderDomain *providerDomain.ModelProviderDomain) *ModelService {
	return &ModelService{ModelDomain: modelDomain, ModelProviderDomain: modelProviderDomain}
}

func (ms *ModelService) SaveModelCredentials(ctx context.Context, tenantId, model, modelTpe, provider string, credentials map[string]interface{}) error {

	providerConfigurations, err := ms.ModelProviderDomain.GetConfigurations(ctx, tenantId)

	if err != nil {
		return err
	}

	providerConfiguration, ok := providerConfigurations.Configurations[provider]

	if !ok {
		return errors.WithCode(code.ErrProviderMapModel, "provider %s not found in map provider configuration", provider)
	}

	err = ms.ModelDomain.AddOrUpdateCustomModelCredentials(ctx, providerConfiguration, credentials, modelTpe, model)

	if err != nil {
		return err
	}

	return nil
}
