// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/infrastructure/log"
	assembler "github.com/lunarianss/Luna/internal/api-server/assembler/app"
	"github.com/lunarianss/Luna/internal/api-server/config"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/provider_app_config"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	biz_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	po_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/app"
	chatDto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime/model_registry"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/lunarianss/Luna/internal/infrastructure/field"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
	"gorm.io/gorm"
)

const MAX_KEYS = 10

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

	tenantRecord, tenantJoin, err := as.accountDomain.GetCurrentTenantOfAccount(ctx, accountRecord.ID)

	if err != nil {
		return nil, err
	}
	if !tenantJoin.IsEditor() {
		return nil, errors.WithCode(code.ErrForbidden, fmt.Sprintf("You don't have the permission for tenant %s", tenantRecord.Name))
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
			log.Warnf("tenant %s doesn't no default type of  %s model", tenantID, common.LLM)
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

func (as AppService) AppDetail(ctx context.Context, accountID string, appID string) (*dto.AppDetail, error) {
	tenantRecord, _, err := as.accountDomain.GetCurrentTenantOfAccount(ctx, accountID)

	if err != nil {
		return nil, err
	}

	appRecord, err := as.appDomain.AppRepo.GetTenantApp(ctx, appID, tenantRecord.ID)

	if err != nil {
		return nil, err
	}

	appConfigRecord, err := as.appDomain.AppRepo.GetAppModelConfigById(ctx, appRecord.AppModelConfigID, appID)
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

	tenant, tenantJoin, err := as.accountDomain.GetCurrentTenantOfAccount(ctx, accountID)

	if err != nil {
		return err
	}

	if !tenantJoin.IsEditor() {
		return errors.WithCode(code.ErrForbidden, fmt.Sprintf("You don't have the permission for %s", tenant.Name))
	}

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

	appRecord, err := as.appDomain.AppRepo.GetTenantApp(ctx, appID, tenant.ID)

	if err != nil {
		return err
	}

	appRecord.AppModelConfigID = appConfigRecord.ID

	if err := as.appDomain.AppRepo.UpdateAppConfigID(ctx, appRecord); err != nil {
		return err
	}

	return nil
}

func (as *AppService) UpdateEnableAppSite(ctx context.Context, accountID string, appID string, enable_site bool) (*dto.AppDetail, error) {
	tenant, tenantJoin, err := as.accountDomain.GetCurrentTenantOfAccount(ctx, accountID)

	if err != nil {
		return nil, err
	}

	if !tenantJoin.IsEditor() {
		return nil, errors.WithCode(code.ErrForbidden, fmt.Sprintf("You don't have the permission for %s", tenant.Name))
	}

	appRecord, err := as.appDomain.AppRepo.GetTenantApp(ctx, appID, tenant.ID)

	if err != nil {
		return nil, err
	}

	appConfigRecord, err := as.appDomain.AppRepo.GetAppModelConfigById(ctx, appRecord.AppModelConfigID, appID)
	if err != nil {
		return nil, err
	}

	siteRecord, err := as.appDomain.WebAppRepo.GetSiteByAppID(ctx, appID)

	if err != nil {
		return nil, err
	}

	enableSite := util.BoolToInt(enable_site)

	if appRecord.EnableSite == field.BitBool(enableSite) {
		return dto.AppRecordToDetail(appRecord, as.config, appConfigRecord, siteRecord), nil
	}

	appRecord.EnableSite = field.BitBool(enableSite)
	appRecord.UpdatedBy = accountID
	appRecord.UpdatedAt = int(time.Now().UTC().Unix())

	appRecord, err = as.appDomain.AppRepo.UpdateEnableAppSite(ctx, appRecord)

	if err != nil {
		return nil, err
	}

	return dto.AppRecordToDetail(appRecord, as.config, appConfigRecord, siteRecord), nil
}

func (as *AppService) UpdateEnableAppApi(ctx context.Context, accountID string, appID string, enable_api bool) (*dto.AppDetail, error) {
	tenant, tenantJoin, err := as.accountDomain.GetCurrentTenantOfAccount(ctx, accountID)

	if err != nil {
		return nil, err
	}

	if !tenantJoin.IsEditor() {
		return nil, errors.WithCode(code.ErrForbidden, fmt.Sprintf("You don't have the permission for %s", tenant.Name))
	}

	appRecord, err := as.appDomain.AppRepo.GetTenantApp(ctx, appID, tenant.ID)

	if err != nil {
		return nil, err
	}

	appConfigRecord, err := as.appDomain.AppRepo.GetAppModelConfigById(ctx, appRecord.AppModelConfigID, appID)
	if err != nil {
		return nil, err
	}

	siteRecord, err := as.appDomain.WebAppRepo.GetSiteByAppID(ctx, appID)

	if err != nil {
		return nil, err
	}

	enableAPI := util.BoolToInt(enable_api)

	if appRecord.EnableSite == field.BitBool(enableAPI) {
		return dto.AppRecordToDetail(appRecord, as.config, appConfigRecord, siteRecord), nil
	}

	appRecord.EnableAPI = field.BitBool(enableAPI)
	appRecord.UpdatedBy = accountID
	appRecord.UpdatedAt = int(time.Now().UTC().Unix())

	appRecord, err = as.appDomain.AppRepo.UpdateEnableAppApi(ctx, appRecord)

	if err != nil {
		return nil, err
	}

	return dto.AppRecordToDetail(appRecord, as.config, appConfigRecord, siteRecord), nil
}

type TaskData struct {
	TaskDescription string
	InputText       string
}

func (as *AppService) GeneratePrompt(ctx context.Context, accountID string, args *dto.GeneratePrompt, rule_config_max_tokens int) (*dto.GeneratePromptResponse, error) {
	tenant, _, err := as.accountDomain.GetCurrentTenantOfAccount(ctx, accountID)

	if err != nil {
		return nil, err
	}

	modelParameters := map[string]interface {
	}{"max_tokens": rule_config_max_tokens,
		"temperature": 0.01}

	modelIns, err := as.providerDomain.GetModelInstance(ctx, tenant.ID, args.ModelConfig.Provider, args.ModelConfig.Name, common.LLM)

	if err != nil {
		return nil, err
	}

	promptGenerate := biz_entity_chat.GetRuleConfigPromptGenerateTemplate()

	promptTemplateParse := biz_entity.NewPromptTemplateParse(promptGenerate, false)

	promptGenerate = promptTemplateParse.Format(map[string]interface{}{
		"TaskDescription": args.Instruction,
	}, false)

	var promptMessages []*po_entity_chat.PromptMessage

	promptMessages = append(promptMessages, po_entity_chat.NewUserMessage(promptGenerate))

	modelCaller := model_registry.NewModelRegisterCaller(args.ModelConfig.Name, string(common.LLM), args.ModelConfig.Provider, modelIns.Credentials, modelIns.ModelTypeInstance)

	llmResult, err := modelCaller.InvokeLLMNonStream(ctx, promptMessages, modelParameters, nil, nil, accountID, nil)

	if err != nil {
		return nil, err
	}

	return &dto.GeneratePromptResponse{
		OpenStatement: "",
		Variables:     make([]string, 0),
		Prompt:        llmResult.Message.Content.(string),
		Error:         "",
	}, nil

}

func (as *AppService) GenerateServiceToken(ctx context.Context, accountID string, appID string) (*dto.GenerateServiceToken, error) {

	tenant, tenantJoin, err := as.accountDomain.GetCurrentTenantOfAccount(ctx, accountID)

	if err != nil {
		return nil, err
	}

	if !tenantJoin.IsEditor() {
		return nil, errors.WithCode(code.ErrForbidden, fmt.Sprintf("You don't have the permission for %s", tenant.Name))
	}

	app, err := as.appDomain.AppRepo.GetTenantApp(ctx, appID, tenant.ID)

	if err != nil {
		return nil, err
	}

	count, err := as.appDomain.AppRepo.GetServiceTokenCount(ctx, app.ID)

	if err != nil {
		return nil, err
	}
	if count >= MAX_KEYS {
		return nil, errors.WithCode(code.ErrAppTokenExceed, "count %d > maxSize %d, exceed max account", count, MAX_KEYS)
	}

	serviceToken, err := as.appDomain.AppRepo.GenerateServiceToken(ctx, 24)

	if err != nil {
		return nil, err
	}

	appToken := &po_entity.ApiToken{
		Token:    serviceToken,
		Type:     "app",
		TenantID: tenant.ID,
		AppID:    appID,
	}

	appTokenRecord, err := as.appDomain.AppRepo.CreateServiceToken(ctx, appToken)

	if err != nil {
		return nil, err
	}

	return &dto.GenerateServiceToken{ID: appToken.AppID, Type: appTokenRecord.Type, Token: appTokenRecord.Token, CreatedAt: appTokenRecord.CreatedAt}, nil
}

func (as *AppService) ListServiceTokens(ctx context.Context, accountID string, appID string) (*dto.DataWrapperResponse[[]*dto.GenerateServiceToken], error) {
	tenant, _, err := as.accountDomain.GetCurrentTenantOfAccount(ctx, accountID)

	if err != nil {
		return nil, err
	}

	app, err := as.appDomain.AppRepo.GetTenantApp(ctx, appID, tenant.ID)

	if err != nil {
		return nil, err
	}

	apiTokens, err := as.appDomain.AppRepo.FindServiceTokens(ctx, app.ID)

	if err != nil {
		return nil, err
	}

	dtoApiTokens := assembler.ConvertToServiceTokens(apiTokens)

	return &dto.DataWrapperResponse[[]*dto.GenerateServiceToken]{
		Data: dtoApiTokens,
	}, nil
}
