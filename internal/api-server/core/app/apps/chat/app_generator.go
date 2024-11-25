package chat

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/google/uuid"
	"github.com/lunarianss/Luna/internal/api-server/core/app"
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config"
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config/model_config"
	"github.com/lunarianss/Luna/internal/api-server/core/app/apps"
	appEntities "github.com/lunarianss/Luna/internal/api-server/core/app/apps/entities"
	"github.com/lunarianss/Luna/internal/api-server/core/app/task_pipeline"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/provider"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
	"github.com/lunarianss/Luna/internal/api-server/entities/message"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
	"github.com/lunarianss/Luna/pkg/log"
)

type ChatAppGenerator struct {
	AppDomain      *appDomain.AppDomain
	ProviderDomain *domain.ModelProviderDomain
	chatDomain     *chatDomain.ChatDomain
}

func NewChatAppGenerator(appDomain *appDomain.AppDomain, providerDomain *domain.ModelProviderDomain, chatDomain *chatDomain.ChatDomain) *ChatAppGenerator {

	return &ChatAppGenerator{
		AppDomain:      appDomain,
		ProviderDomain: providerDomain,
		chatDomain:     chatDomain,
	}

}

func (g *ChatAppGenerator) getAppModelConfig(ctx context.Context, appModel *model.App, conversation *model.Conversation) (*model.AppModelConfig, error) {
	if conversation == nil {
		if appModel.AppModelConfigID == "" {
			return nil, errors.WithCode(code.ErrAppNotFoundRelatedConfig, fmt.Sprintf("app %s not found related config", appModel.Name))
		}

		return g.AppDomain.AppRepo.GetAppModelConfigById(ctx, appModel.AppModelConfigID)
	} else {
		return nil, errors.New("todo")
	}
}

func (g *ChatAppGenerator) Generate(c context.Context, appModel *model.App, user model.BaseAccount, args *dto.CreateChatMessageBody, invokeFrom appEntities.InvokeForm, stream bool) error {

	var (
		conversationRecord     *model.Conversation
		messageRecord          *model.Message
		extras                 map[string]interface{}
		overrideModelConfigMap map[string]interface{}
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

	appModelConfig, err := g.getAppModelConfig(c, appModel, conversationRecord)

	if err != nil {
		return err
	}

	modelConfigManager := NewChatAppConfigManager(g.ProviderDomain)
	if args.ModelConfig != nil {
		if invokeFrom != appEntities.DEBUGGER {
			return errors.WithCode(code.ErrOnlyOverrideConfigInDebugger, fmt.Sprintf("mode %s is not debugger, so it cannot override", invokeFrom))
		}

		overrideModelConfigMap, err = modelConfigManager.ConfigValidate(c, appModel.TenantID, args.ModelConfig)
		if err != nil {
			return err
		}

		overrideModelConfigMap["retriever_resource"] = map[string]any{
			"enabled": true,
		}
	}

	appConfig, err := modelConfigManager.getAppConfig(c, appModel, appModelConfig, conversationRecord, overrideModelConfigMap)

	if err != nil {
		return err
	}

	var conversationID string

	if conversationRecord != nil {
		if conversationRecord.ID != "" {
			conversationID = conversationRecord.ID
		}
	}

	modelConverter := model_config.NewModelConfigConverter(g.ProviderDomain)
	modelConf, err := modelConverter.Convert(c, appConfig.EasyUIBasedAppConfig, true)

	if err != nil {
		return err
	}

	applicationGenerateEntity := &app.ChatAppGenerateEntity{
		ConversationID:  conversationID,
		ParentMessageID: args.ParentMessageId,
		EasyUIBasedAppGenerateEntity: &app.EasyUIBasedAppGenerateEntity{
			AppConfig: appConfig.EasyUIBasedAppConfig,
			ModelConf: modelConf,
			AppGenerateEntity: &app.AppGenerateEntity{
				TaskID:     uuid.NewString(),
				AppConfig:  appConfig.AppConfig,
				Stream:     stream,
				Inputs:     inputs,
				UserID:     user.GetAccountID(),
				InvokeFrom: app.InvokeFrom(invokeFrom),
				Extras:     extras,
			},
			Query: query,
		},
	}

	conversationRecord, messageRecord, err = g.InitGenerateRecords(c, applicationGenerateEntity, conversationRecord)

	if err != nil {
		return err
	}

	queueManager, streamResultChunkQueue, streamFinalChunkQueue := model_runtime.NewStreamGenerateQueue(
		uuid.NewString(),
		applicationGenerateEntity.UserID,
		conversationRecord.ID,
		messageRecord.ID,
		model.AppMode("chat"),
		invokeFrom)

	go g.generateGoRoutine(c, applicationGenerateEntity, conversationRecord.ID, messageRecord.ID, queueManager)

	go g.ListenQueue(queueManager)

	task_pipeline.NewChatAppTaskPipeline(applicationGenerateEntity, streamResultChunkQueue, streamFinalChunkQueue, g.chatDomain.MessageRepo, messageRecord, conversationRecord.ID).Process(c, true)

	return nil
}

func (g *ChatAppGenerator) ListenQueue(queueManager *model_runtime.StreamGenerateQueue) {
	queueManager.Listen()
}

func (g *ChatAppGenerator) generateGoRoutine(ctx context.Context, applicationGenerateEntity *app.ChatAppGenerateEntity, conversationID string, messageID string, queueManager *model_runtime.StreamGenerateQueue) {

	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovered from generateGoRoutine panic: %+v", r)
			log.Errorf("Stack trace: %s", debug.Stack())
		}
	}()

	appRunner := &apps.AppRunner{
		AppDomain: g.AppDomain,
	}

	message, err := g.AppDomain.AppRepo.GetMessageByID(ctx, messageID)

	if err != nil {
		queueManager.PushErr(err)
		return
	}

	conversation, err := g.AppDomain.AppRepo.GetConversationByID(ctx, conversationID)

	if err != nil {
		queueManager.PushErr(err)
		return
	}

	appRunner.Run(ctx, applicationGenerateEntity, message, conversation, queueManager)
}

