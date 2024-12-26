// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"

	"github.com/lunarianss/Luna/infrastructure/errors"
	assembler "github.com/lunarianss/Luna/internal/api-server/assembler/chat"
	"github.com/lunarianss/Luna/internal/api-server/core/app_chat/app_chat_generator"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/provider_app_config"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	biz_entity_provider "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"
	biz_entity_app_generate "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_generate"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime/model_registry"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
	"gorm.io/gorm"
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

func (s *ChatService) TextToAudio(ctx context.Context, appID, text, messageID, voice, accountID string) error {
	accountRecord, err := s.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return err
	}

	tenant, _, err := s.accountDomain.GetCurrentTenantOfAccount(ctx, accountRecord.ID)

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

		ttsModelIntegratedInstance, err := s.providerDomain.GetDefaultModelInstance(ctx, tenant.ID, biz_entity_provider.TTS)

		if err != nil {
			return err
		}

		modelRegistryCaller := model_registry.NewModelRegisterCaller(ttsModelIntegratedInstance.Model, string(biz_entity_provider.TTS), ttsModelIntegratedInstance.Provider, ttsModelIntegratedInstance.Credentials, ttsModelIntegratedInstance.ModelTypeInstance)

		err = modelRegistryCaller.InvokeTextToSpeech(ctx, nil, accountRecord.ID, "longxiaochun", "", []string{text})

		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ChatService) AudioToText(ctx context.Context, audioFileContent []byte, filename, appID, accountID string) (*dto.Speech2TextResp, error) {

	accountRecord, err := s.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, err
	}

	tenant, _, err := s.accountDomain.GetCurrentTenantOfAccount(ctx, accountRecord.ID)

	if err != nil {
		return nil, err
	}

	audioModelIntegratedInstance, err := s.providerDomain.GetDefaultModelInstance(ctx, tenant.ID, biz_entity_provider.SPEECH2TEXT)

	if err != nil {
		return nil, err
	}

	modelRegistryCaller := model_registry.NewModelRegisterCaller(audioModelIntegratedInstance.Model, string(biz_entity_provider.SPEECH2TEXT), audioModelIntegratedInstance.Provider, audioModelIntegratedInstance.Credentials, audioModelIntegratedInstance.ModelTypeInstance)

	transStr, err := modelRegistryCaller.InvokeSpeechToText(ctx, audioFileContent, accountID, filename)

	if err != nil {
		return nil, err
	}

	return &dto.Speech2TextResp{
		Text: transStr,
	}, nil
}

func (s *ChatService) Generate(ctx context.Context, appID, accountID string, args *dto.CreateChatMessageBody, invokeFrom biz_entity_app_generate.InvokeFrom, streaming bool) error {

	accountRecord, err := s.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return err
	}

	tenant, _, err := s.accountDomain.GetCurrentTenantOfAccount(ctx, accountRecord.ID)

	if err != nil {
		return err
	}

	appModel, err := s.appDomain.AppRepo.GetTenantApp(ctx, appID, tenant.ID)

	if err != nil {
		return err
	}

	chatAppGenerator := app_chat_generator.NewChatAppGenerator(s.appDomain, s.providerDomain, s.chatDomain)

	if err := chatAppGenerator.Generate(ctx, appModel, accountRecord, args, invokeFrom, true); err != nil {
		return err
	}

	return nil
}

func (s *ChatService) ListConsoleMessagesOfConversation(ctx context.Context, accountID, appID string, args *dto.ListChatMessageQuery) (*dto.ListChatMessagesResponse, error) {

	accountRecord, err := s.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, err
	}

	tenant, tenantJoin, err := s.accountDomain.GetCurrentTenantOfAccount(ctx, accountRecord.ID)

	if err != nil {
		return nil, err
	}

	if !tenantJoin.IsEditor() {
		return nil, errors.WithCode(code.ErrForbidden, "You don't have the permission for %s", tenant.Name)
	}

	app, err := s.appDomain.AppRepo.GetTenantApp(ctx, appID, tenant.ID)

	if err != nil {
		return nil, err
	}

	conversation, err := s.chatDomain.MessageRepo.GetConversationByApp(ctx, args.ConversationID, app.ID)

	if err != nil {
		return nil, err
	}

	messageRecords, count, err := s.chatDomain.MessageRepo.FindConsoleAppMessages(ctx, conversation.ID, args.Limit, args.FirstID)

	if err != nil {
		return nil, err
	}
	messageItems := make([]*dto.ListChatMessageItem, 0, 10)

	hasMore := true

	for _, mr := range messageRecords {
		annotation, err := s.chatDomain.AnnotationRepo.GetMessageAnnotation(ctx, mr.ID)

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		messageDto := assembler.ConvertToListMessageDto(mr, annotation)
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
	rets := make([]*dto.ListChatConversationItem, 0, 10)
	var sessionID string

	accountRecord, err := s.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, err
	}

	tenant, tenantJoin, err := s.accountDomain.GetCurrentTenantOfAccount(ctx, accountRecord.ID)

	if err != nil {
		return nil, err
	}

	if !tenantJoin.IsEditor() {
		return nil, errors.WithCode(code.ErrForbidden, "You don't have the permission for %s", tenant.Name)
	}

	app, err := s.appDomain.AppRepo.GetTenantApp(ctx, appID, tenant.ID)

	if err != nil {
		return nil, err
	}

	conversationRecords, count, err := s.chatDomain.MessageRepo.FindConversationsInConsole(ctx, args.Page, args.Limit, app.ID, args.Start, args.End, args.SortBy, args.Keyword, accountRecord.Timezone)

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

func (s *ChatService) DetailConversation(ctx context.Context, accountID string, cID string, appID string) (*dto.ListChatConversationItem, error) {
	var sessionID string

	accountRecord, err := s.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, err
	}

	tenant, tenantJoin, err := s.accountDomain.GetCurrentTenantOfAccount(ctx, accountRecord.ID)

	if err != nil {
		return nil, err
	}

	if !tenantJoin.IsEditor() {
		return nil, errors.WithCode(code.ErrForbidden, "You don't have the permission for %s", tenant.Name)
	}

	app, err := s.appDomain.AppRepo.GetTenantApp(ctx, appID, tenant.ID)

	if err != nil {
		return nil, err
	}

	conversationRecord, err := s.chatDomain.MessageRepo.GetConversationByApp(ctx, cID, app.ID)

	if err != nil {
		return nil, err
	}

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

	conversationJoin.MessageCount = msgCount

	if conversationRecord.OverrideModelConfigs != nil {
		conversationJoin.ModelConfig = conversationRecord.OverrideModelConfigs
		conversationJoin.ModelConfig.ModelID = conversationRecord.ModelID
		conversationJoin.ModelConfig.Provider = conversationRecord.ModelProvider
	} else {
		appConf, err := s.appDomain.AppRepo.GetAppModelConfigById(ctx, conversationRecord.AppModelConfigID, app.ID)
		if err != nil {
			return nil, err
		}
		conversationJoin.ModelConfig = biz_entity.ConvertToAppConfigBizEntity(appConf, nil)
		conversationJoin.ModelConfig.ModelID = conversationRecord.ModelID
		conversationJoin.ModelConfig.Provider = conversationRecord.ModelProvider
	}

	conversationJoin.FromAccountName = account.Name
	conversationJoin.UserFeedbackStats = dto.NewFeedBackStats()
	conversationJoin.AdminFeedbackStats = dto.NewFeedBackStats()

	if conversationRecord.FromEndUserID != "" {
		conversationJoin.FromEndUserSessionID = sessionID
	}

	return conversationJoin, nil
}
