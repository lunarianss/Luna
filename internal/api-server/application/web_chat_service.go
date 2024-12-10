// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/internal/api-server/config"
	"github.com/lunarianss/Luna/internal/api-server/core/app_chat/app_chat_generator"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	biz_entity_provider "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"
	biz_entity_app_generate "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_generate"
	webAppDomain "github.com/lunarianss/Luna/internal/api-server/domain/web_app/domain_service"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime/model_registry"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
)

type WebChatService struct {
	webAppDomain   *webAppDomain.WebAppDomain
	accountDomain  *accountDomain.AccountDomain
	appDomain      *appDomain.AppDomain
	chatDomain     *chatDomain.ChatDomain
	providerDomain *domain_service.ProviderDomain
	config         *config.Config
}

func NewWebChatService(webAppDomain *webAppDomain.WebAppDomain, accountDomain *accountDomain.AccountDomain, appDomain *appDomain.AppDomain, config *config.Config, providerDomain *domain_service.ProviderDomain, chatDomain *chatDomain.ChatDomain) *WebChatService {
	return &WebChatService{
		webAppDomain:   webAppDomain,
		accountDomain:  accountDomain,
		appDomain:      appDomain,
		config:         config,
		providerDomain: providerDomain,
		chatDomain:     chatDomain,
	}
}

func (s *WebChatService) Chat(ctx context.Context, appID, endUserID string, args *dto.CreateChatMessageBody, invokeFrom biz_entity_app_generate.InvokeFrom, streaming bool) error {

	appModel, err := s.appDomain.AppRepo.GetAppByID(ctx, appID)

	if err != nil {
		return err
	}

	endUserRecord, err := s.webAppDomain.WebAppRepo.GetEndUserByID(ctx, endUserID)

	if err != nil {
		return err
	}

	chatAppGenerator := app_chat_generator.NewChatAppGenerator(s.appDomain, s.providerDomain, s.chatDomain)

	if err := chatAppGenerator.Generate(ctx, appModel, endUserRecord, args, invokeFrom, true); err != nil {
		return err
	}
	return nil
}

func (s *WebChatService) AudioToText(ctx context.Context, audioFileContent []byte, filename, appID, endUserID string) (*dto.Speech2TextResp, error) {

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

func (s *WebChatService) TextToAudio(ctx context.Context, appID, text, messageID, voice, endUserID string) error {

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
