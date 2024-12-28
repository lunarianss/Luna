// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"

	"github.com/lunarianss/Luna/infrastructure/errors"
	assembler "github.com/lunarianss/Luna/internal/api-server/assembler/service"
	"github.com/lunarianss/Luna/internal/api-server/config"
	"github.com/lunarianss/Luna/internal/api-server/core/app_chat/app_chat_generator"
	"github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_registry"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	po_account "github.com/lunarianss/Luna/internal/api-server/domain/account/entity/po_entity"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	datasetDomain "github.com/lunarianss/Luna/internal/api-server/domain/dataset/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	biz_entity_provider "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"
	biz_entity_app_generate "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_generate"
	webAppDomain "github.com/lunarianss/Luna/internal/api-server/domain/web_app/domain_service"
	po_webapp "github.com/lunarianss/Luna/internal/api-server/domain/web_app/entity/po_entity"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ServiceChatService struct {
	webAppDomain   *webAppDomain.WebAppDomain
	accountDomain  *accountDomain.AccountDomain
	appDomain      *appDomain.AppDomain
	chatDomain     *chatDomain.ChatDomain
	providerDomain *domain_service.ProviderDomain
	config         *config.Config
	datasetDomain  *datasetDomain.DatasetDomain
	redis          *redis.Client
}

func NewServiceChatService(webAppDomain *webAppDomain.WebAppDomain, accountDomain *accountDomain.AccountDomain, appDomain *appDomain.AppDomain, config *config.Config, providerDomain *domain_service.ProviderDomain, chatDomain *chatDomain.ChatDomain) *ServiceChatService {
	return &ServiceChatService{
		webAppDomain:   webAppDomain,
		accountDomain:  accountDomain,
		appDomain:      appDomain,
		config:         config,
		providerDomain: providerDomain,
		chatDomain:     chatDomain,
	}
}

func (s *ServiceChatService) baseChat(ctx context.Context, app *po_entity.App, tenant *po_account.Tenant, args *dto.ServiceCreateChatMessageBody) (*po_webapp.EndUser, app_chat_generator.IChatAppGenerator, *dto.CreateChatMessageBody, error) {

	var (
		endUserRecord *po_webapp.EndUser
		err           error
	)

	endUserRecord, err = s.webAppDomain.WebAppRepo.GetEndUserByInfo(ctx, args.User, tenant.ID, app.ID, "service_api")

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			endUser := &po_webapp.EndUser{
				TenantID:    tenant.ID,
				AppID:       app.ID,
				Type:        "service_api",
				SessionID:   args.User,
				IsAnonymous: 0,
			}

			endUserRecord, err = s.webAppDomain.WebAppRepo.CreateEndUser(ctx, endUser, nil)
			if err != nil {
				return nil, nil, nil, errors.WithCode(code.ErrDatabase, "create end user error: %s", err.Error())
			}
		} else {
			return nil, nil, nil, errors.WithCode(code.ErrDatabase, "get end user error: %s", err.Error())
		}
	}

	chatAppGenerator := app_chat_generator.NewChatAppGenerator(s.appDomain, s.providerDomain, s.chatDomain, s.datasetDomain, s.redis)

	chatMessageBodyDto := assembler.ConvertToCreateChatMessageBody(args)

	return endUserRecord, chatAppGenerator, chatMessageBodyDto, nil
}

func (s *ServiceChatService) Chat(ctx context.Context, app *po_entity.App, tenant *po_account.Tenant, args *dto.ServiceCreateChatMessageBody, invokeFrom biz_entity_app_generate.InvokeFrom) error {

	endUserRecord, chatAppGenerator, chatMessageBodyDto, err := s.baseChat(ctx, app, tenant, args)

	if err != nil {
		return err
	}

	if err := chatAppGenerator.Generate(ctx, app, endUserRecord, chatMessageBodyDto, invokeFrom, true); err != nil {
		return err
	}

	return nil
}

func (s *ServiceChatService) ChatNonStream(ctx context.Context, app *po_entity.App, tenant *po_account.Tenant, args *dto.ServiceCreateChatMessageBody, invokeFrom biz_entity_app_generate.InvokeFrom) (*dto.ServiceChatCompletionResponse, error) {

	endUserRecord, chatAppGenerator, chatMessageBodyDto, err := s.baseChat(ctx, app, tenant, args)

	if err != nil {
		return nil, err
	}

	return chatAppGenerator.GenerateNonStream(ctx, app, endUserRecord, chatMessageBodyDto, invokeFrom, true)
}

func (s *ServiceChatService) AudioToText(ctx context.Context, audioFileContent []byte, filename, appID, endUserID string) (*dto.Speech2TextResp, error) {

	appModel, err := s.appDomain.AppRepo.GetAppByID(ctx, appID)

	if err != nil {
		return nil, err
	}

	endUserRecord, err := s.webAppDomain.WebAppRepo.GetEndUserByID(ctx, endUserID)

	if err != nil {
		return nil, err
	}

	audioModelIntegratedInstance, err := s.providerDomain.GetDefaultModelInstance(ctx, appModel.TenantID, biz_entity_provider.SPEECH2TEXT)

	if err != nil {
		return nil, err
	}

	modelRegistryCaller := model_registry.NewModelRegisterCaller(audioModelIntegratedInstance.Model, string(biz_entity_provider.SPEECH2TEXT), audioModelIntegratedInstance.Provider, audioModelIntegratedInstance.Credentials, audioModelIntegratedInstance.ModelTypeInstance)

	transStr, err := modelRegistryCaller.InvokeSpeechToText(ctx, audioFileContent, endUserRecord.ID, filename)

	if err != nil {
		return nil, err
	}

	return &dto.Speech2TextResp{
		Text: transStr,
	}, nil
}

func (s *ServiceChatService) TextToAudio(ctx context.Context, appID, text, messageID, voice, endUserID string) error {

	appModel, err := s.appDomain.AppRepo.GetAppByID(ctx, appID)

	if err != nil {
		return err
	}

	endUserRecord, err := s.appDomain.WebAppRepo.GetEndUserByID(ctx, endUserID)

	if err != nil {
		return err
	}
	if messageID != "" {
		message, err := s.chatDomain.MessageRepo.GetMessageByID(ctx, messageID)

		if err != nil {
			return err
		}

		if message.Answer == "" && message.Status == "normal" {
			return errors.WithCode(code.ErrAudioTextEmpty, "")
		}

		ttsModelIntegratedInstance, err := s.providerDomain.GetDefaultModelInstance(ctx, appModel.TenantID, biz_entity_provider.TTS)

		if err != nil {
			return err
		}

		modelRegistryCaller := model_registry.NewModelRegisterCaller(ttsModelIntegratedInstance.Model, string(biz_entity_provider.TTS), ttsModelIntegratedInstance.Provider, ttsModelIntegratedInstance.Credentials, ttsModelIntegratedInstance.ModelTypeInstance)

		err = modelRegistryCaller.InvokeTextToSpeech(ctx, nil, endUserRecord.ID, "longxiaochun", "", []string{text})

		if err != nil {
			return err
		}
	}
	return nil
}
