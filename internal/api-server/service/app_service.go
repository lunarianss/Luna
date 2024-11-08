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
	"github.com/lunarianss/Luna/internal/api-server/entities/base"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/pkg/template"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/internal/pkg/util"
	"github.com/lunarianss/Luna/pkg/errors"
	"github.com/lunarianss/Luna/pkg/log"
)

type AppService struct {
	AppDomain      *domain.AppDomain
	ModelDomain    *modelDomain.ModelDomain
	ProviderDomain *providerDomain.ModelProviderDomain
}

func NewAppService(appDomain *domain.AppDomain, modelDomain *modelDomain.ModelDomain, providerDomain *providerDomain.ModelProviderDomain) *AppService {
	return &AppService{AppDomain: appDomain, ModelDomain: modelDomain, ProviderDomain: providerDomain}
}

func (as *AppService) CreateApp(ctx context.Context, tenantID, accountID string, createAppRequest *dto.CreateAppRequest) (*dto.CreateAppResponse, error) {

	var (
		retApp *model.App
	)
	appTemplate, ok := template.DefaultAppTemplates[template.AppMode(createAppRequest.Mode)]

	if !ok {
		return nil, errors.WithCode(code.ErrAppMapMode, fmt.Sprintf("Invalid node template: %v", createAppRequest.Mode))
	}
	defaultModelConfig := &template.ModelConfig{}
	defaultModel := &template.Model{}

	util.DeepCopyUsingJSON(appTemplate.ModelConfig, defaultModelConfig)

	if defaultModelConfig.Model.Name != "" {
		modelInstance, err := as.ProviderDomain.GetDefaultModelInstance(ctx, tenantID, base.LLM)

		if err != nil && errors.IsCode(err, code.ErrDefaultModelNotFound) {
			log.Warnf("%s doesn't no default type of  %s model", tenantID, base.LLM)
		}

		if modelInstance != nil {
			if modelInstance.Model == defaultModelConfig.Model.Name && modelInstance.Provider == defaultModelConfig.Model.Provider {
				defaultModel = &defaultModelConfig.Model
			} else {
				modelSchema, err := as.ProviderDomain.GetModelSchema(ctx, modelInstance.Model, modelInstance.Credentials, modelInstance.ModelTypeInstance)
				if err != nil {
					return nil, err
				}

				defaultModel.Provider = modelInstance.Provider
				defaultModel.Name = modelInstance.Model

				if v, ok := modelSchema.ModelProperties[base.MODE].(string); ok {
					defaultModel.Mode = v
				}
				defaultModel.CompletionParams = make(map[string]interface{})
			}
		} else {
			provider, model, err := as.ProviderDomain.GetFirstProviderFirstModel(ctx, tenantID, string(base.LLM))
			if err != nil {
				return nil, err
			}
			defaultModelConfig.Model.Provider = provider
			defaultModelConfig.Model.Name = model
			defaultModel = &defaultModelConfig.Model
		}
		defaultModelConfig.Model = *defaultModel
	}

	app := &model.App{
		Name:           createAppRequest.Name,
		Description:    createAppRequest.Description,
		Mode:           createAppRequest.Mode,
		Icon:           createAppRequest.Icon,
		IconBackground: createAppRequest.IconBackground,
		TenantID:       tenantID,
		CreatedBy:      accountID,
		UpdatedBy:      accountID,
		APIRpm:         createAppRequest.ApiRpm,
		APIRph:         createAppRequest.ApiRph,
	}

	appConfig := &model.AppModelConfig{
		CreatedBy:     accountID,
		UpdatedBy:     accountID,
		UserInputForm: defaultModelConfig.UserInputForm,
		PrePrompt:     defaultModelConfig.PrePrompt,
		Provider:      defaultModelConfig.Model.Provider,
		Model:         defaultModelConfig.Model,
		PromptType:    "simple",
	}

	if createAppRequest.IconType == "" {
		app.IconType = "emoji"
	}

	if defaultModelConfig.Model.Provider != "" && defaultModelConfig.Model.Name != "" {
		app, err := as.AppDomain.AppRepo.CreateAppWithConfig(ctx, app, appConfig)
		if err != nil {
			return nil, err
		}
		retApp = app
	} else {
		app, err := as.AppDomain.AppRepo.CreateApp(ctx, app)
		if err != nil {
			return nil, err
		}
		retApp = app
	}

	return &dto.CreateAppResponse{
		ModelConfig: appConfig,
		App:         retApp,
	}, nil
}
