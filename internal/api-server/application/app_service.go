// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/infrastructure/log"
	assembler "github.com/lunarianss/Luna/internal/api-server/assembler/app"
	"github.com/lunarianss/Luna/internal/api-server/config"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/provider_app_config"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/app"
	chatDto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/lunarianss/Luna/internal/infrastructure/field"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
	"gorm.io/gorm"
)

type AppService struct {
	appDomain      *appDomain.AppDomain
	providerDomain *domain_service.ProviderDomain
	accountDomain  *accountDomain.AccountDomain
	db             *gorm.DB
	config         *config.Config
}

func NewAppService(appDomain *appDomain.AppDomain, providerDomain *domain_service.ProviderDomain, accountDomain *accountDomain.AccountDomain, db *gorm.DB, config *config.Config) *AppService {
	return &AppService{appDomain: appDomain, providerDomain: providerDomain, accountDomain: accountDomain, db: db, config: config}
}

func (as *AppService) CreateApp(ctx context.Context, accountID string, createAppRequest *dto.CreateAppRequest) (*dto.CreateAppResponse, error) {

	accountRecord, err := as.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, err
	}

	tenantRecord, _, err := as.accountDomain.GetCurrentTenantOfAccount(ctx, accountID)

	if err != nil {
		return nil, err
	}
	tenantID := tenantRecord.ID

	appTemplate, err := as.appDomain.GetTemplate(ctx, createAppRequest.Mode)

	if err != nil {
		return nil, err
	}

	defaultModelConfig := &biz_entity.ModelConfig{}
	defaultModel := &biz_entity.ModelInfo{}

	util.DeepCopyUsingJSON(appTemplate.ModelConfig, defaultModelConfig)

	if defaultModelConfig.Model.Name != "" {
		modelInstance, err := as.providerDomain.GetDefaultModelInstance(ctx, tenantID, common.LLM)

		if err != nil && errors.IsCode(err, code.ErrDefaultModelNotFound) {
			log.Warnf("%s doesn't no default type of  %s model", tenantID, common.LLM)
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

				if v, ok := modelSchema.ModelProperties[common.MODE].(string); ok {
					defaultModel.Mode = v
				}
				defaultModel.CompletionParams = make(map[string]interface{})
			}
		} else {
			provider, model, err := as.providerDomain.GetFirstProviderFirstModel(ctx, tenantID, string(common.LLM))
			if err != nil {
				return nil, err
			}
			defaultModelConfig.Model.Provider = provider
			defaultModelConfig.Model.Name = model
			defaultModel = &defaultModelConfig.Model
		}
		defaultModelConfig.Model = *defaultModel
	}

	app := &po_entity.App{
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

	appConfig := &po_entity.AppModelConfig{
		CreatedBy:     accountID,
		UpdatedBy:     accountID,
		UserInputForm: biz_entity.ConvertToUserInputPoEntity(defaultModelConfig.UserInputForm),
		PrePrompt:     defaultModelConfig.PrePrompt,
		Provider:      defaultModelConfig.Model.Provider,
		Model:         biz_entity.ConvertToModelPoEntity(defaultModelConfig.Model),
		PromptType:    "simple",
	}

	if createAppRequest.IconType == "" {
		app.IconType = "emoji"
	}

	app, appConfig, err = as.appDomain.CreateApp(ctx, app, appConfig, defaultModelConfig.Model.Provider, defaultModelConfig.Model.Name, accountRecord.InterfaceLanguage)

	if err != nil {
		return nil, err
	}

	return &dto.CreateAppResponse{
		ModelConfig: appConfig,
		App:         app,
	}, nil
}

func (as *AppService) ListTenantApps(ctx context.Context, params *dto.ListAppRequest, accountID string) (*dto.ListAppsResponse, error) {

	tenantRecord, _, err := as.accountDomain.GetCurrentTenantOfAccount(ctx, accountID)

	if err != nil {
		return nil, err
	}
	appRecords, appCount, err := as.appDomain.AppRepo.FindTenantApps(ctx, tenantRecord.ID, params.Page, params.PageSize)

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

	appConfigRecord, err := as.appDomain.AppRepo.GetAppModelConfigById(ctx, appRecord.AppModelConfigID)
	if err != nil {
		return nil, err
	}

	siteRecord, err := as.appDomain.WebAppRepo.GetSiteByAppID(ctx, appID)

	if err != nil {
		return nil, err
	}

	return dto.AppRecordToDetail(appRecord, as.config, appConfigRecord, siteRecord), nil
}

func (as *AppService) UpdateAppModelConfig(ctx context.Context, modelConfig *chatDto.AppModelConfigDto, appID string, accountID string) error {
	configEntity := assembler.ConvertToConfigEntity(modelConfig)
	configRecord := configEntity.ConvertToAppConfigPoEntity()
	configRecord.AppID = appID
	configRecord.Provider = configEntity.Model.Provider
	configRecord.UpdatedBy = accountID
	configRecord.CreatedBy = accountID

	appConfigRecord, err := as.appDomain.AppRepo.CreateAppConfig(ctx, nil, configRecord)

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
