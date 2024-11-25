// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/core/app/apps/chat"
	"github.com/lunarianss/Luna/internal/api-server/core/app/apps/entities"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/app"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat"
	providerDomain "github.com/lunarianss/Luna/internal/api-server/domain/provider"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
)

type ChatService struct {
	appDomain      *domain.AppDomain
	providerDomain *providerDomain.ModelProviderDomain
	accountDomain  *accountDomain.AccountDomain
	chatDomain     *chatDomain.ChatDomain
}

func NewChatService(appDomain *domain.AppDomain, providerDomain *providerDomain.ModelProviderDomain, accountDomain *accountDomain.AccountDomain, chatDomain *chatDomain.ChatDomain) *ChatService {
	return &ChatService{
		appDomain:      appDomain,
		providerDomain: providerDomain,
		accountDomain:  accountDomain,
		chatDomain:     chatDomain,
	}
}

func (s *ChatService) Generate(ctx context.Context, appID, accountID string, args *dto.CreateChatMessageBody, invokeFrom entities.InvokeForm, streaming bool) error {

	appModel, err := s.appDomain.AppRepo.GetAppByID(ctx, appID)

	if err != nil {
		return err
	}

	accountRecord, err := s.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return err
	}

	chatAppGenerator := chat.NewChatAppGenerator(s.appDomain, s.providerDomain, s.chatDomain)

	if err := chatAppGenerator.Generate(ctx, appModel, accountRecord, args, invokeFrom, true); err != nil {
		return err
	}

	return nil
}
