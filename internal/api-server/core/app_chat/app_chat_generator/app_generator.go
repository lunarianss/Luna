// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package app_chat_generator

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/google/uuid"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/internal/api-server/core/app_chat/app_chat_runner"
	"github.com/lunarianss/Luna/internal/api-server/core/app_chat/task_pipeline"
	"github.com/lunarianss/Luna/internal/api-server/core/app_config/app_config"
	"github.com/lunarianss/Luna/internal/api-server/core/app_config/app_model_config"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	biz_entity_app_config "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/provider_app_config"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	biz_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	po_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/common/repository"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	biz_entity_app_generate "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_generate"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"

	"github.com/lunarianss/Luna/internal/infrastructure/code"
)

type IChatAppGenerator interface {
	Generate(c context.Context, appModel *po_entity.App, user repository.BaseAccount, args *dto.CreateChatMessageBody, invokeFrom biz_entity_app_generate.InvokeFrom, stream bool) error

	GenerateNonStream(c context.Context, appModel *po_entity.App, user repository.BaseAccount, args *dto.CreateChatMessageBody, invokeFrom biz_entity_app_generate.InvokeFrom, stream bool) (*biz_entity_chat.LLMResult, error)
}

type ChatAppGenerator struct {
	AppDomain      *appDomain.AppDomain
	ProviderDomain *domain_service.ProviderDomain
	chatDomain     *chatDomain.ChatDomain
}

func NewChatAppGenerator(appDomain *appDomain.AppDomain, providerDomain *domain_service.ProviderDomain, chatDomain *chatDomain.ChatDomain) *ChatAppGenerator {

	return &ChatAppGenerator{
		AppDomain:      appDomain,
		ProviderDomain: providerDomain,
		chatDomain:     chatDomain,
	}

}

func (g *ChatAppGenerator) getAppModelConfig(ctx context.Context, appModel *po_entity.App, conversation *po_entity_chat.Conversation) (*po_entity.AppModelConfig, error) {
	if conversation == nil {
		if appModel.AppModelConfigID == "" {
			return nil, errors.WithCode(code.ErrAppNotFoundRelatedConfig, fmt.Sprintf("app %s not found related config", appModel.Name))
		}
		return g.AppDomain.AppRepo.GetAppModelConfigById(ctx, appModel.AppModelConfigID, appModel.ID)
	} else {
		if appModel.AppModelConfigID == "" {
			return nil, errors.WithCode(code.ErrAppNotFoundRelatedConfig, fmt.Sprintf("conversation %s not found related config", appModel.Name))
		}
		return g.AppDomain.AppRepo.GetAppModelConfigById(ctx, conversation.AppModelConfigID, appModel.ID)
	}
}

func (g *ChatAppGenerator) baseGenerate(c context.Context, appModel *po_entity.App, user repository.BaseAccount, args *dto.CreateChatMessageBody, invokeFrom biz_entity_app_generate.InvokeFrom, stream bool) (*biz_entity_app_generate.ChatAppGenerateEntity, *po_entity_chat.Conversation, *po_entity_chat.Message, error) {

	var (
		conversationRecord     *po_entity_chat.Conversation
		messageRecord          *po_entity_chat.Message
		extras                 map[string]interface{}
		overrideModelConfigMap *dto.AppModelConfigDto
		conversationID         string
		err                    error
	)

	query := args.Query
	inputs := args.Inputs
	// role := model.AccountCreatedByRole
	extras = make(map[string]interface{})

	if !args.AutoGenerateConversationName {
		extras["auto_generate_conversation_name"] = true
	} else {
		extras["auto_generate_conversation_name"] = args.AutoGenerateConversationName
	}

	if args.ConversationID != "" {
		conversationRecord, err = g.chatDomain.MessageRepo.GetConversationByUser(c, appModel.ID, args.ConversationID, user)

		if err != nil {
			return nil, nil, nil, err
		}
	}

	appModelConfig, err := g.getAppModelConfig(c, appModel, conversationRecord)

	if err != nil {
		return nil, nil, nil, err
	}

	modelConfigManager := app_config.NewChatAppConfigManager(g.ProviderDomain)
	if args.ModelConfig.AppID != "" {
		if invokeFrom != biz_entity_app_generate.Debugger {
			return nil, nil, nil, errors.WithCode(code.ErrOnlyOverrideConfigInDebugger, fmt.Sprintf("mode %s is not debugger, so it cannot override", invokeFrom))
		}

		overrideModelConfigMap, err = modelConfigManager.ConfigValidate(c, appModel.TenantID, &args.ModelConfig)
		if err != nil {
			return nil, nil, nil, err
		}

		overrideModelConfigMap.RetrieverResource.Enabled = true
	}

	appConfig, err := modelConfigManager.GetAppConfig(c, appModel, appModelConfig, conversationRecord, overrideModelConfigMap)

	if err != nil {
		return nil, nil, nil, err
	}

	if conversationRecord != nil {
		conversationID = conversationRecord.ID
	}

	modelConverter := app_model_config.NewModelConfigConverter(g.ProviderDomain)
	modelConf, err := modelConverter.Convert(c, appConfig.EasyUIBasedAppConfig, true)

	if err != nil {
		return nil, nil, nil, err
	}

	applicationGenerateEntity := &biz_entity_app_generate.ChatAppGenerateEntity{
		ConversationID:  conversationID,
		ParentMessageID: args.ParentMessageId,
		EasyUIBasedAppGenerateEntity: &biz_entity_app_generate.EasyUIBasedAppGenerateEntity{
			AppConfig: appConfig.EasyUIBasedAppConfig,
			ModelConf: modelConf,
			AppGenerateEntity: &biz_entity_app_generate.AppGenerateEntity{
				TaskID:     uuid.NewString(),
				AppConfig:  appConfig.AppConfig,
				Stream:     stream,
				Inputs:     inputs,
				UserID:     user.GetAccountID(),
				InvokeFrom: biz_entity_app_generate.InvokeFrom(invokeFrom),
				Extras:     extras,
			},
			Query: query,
		},
	}

	conversationRecord, messageRecord, err = g.InitGenerateRecords(c, applicationGenerateEntity, conversationRecord)

	if err != nil {
		return nil, nil, nil, err
	}
	return applicationGenerateEntity, conversationRecord, messageRecord, nil
}

