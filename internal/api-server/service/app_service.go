// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"fmt"

	"github.com/lunarianss/Luna/internal/api-server/config"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/app"
	modelDomain "github.com/lunarianss/Luna/internal/api-server/domain/model"
	providerDomain "github.com/lunarianss/Luna/internal/api-server/domain/provider"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/app"
	"github.com/lunarianss/Luna/internal/api-server/entities/base"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/pkg/template"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/internal/pkg/field"
	"github.com/lunarianss/Luna/internal/pkg/util"
	"github.com/lunarianss/Luna/pkg/errors"
	"github.com/lunarianss/Luna/pkg/log"
	"gorm.io/gorm"
)

type AppService struct {
	appDomain      *domain.AppDomain
	modelDomain    *modelDomain.ModelDomain
	providerDomain *providerDomain.ModelProviderDomain
	accountDomain  *accountDomain.AccountDomain
	db             *gorm.DB
	config         *config.Config
}

func NewAppService(appDomain *domain.AppDomain, modelDomain *modelDomain.ModelDomain, providerDomain *providerDomain.ModelProviderDomain, accountDomain *accountDomain.AccountDomain, db *gorm.DB, config *config.Config) *AppService {
	return &AppService{appDomain: appDomain, modelDomain: modelDomain, providerDomain: providerDomain, accountDomain: accountDomain, db: db, config: config}
}

func (as *AppService) CreateApp(ctx context.Context, accountID string, createAppRequest *dto.CreateAppRequest) (*dto.CreateAppResponse, error) {

	var (
		retApp *model.App
	)

	accountRecord, err := as.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, err
	}

	tenantRecord, _, err := as.accountDomain.GetCurrentTenantOfAccount(ctx, accountID)

	if err != nil {
		return nil, err
	}
	tenantID := tenantRecord.ID

	appTemplate, ok := template.DefaultAppTemplates[template.AppMode(createAppRequest.Mode)]

	if !ok {
		return nil, errors.WithCode(code.ErrAppMapMode, fmt.Sprintf("Invalid node template: %v", createAppRequest.Mode))
	}
	defaultModelConfig := &template.ModelConfig{}
	defaultModel := &template.Model{}

	util.DeepCopyUsingJSON(appTemplate.ModelConfig, defaultModelConfig)

	if defaultModelConfig.Model.Name != "" {
		modelInstance, err := as.providerDomain.GetDefaultModelInstance(ctx, tenantID, base.LLM)

		if err != nil && errors.IsCode(err, code.ErrDefaultModelNotFound) {
			log.Warnf("%s doesn't no default type of  %s model", tenantID, base.LLM)
		}

		if modelInstance != nil {
			if modelInstance.Model == defaultModelConfig.Model.Name && modelInstance.Provider == defaultModelConfig.Model.Provider {
				defaultModel = &defaultModelConfig.Model
			} else {
				modelSchema, err := as.providerDomain.GetModelSchema(ctx, modelInstance.Model, modelInstance.Credentials, modelInstance.ModelTypeInstance)
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
			provider, model, err := as.providerDomain.GetFirstProviderFirstModel(ctx, tenantID, string(base.LLM))
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
		Mode:           appTemplate.App.Mode,
		EnableSite:     field.BitBool(appTemplate.App.EnableSite),
		EnableAPI:      field.BitBool(appTemplate.App.EnableAPI),
		Icon:           createAppRequest.Icon,
		IconBackground: createAppRequest.IconBackground,
		IconType:       createAppRequest.IconType,
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

	tx := as.db.Begin()

	if defaultModelConfig.Model.Provider != "" && defaultModelConfig.Model.Name != "" {
		app, err := as.appDomain.AppRepo.CreateAppWithConfig(ctx, tx, app, appConfig)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		retApp = app
	} else {
		app, err := as.appDomain.AppRepo.CreateApp(ctx, tx, app)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		retApp = app
	}

	installApp := &model.InstalledApp{
		TenantID:         app.TenantID,
		AppID:            app.ID,
		AppOwnerTenantID: app.TenantID,
	}

	if _, err := as.appDomain.AppRunningRepo.CreateInstallApp(ctx, installApp, tx); err != nil {
		tx.Rollback()
		return nil, err
	}

	siteCode, err := as.appDomain.AppRunningRepo.GenerateUniqueCodeForSite(ctx)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	site := &model.Site{
		AppID:                  app.ID,
		Title:                  app.Name,
		IconType:               app.IconType,
		Icon:                   app.Icon,
		IconBackground:         app.IconBackground,
		DefaultLanguage:        accountRecord.InterfaceLanguage,
		CustomizeTokenStrategy: "not_allowed",
		Code:                   siteCode,
		CreatedBy:              app.CreatedBy,
		UpdatedBy:              app.UpdatedBy,
	}

	if _, err := as.appDomain.AppRunningRepo.CreateSite(ctx, site, tx); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &dto.CreateAppResponse{
		ModelConfig: appConfig,
		App:         retApp,
	}, nil
}

func (as *AppService) ListTenantApps(ctx context.Context, params *dto.ListAppRequest, accountID string) (*dto.ListAppsResponse, error) {

	tenantRecord, _, err := as.accountDomain.GetCurrentTenantOfAccount(ctx, accountID)

	if err != nil {
		return nil, err
	}
	appRecords, appCount, err := as.appDomain.AppRepo.FindTenantApps(ctx, tenantRecord, params.Page, params.PageSize)

	if err != nil {
		return nil, err
	}

	appItems := make([]*dto.ListAppItem, 0, 5)

	for _, app := range appRecords {
		appItems = append(appItems, dto.ListAppRecordToItem(app))
	}

	hasMore := 1

	if len(appRecords) < params.PageSize {
		hasMore = 0
	}

	return &dto.ListAppsResponse{
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    appCount,
		Data:     appItems,
		HasMore:  hasMore,
	}, nil

}

func (as AppService) AppDetail(ctx context.Context, appID string) (*dto.AppDetail, error) {

	appRecord, err := as.appDomain.AppRepo.GetAppByID(ctx, appID)
	if err != nil {
		return nil, err
	}

	appConfigRecord, err := as.appDomain.AppRepo.GetAppModelConfigByAppID(ctx, appID)
	if err != nil {
		return nil, err
	}

	siteRecord, err := as.appDomain.AppRunningRepo.GetSiteByAppID(ctx, appID)

	if err != nil {
		return nil, err
	}

	return dto.AppRecordToDetail(appRecord, as.config, appConfigRecord, siteRecord), nil
}

func (as *AppService) UpdateAppModelConfig(ctx context.Context, modelConfig *dto.UpdateModelConfig, appID string, accountID string) error {
	appConfig := &model.AppModelConfig{
		AppID:      appID,
		CreatedBy:  accountID,
		UpdatedBy:  accountID,
		Model:      modelConfig.Model,
		Provider:   modelConfig.Model.Provider,
		PromptType: "simple",
	}

	appConfigRecord, err := as.appDomain.AppRepo.CreateAppConfig(ctx, nil, appConfig)

	if err != nil {
		return err
	}

	appRecord, err := as.appDomain.AppRepo.GetAppByID(ctx, appID)

	if err != nil {
		return err
	}

	appRecord.AppModelConfigID = appConfigRecord.ID

	if err := as.appDomain.AppRepo.UpdateAppConfigID(ctx, appRecord); err != nil {
		return err
	}

	return nil
}
