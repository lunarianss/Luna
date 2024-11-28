// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"time"

	"github.com/lunarianss/Luna/internal/api-server/config"
	"github.com/lunarianss/Luna/internal/api-server/core/app/apps/entities"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	providerDomain "github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	webAppDomain "github.com/lunarianss/Luna/internal/api-server/domain/web_app/domain_service"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/web_app"
	"github.com/lunarianss/Luna/pkg/errors"
	"gorm.io/gorm"
)

type WebMessageService struct {
	webAppDomain   *webAppDomain.WebAppDomain
	accountDomain  *accountDomain.AccountDomain
	appDomain      *appDomain.AppDomain
	chatDomain     *chatDomain.ChatDomain
	providerDomain *providerDomain.ProviderDomain
	config         *config.Config
}

func NewWebMessageService(webAppDomain *webAppDomain.WebAppDomain, accountDomain *accountDomain.AccountDomain, appDomain *appDomain.AppDomain, config *config.Config, providerDomain *providerDomain.ProviderDomain, chatDomain *chatDomain.ChatDomain) *WebMessageService {
	return &WebMessageService{
		webAppDomain:   webAppDomain,
		accountDomain:  accountDomain,
		appDomain:      appDomain,
		config:         config,
		providerDomain: providerDomain,
		chatDomain:     chatDomain,
	}
}

func (s *WebMessageService) ListConversations(ctx context.Context, appID, endUserID string, args *dto.ListConversationQuery, invokeFrom entities.InvokeForm) (*dto.ListConversationResponse, error) {

	var (
		includeIDs            []string
		excludeIDs            []string
		pinnedConversationIDs []string
	)

	endUser, err := s.webAppDomain.WebAppRepo.GetEndUserByID(ctx, endUserID)

	if err != nil {
		return nil, err
	}

	pinnedConversations, err := s.chatDomain.MessageRepo.FindPinnedConversationByUser(ctx, appID, endUser)

	if err != nil {
		return nil, err
	}

	for _, pinnedConversation := range pinnedConversations {
		pinnedConversationIDs = append(pinnedConversationIDs, pinnedConversation.ConversationID)
	}

	if *args.Pinned {
		if len(pinnedConversationIDs) > 0 {
			includeIDs = pinnedConversationIDs
		} else {
			includeIDs = append(includeIDs, "")
		}
	} else {
		excludeIDs = pinnedConversationIDs
	}

	conversations, count, err := s.chatDomain.MessageRepo.FindEndUserConversationsOrderByUpdated(ctx, appID, string(invokeFrom), endUser, args.Limit, includeIDs, excludeIDs, args.LastID, args.SortBy)

	if err != nil {
		return nil, err
	}

	conversationList := make([]*dto.WebConversationDetail, 0)

	for _, conversation := range conversations {
		conversationList = append(conversationList, dto.ConversationRecordToDetail(conversation))
	}

	hasMore := 0

	if len(conversations) == args.Limit {
		hasMore = 1
	}

	return &dto.ListConversationResponse{
		Data:    conversationList,
		Limit:   args.Limit,
		HasMore: hasMore,
		Count:   count,
	}, nil

}

func (s *WebMessageService) ListMessages(ctx context.Context, appID, endUserID string, args *dto.ListMessageQuery, invokeFrom entities.InvokeForm) (*dto.ListMessageResponse, error) {

	endUser, err := s.webAppDomain.WebAppRepo.GetEndUserByID(ctx, endUserID)

	if err != nil {
		return nil, err
	}

	messages, count, err := s.chatDomain.MessageRepo.FindEndUserMessages(ctx, appID, endUser, args.ConversationID, args.FirstID, args.Limit, "DESC")

	if err != nil {
		return nil, err
	}

	messagesList := make([]*dto.WebMessageDetail, 0)

	for _, message := range messages {
		messagesList = append(messagesList, dto.MessageRecordToDetail(message))
	}

	hasMore := 0

	if len(messagesList) == args.Limit {
		hasMore = 1
	}

	return &dto.ListMessageResponse{
		Data:    messagesList,
		Limit:   args.Limit,
		HasMore: hasMore,
		Count:   count,
	}, nil

}

func (s *WebMessageService) UnPinnedConversation(ctx context.Context, appID, endUserID, conversationID string) error {
	endUser, err := s.webAppDomain.WebAppRepo.GetEndUserByID(ctx, endUserID)

	if err != nil {
		return err
	}

	pinnedConversation, err := s.chatDomain.MessageRepo.GetPinnedConversationByConversation(ctx, appID, conversationID, endUser)

	if err != nil {
		return err
	}

	if err := s.chatDomain.MessageRepo.DeletePinnedConversation(ctx, pinnedConversation.ID); err != nil {
		return nil
	}
	return nil
}

func (s *WebMessageService) PinnedConversation(ctx context.Context, appID, endUserID, conversationID string) error {

	endUser, err := s.webAppDomain.WebAppRepo.GetEndUserByID(ctx, endUserID)

	if err != nil {
		return err
	}

	_, err = s.chatDomain.MessageRepo.GetPinnedConversationByConversation(ctx, appID, conversationID, endUser)

	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		conversation, err := s.chatDomain.MessageRepo.GetConversationByUser(ctx, appID, conversationID, endUser)
		if err != nil {
			return err
		}

		pinnedConversation := &po_entity.PinnedConversation{
			AppID:          appID,
			ConversationID: conversation.ID,
			CreatedByRole:  endUser.GetAccountType(),
			CreatedBy:      endUser.GetAccountID(),
		}

		if _, err := s.chatDomain.MessageRepo.CreatePinnedConversation(ctx, pinnedConversation); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func (s *WebMessageService) DeleteConversation(ctx context.Context, appID, endUserID, conversationID string) error {

	endUser, err := s.webAppDomain.WebAppRepo.GetEndUserByID(ctx, endUserID)

	if err != nil {
		return err
	}

	conversationRecord, err := s.chatDomain.MessageRepo.GetConversationByUser(ctx, appID, conversationID, endUser)

	if err != nil {
		return err
	}

	if err := s.chatDomain.MessageRepo.LogicalDeleteConversation(ctx, conversationRecord); err != nil {
		return err
	}
	return nil
}

func (s *WebMessageService) RenameConversation(ctx context.Context, appID, endUserID, conversationID string, params *dto.RenameConversationRequest) error {
	endUser, err := s.webAppDomain.WebAppRepo.GetEndUserByID(ctx, endUserID)

	if err != nil {
		return err
	}

	conversationRecord, err := s.chatDomain.MessageRepo.GetConversationByUser(ctx, appID, conversationID, endUser)

	if err != nil {
		return err
	}

	if !params.AutoGenerate {
		conversationRecord.Name = params.Name
		conversationRecord.UpdatedAt = time.Now().UTC().Unix()
		if err := s.chatDomain.MessageRepo.UpdateConversationName(ctx, conversationRecord); err != nil {
			return err
		}
	}

	return err
}
