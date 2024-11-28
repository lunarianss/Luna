// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/config"
	"github.com/lunarianss/Luna/internal/api-server/core/app_chat/app_chat_generator"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	biz_entity_app_generate "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_generate"
	webAppDomain "github.com/lunarianss/Luna/internal/api-server/domain/web_app/domain_service"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
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