func (g *ChatAppGenerator) InitGenerateRecords(ctx context.Context, chatAppGenerateEntity *app.ChatAppGenerateEntity, conversation *model.Conversation) (*model.Conversation, *model.Message, error) {

	appConfig := chatAppGenerateEntity.AppConfig

	var (
		fromSource          string
		accountID           string
		conversationRecord  *model.Conversation
		messageRecord       *model.Message
		endUserID           string
		appModelConfigID    string
		modelProvider       string
		modelID             string
		overrideModelConfig map[string]interface{}
	)

	if chatAppGenerateEntity.InvokeFrom == app.WebApp || chatAppGenerateEntity.InvokeFrom == app.ServiceAPI {
		fromSource = "api"
		endUserID = chatAppGenerateEntity.UserID
	} else {
		fromSource = "console"
		accountID = chatAppGenerateEntity.UserID
	}

	appModelConfigID = appConfig.AppModelConfigID
	modelProvider = chatAppGenerateEntity.ModelConf.Provider
	modelID = chatAppGenerateEntity.ModelConf.Model

	if appConfig.AppModelConfigFrom == app_config.Args && (appConfig.AppMode == string(model.CHAT) || appConfig.AppMode == string(model.AGENT_CHAT) || appConfig.AppMode == string(model.COMPLETION)) {
		overrideModelConfig = appConfig.AppModelConfigDict
	}

	if conversation == nil {
		var err error
		conversationRecord = &model.Conversation{
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
		conversationRecord, err = g.AppDomain.AppRepo.CreateConversation(ctx, conversationRecord)

		if err != nil {
			return nil, nil, err
		}
	}

	message := &model.Message{
		AppID:                   appConfig.AppID,
		ModelProvider:           modelProvider,
		ModelID:                 modelID,
		OverrideModelConfigs:    overrideModelConfig,
		ConversationID:          conversationRecord.ID,
		Inputs:                  chatAppGenerateEntity.Inputs,
		Query:                   chatAppGenerateEntity.Query,
		Message:                 make([]*message.PromptMessage, 0),
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

	messageRecord, err := g.AppDomain.AppRepo.CreateMessage(ctx, message)

	if err != nil {
		return nil, nil, err
	}

	return conversationRecord, messageRecord, nil
}