func (g *ChatAppGenerator) GenerateNonStream(c context.Context, appModel *po_entity.App, user repository.BaseAccount, args *dto.CreateChatMessageBody, invokeFrom biz_entity_app_generate.InvokeFrom, stream bool) (*biz_entity_chat.LLMResult, error) {

	applicationGenerateEntity, conversationRecord, messageRecord, err := g.baseGenerate(c, appModel, user, args, invokeFrom, stream)

	if err != nil {
		return nil, err
	}

	return g.generateNonStream(c, applicationGenerateEntity, conversationRecord.ID, messageRecord.ID)
}

func (g *ChatAppGenerator) Generate(c context.Context, appModel *po_entity.App, user repository.BaseAccount, args *dto.CreateChatMessageBody, invokeFrom biz_entity_app_generate.InvokeFrom, stream bool) error {

	applicationGenerateEntity, conversationRecord, messageRecord, err := g.baseGenerate(c, appModel, user, args, invokeFrom, stream)

	if err != nil {
		return err
	}

	queueManager, streamResultChunkQueue, streamFinalChunkQueue := biz_entity_chat.NewStreamGenerateQueue(
		uuid.NewString(),
		applicationGenerateEntity.UserID,
		conversationRecord.ID,
		messageRecord.ID,
		po_entity.AppMode("chat"),
		string(invokeFrom))

	go g.generateGoRoutine(c, applicationGenerateEntity, conversationRecord.ID, messageRecord.ID, queueManager)

	go g.ListenQueue(queueManager)

	task_pipeline.NewChatAppTaskPipeline(applicationGenerateEntity, streamResultChunkQueue, streamFinalChunkQueue, g.chatDomain.MessageRepo, messageRecord).Process(c)

	return nil
}

func (g *ChatAppGenerator) ListenQueue(queueManager *biz_entity_chat.StreamGenerateQueue) {
	queueManager.Listen()
}

func (g *ChatAppGenerator) generateNonStream(ctx context.Context, applicationGenerateEntity *biz_entity_app_generate.ChatAppGenerateEntity, conversationID string, messageID string) (*biz_entity_chat.LLMResult, error) {

	appRunner := app_chat_runner.NewAppChatRunner(app_chat_runner.NewAppBaseChatRunner(), g.AppDomain, g.chatDomain)

	message, err := g.chatDomain.MessageRepo.GetMessageByID(ctx, messageID)

	if err != nil {

		return nil, err
	}

	conversation, err := g.chatDomain.MessageRepo.GetConversationByID(ctx, conversationID)

	if err != nil {
		return nil, err
	}

	return appRunner.RunNonStream(ctx, applicationGenerateEntity, message, conversation)
}

