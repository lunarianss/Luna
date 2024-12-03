// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"

	assembler "github.com/lunarianss/Luna/internal/api-server/assembler/chat"
	"github.com/lunarianss/Luna/internal/api-server/core/app_chat/app_chat_generator"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	biz_entity_app_generate "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_generate"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
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

func (s *ChatService) ListConsoleMessagesOfConversation(ctx context.Context, appID string, args *dto.ListChatMessageQuery) (*dto.ListChatMessagesResponse, error) {
	conversation, err := s.chatDomain.MessageRepo.GetConversationByApp(ctx, args.ConversationID, appID)

	if err != nil {
		return nil, err
	}

	messageRecords, count, err := s.chatDomain.MessageRepo.FindConsoleAppMessages(ctx, conversation.ID, args.Limit)

	if err != nil {
		return nil, err
	}
	var messageItems []*dto.ListChatMessageItem

	hasMore := true

	for _, mr := range messageRecords {
		messageDto := assembler.ConvertToListMessageDto(mr)
		messageItems = append(messageItems, messageDto)
	}

	if len(messageRecords) < 10 {
		hasMore = false
	}

	util.SliceReverse(messageItems)

	return &dto.ListChatMessagesResponse{
		Limit:   args.Limit,
		HasMore: hasMore,
		Data:    messageItems,
		Count:   count,
	}, nil
}

func (s *ChatService) ListConversations(ctx context.Context, accountID string, appID string, args *dto.ListChatConversationQuery) (*dto.ListChatConversationResponse, error) {
	var rets []*dto.ListChatConversationItem
	var sessionID string

	conversationRecords, count, err := s.chatDomain.MessageRepo.FindConversationsInConsole(ctx, args.Page, args.Limit, appID, args.Start, args.End, args.SortBy, args.Keyword)

	if err != nil {
		return nil, err
	}

	for _, conversationRecord := range conversationRecords {
		conversationJoin := assembler.ConvertToConversationJoins(conversationRecord)

		msgCount, err := s.chatDomain.MessageRepo.GetMessageCountOfConversation(ctx, conversationRecord.ID)

		if err != nil {
			return nil, err
		}

		account, err := s.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

		if err != nil {
			return nil, err
		}

		if conversationRecord.FromEndUserID != "" {
			endUser, err := s.appDomain.WebAppRepo.GetEndUserByID(ctx, conversationRecord.FromEndUserID)

			if err != nil {
				return nil, err
			}
			sessionID = endUser.SessionID
		}

		if err != nil {
			return nil, err
		}

		conversationJoin.MessageCount = msgCount
		conversationJoin.FromAccountName = account.Name
		conversationJoin.UserFeedbackStats = dto.NewFeedBackStats()
		conversationJoin.AdminFeedbackStats = dto.NewFeedBackStats()

		if conversationRecord.FromEndUserID != "" {
			conversationJoin.FromEndUserSessionID = sessionID
		}

		rets = append(rets, conversationJoin)

	}

	hasMore := false

	if len(conversationRecords) == args.Limit {
		hasMore = true
	}

	return &dto.ListChatConversationResponse{
		Page:    args.Page,
		Limit:   args.Limit,
		Data:    rets,
		HasMore: hasMore,
		Total:   count,
	}, nil
}
