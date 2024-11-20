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
	providerDomain "github.com/lunarianss/Luna/internal/api-server/domain/provider"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
)

type ChatService struct {
	AppDomain      *domain.AppDomain
	ProviderDomain *providerDomain.ModelProviderDomain
	AccountDomain  *accountDomain.AccountDomain
}

func NewChatService(appDomain *domain.AppDomain, providerDomain *providerDomain.ModelProviderDomain, accountDomain *accountDomain.AccountDomain) *ChatService {
	return &ChatService{
		AppDomain:      appDomain,
		ProviderDomain: providerDomain,
		AccountDomain:  accountDomain,
	}
}

func (s *ChatService) Generate(ctx context.Context, appID, accountID string, args *dto.CreateChatMessageBody, invokeFrom entities.InvokeForm, streaming bool) error {

	appModel, err := s.AppDomain.AppRepo.GetAppByID(ctx, appID)

	if err != nil {
		return err
	}

	accountRecord, err := s.AccountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return err
	}

	chatAppGenerator := &chat.ChatAppGenerator{
		AppDomain:      s.AppDomain,
		ProviderDomain: s.ProviderDomain,
	}

	if err := chatAppGenerator.Generate(ctx, appModel, accountRecord, args, invokeFrom, true); err != nil {
		return err
	}

	return nil
}