func (g *ChatAppGenerator) generateGoRoutine(ctx context.Context, applicationGenerateEntity *biz_entity_app_generate.ChatAppGenerateEntity, conversationID string, messageID string, queueManager *biz_entity_chat.StreamGenerateQueue) {

	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovered from generateGoRoutine panic: %+v", r)
			log.Errorf("Stack trace: %s", debug.Stack())
		}
	}()

	appRunner := app_chat_runner.NewAppChatRunner(app_chat_runner.NewAppBaseChatRunner(), g.AppDomain, g.chatDomain)

	message, err := g.chatDomain.MessageRepo.GetMessageByID(ctx, messageID)

	if err != nil {
		queueManager.PushErr(err)
		return
	}

	conversation, err := g.chatDomain.MessageRepo.GetConversationByID(ctx, conversationID)

	if err != nil {
		queueManager.PushErr(err)
		return
	}

	appRunner.Run(ctx, applicationGenerateEntity, message, conversation, queueManager)
}

func (g *ChatAppGenerator) InitGenerateRecords(ctx context.Context, chatAppGenerateEntity *biz_entity_app_generate.ChatAppGenerateEntity, conversation *po_entity_chat.Conversation) (*po_entity_chat.Conversation, *po_entity_chat.Message, error) {

	appConfig := chatAppGenerateEntity.AppConfig

	var (
		fromSource          string
		accountID           string
		conversationRecord  *po_entity_chat.Conversation
		messageRecord       *po_entity_chat.Message
		endUserID           string
		appModelConfigID    string
		modelProvider       string
		modelID             string
		overrideModelConfig *biz_entity_app_config.AppModelConfig
	)

	if chatAppGenerateEntity.InvokeFrom == biz_entity_app_generate.WebApp || chatAppGenerateEntity.InvokeFrom == biz_entity_app_generate.ServiceAPI {
		fromSource = "api"
		endUserID = chatAppGenerateEntity.UserID
	} else {
		fromSource = "console"
		accountID = chatAppGenerateEntity.UserID
	}

	appModelConfigID = appConfig.AppModelConfigID
	modelProvider = chatAppGenerateEntity.ModelConf.Provider
	modelID = chatAppGenerateEntity.ModelConf.Model

	if appConfig.AppModelConfigFrom == biz_entity_app_config.Args && (appConfig.AppMode == string(po_entity.CHAT) || appConfig.AppMode == string(po_entity.AGENT_CHAT) || appConfig.AppMode == string(po_entity.COMPLETION)) {
		overrideModelConfig = appConfig.AppModelConfig
	}

	if conversation == nil {
		var err error
		conversationRecord = &po_entity_chat.Conversation{
			AppID:                   appConfig.AppID,
			AppModelConfigID:        appModelConfigID,
			ModelProvider:           modelProvider,
			ModelID:                 modelID,
			OverrideModelConfigs:    overrideModelConfig,
			Mode:                    appConfig.AppMode,
			Name:                    "New Conversation",
			Inputs:                  chatAppGenerateEntity.Inputs,
			Introduction:            "",
			SystemInstruction:       "",
			SystemInstructionTokens: 0,
			Status:                  "normal",
			InvokeFrom:              string(chatAppGenerateEntity.InvokeFrom),
			FromSource:              fromSource,
			FromEndUserID:           endUserID,
			FromAccountID:           accountID,
		}
		conversationRecord, err = g.chatDomain.MessageRepo.CreateConversation(ctx, conversationRecord)

		if err != nil {
			return nil, nil, err
		}
		chatAppGenerateEntity.ConversationID = conversationRecord.ID
	} else {
		conversationRecord = conversation
		conversation.UpdatedAt = time.Now().UTC().Unix()
		if err := g.chatDomain.MessageRepo.UpdateConversationUpdateAt(ctx, conversation.AppID, conversation); err != nil {
			return nil, nil, err
		}
	}

	message := &po_entity_chat.Message{
		AppID:                   appConfig.AppID,
		ModelProvider:           modelProvider,
		ModelID:                 modelID,
		OverrideModelConfigs:    overrideModelConfig,
		ConversationID:          conversationRecord.ID,
		Inputs:                  chatAppGenerateEntity.Inputs,
		Query:                   chatAppGenerateEntity.Query,
		Message:                 make([]*po_entity_chat.PromptMessage, 0),
		MessageTokens:           0,
		MessageUnitPrice:        0,
		MessagePriceUnit:        0,
		Answer:                  "",
		AnswerTokens:            0,
		AnswerUnitPrice:         0,
		AnswerPriceUnit:         0,
		ParentMessageID:         chatAppGenerateEntity.ParentMessageID,
		ProviderResponseLatency: 0,
		TotalPrice:              0,
		Currency:                "USD",
		InvokeFrom:              string(chatAppGenerateEntity.InvokeFrom),
		FromSource:              fromSource,
		FromEndUserID:           endUserID,
		FromAccountID:           accountID,
	}

	messageRecord, err := g.chatDomain.MessageRepo.CreateMessage(ctx, message)

	if err != nil {
		return nil, nil, err
	}

	return conversationRecord, messageRecord, nil
}
