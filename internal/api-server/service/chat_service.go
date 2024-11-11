// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/core/app/apps/chat"
	"github.com/lunarianss/Luna/internal/api-server/core/app/apps/entities"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/app"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
)

type ChatService struct {
	AppDomain *domain.AppDomain
}

func NewChatService(appDomain *domain.AppDomain) *ChatService {
	return &ChatService{
		AppDomain: appDomain,
	}
}

func (s *ChatService) Generate(ctx context.Context, appID, accountID string, args *dto.CreateChatMessageBody, invokeFrom entities.InvokeForm, streaming bool) error {

	appModel, err := s.AppDomain.AppRepo.GetAppByID(ctx, appID)

	if err != nil {
		return err
	}

	chatAppGenerator := chat.ChatAppGenerator{}

	chatAppGenerator.Generate(ctx, appModel, nil, args, invokeFrom, true)

	return nil
}
