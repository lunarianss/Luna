// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/core/app_chat/app_chat_generator"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	biz_entity_app_generate "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_generate"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
)

type ChatService struct {
	appDomain      *appDomain.AppDomain
	providerDomain *domain_service.ProviderDomain
	accountDomain  *accountDomain.AccountDomain
	chatDomain     *chatDomain.ChatDomain
}

func NewChatService(appDomain *appDomain.AppDomain, providerDomain *domain_service.ProviderDomain, accountDomain *accountDomain.AccountDomain, chatDomain *chatDomain.ChatDomain) *ChatService {
	return &ChatService{
		appDomain:      appDomain,
		providerDomain: providerDomain,
		accountDomain:  accountDomain,
		chatDomain:     chatDomain,
	}
}

func (s *ChatService) Generate(ctx context.Context, appID, accountID string, args *dto.CreateChatMessageBody, invokeFrom biz_entity_app_generate.InvokeFrom, streaming bool) error {

	appModel, err := s.appDomain.AppRepo.GetAppByID(ctx, appID)

	if err != nil {
		return err
	}

	accountRecord, err := s.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return err
	}

	chatAppGenerator := app_chat_generator.NewChatAppGenerator(s.appDomain, s.providerDomain, s.chatDomain)

	if err := chatAppGenerator.Generate(ctx, appModel, accountRecord, args, invokeFrom, true); err != nil {
		return err
	}

	return nil
}
