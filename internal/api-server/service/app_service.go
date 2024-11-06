// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"fmt"

	domain "github.com/lunarianss/Luna/internal/api-server/domain/app"
	modelDomain "github.com/lunarianss/Luna/internal/api-server/domain/model"
	providerDomain "github.com/lunarianss/Luna/internal/api-server/domain/provider"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/app"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/pkg/template"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
)

type AppService struct {
	AppDomain      *domain.AppDomain
	ModelDomain    *modelDomain.ModelDomain
	ProviderDomain *providerDomain.ModelProviderDomain
}

func NewAppService(appDomain *domain.AppDomain, modelDomain *modelDomain.ModelDomain, providerDomain *providerDomain.ModelProviderDomain) *AppService {
	return &AppService{AppDomain: appDomain, ModelDomain: modelDomain, ProviderDomain: providerDomain}
}

func (as *AppService) CreateApp(ctx context.Context, tenantID, accountID string, createAppRequest *dto.CreateAppRequest) (*model.App, error) {

	appTemplate, ok := template.DefaultAppTemplates[model.AppMode(createAppRequest.Mode)]

	if !ok {
		return nil, errors.WithCode(code.ErrAppMapMode, fmt.Sprintf("Invalid node template: %v", createAppRequest.Mode))
	}

	defaultTemplateModelConfig := appTemplate.ModelConfig

	return nil, nil
}
