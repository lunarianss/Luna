// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"errors"

	assembler "github.com/lunarianss/Luna/internal/api-server/assembler/chat"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
	"gorm.io/gorm"
)

type AnnotationService struct {
	appDomain      *appDomain.AppDomain
	providerDomain *domain_service.ProviderDomain
	accountDomain  *accountDomain.AccountDomain
	chatDomain     *chatDomain.ChatDomain
}

func NewAnnotationService(appDomain *appDomain.AppDomain, providerDomain *domain_service.ProviderDomain, accountDomain *accountDomain.AccountDomain, chatDomain *chatDomain.ChatDomain) *AnnotationService {
	return &AnnotationService{
		appDomain:      appDomain,
		providerDomain: providerDomain,
		accountDomain:  accountDomain,
		chatDomain:     chatDomain,
	}
}

func (as *AnnotationService) InsertAnnotationFromMessage(ctx context.Context, accountID string, appID string, args *dto.InsertAnnotationFormMessage) (*dto.MessageAnnotation, error) {

	var (
		bizMessageAnnotation *biz_entity.BizMessageAnnotation
	)

	accountRecord, err := as.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, err
	}

	tenant, _, err := as.accountDomain.GetCurrentTenantOfAccount(ctx, accountRecord.ID)

	if err != nil {
		return nil, err
	}

	app, err := as.appDomain.AppRepo.GetTenantApp(ctx, appID, tenant.ID)

	if err != nil {
		return nil, err
	}

	if args.MessageID != "" {
		message, err := as.chatDomain.MessageRepo.GetMessageByApp(ctx, args.MessageID, app.ID)

		if err != nil {
			return nil, err
		}

		bizMessageAnnotation, err = as.chatDomain.AnnotationRepo.GetMessageAnnotation(ctx, message.ID)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				annotation := &po_entity.MessageAnnotation{
					AppID:          app.ID,
					ConversationID: message.ConversationID,
					MessageID:      message.ID,
					Content:        args.Answer,
					Question:       args.Question,
					AccountID:      accountID,
				}
				bizMessageAnnotation, err = as.chatDomain.AnnotationRepo.CreateMessageAnnotation(ctx, annotation)

				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			bizMessageAnnotation.Content = args.Answer
			bizMessageAnnotation.Question = args.Question
			if err := as.chatDomain.AnnotationRepo.UpdateMessageAnnotation(ctx, bizMessageAnnotation.ID, bizMessageAnnotation.Content, bizMessageAnnotation.Question); err != nil {
				return nil, err
			}
		}
	} else {
		annotation := &po_entity.MessageAnnotation{
			AppID:     appID,
			Content:   args.Answer,
			Question:  args.Question,
			AccountID: accountID,
		}
		bizMessageAnnotation, err = as.chatDomain.AnnotationRepo.CreateMessageAnnotation(ctx, annotation)
		if err != nil {
			return nil, err
		}
	}

	_, err = as.chatDomain.AnnotationRepo.GetAnnotationSetting(ctx, app.ID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assembler.ConvertToAnnotation(bizMessageAnnotation), nil
		} else {
			return nil, err
		}
	} else {
		// todo rocketmq 异步消息
	}

	return nil, nil

}
