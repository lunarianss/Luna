package app_agent_chat_generator

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/internal/api-server/config"
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_agent_chat_runner"
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_chat_runner"
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config/app_agent_config"
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config/app_model_config"
	"github.com/lunarianss/Luna/internal/api-server/core/app/task_pipeline"
	agentDomain "github.com/lunarianss/Luna/internal/api-server/domain/agent/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/agent/entity/biz_entity"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	biz_entity_app_config "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/provider_app_config"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	biz_entity_agent_generator "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity/stream_agent_generator"
	biz_entity_base_stream_generator "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity/stream_base_generator"
	po_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/common/repository"
	datasetDomain "github.com/lunarianss/Luna/internal/api-server/domain/dataset/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	biz_entity_app_generate "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_generate"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/redis/go-redis/v9"
)

const UUID_NIL = "00000000-0000-0000-0000-000000000000"

type IAgentChatAppGenerator interface {
	Generate(c context.Context, appModel *po_entity.App, user repository.BaseAccount, args *dto.CreateChatMessageBody, invokeFrom biz_entity_app_generate.InvokeFrom, stream bool) error
}

type AgentChatGenerator struct {
	AppDomain      *appDomain.AppDomain
	ProviderDomain *domain_service.ProviderDomain
	chatDomain     *chatDomain.ChatDomain
	DatasetDomain  *datasetDomain.DatasetDomain
	redis          *redis.Client
	agentDomain    *agentDomain.AgentDomain
	appConfig      *biz_entity_app_config.AgentChatAppConfig
	config         *config.Config
}

func NewChatAppGenerator(appDomain *appDomain.AppDomain, providerDomain *domain_service.ProviderDomain, chatDomain *chatDomain.ChatDomain, datasetDomain *datasetDomain.DatasetDomain, redis *redis.Client, agentDomain *agentDomain.AgentDomain, config *config.Config) *AgentChatGenerator {

	return &AgentChatGenerator{
		AppDomain:      appDomain,
		ProviderDomain: providerDomain,
		chatDomain:     chatDomain,
		DatasetDomain:  datasetDomain,
		redis:          redis,
		config:         config,
		agentDomain:    agentDomain,
	}
}

func (acg *AgentChatGenerator) getAppModelConfig(ctx context.Context, appModel *po_entity.App, conversation *po_entity_chat.Conversation) (*po_entity.AppModelConfig, error) {
	if conversation == nil {
		if appModel.AppModelConfigID == "" {
			return nil, errors.WithCode(code.ErrAppNotFoundRelatedConfig, "app %s not found related config", appModel.Name)
		}
		return acg.AppDomain.AppRepo.GetAppModelConfigById(ctx, appModel.AppModelConfigID, appModel.ID)
	} else {
		if appModel.AppModelConfigID == "" {
			return nil, errors.WithCode(code.ErrAppNotFoundRelatedConfig, "conversation %s not found related config", appModel.Name)
		}
		return acg.AppDomain.AppRepo.GetAppModelConfigById(ctx, conversation.AppModelConfigID, appModel.ID)
	}
}

func (acg *AgentChatGenerator) Generate(c context.Context, appModel *po_entity.App, user repository.BaseAccount, args *dto.CreateChatMessageBody, invokeFrom biz_entity_app_generate.InvokeFrom, stream bool) error {

	if !stream {
		return errors.WithSCode(code.ErrNotStreamAgent, "")
	}

	var (
		conversationRecord     *po_entity_chat.Conversation
		messageRecord          *po_entity_chat.Message
		extras                 map[string]interface{}
		overrideModelConfigMap *dto.AppModelConfigDto
		conversationID         string
		err                    error
		parentMessageID        string
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
		conversationRecord, err = acg.chatDomain.MessageRepo.GetConversationByUser(c, appModel.ID, args.ConversationID, user)

		if err != nil {
			return err
		}
	}

	appModelConfig, err := acg.getAppModelConfig(c, appModel, conversationRecord)

	if err != nil {
		return err
	}

	modelConfigManager := app_agent_config.NewAgentChatAppConfigManager(acg.ProviderDomain)

	if args.ModelConfig.AppID != "" {
		if invokeFrom != biz_entity_app_generate.Debugger {
			return errors.WithCode(code.ErrOnlyOverrideConfigInDebugger, "mode %s is not debugger, so it cannot override", invokeFrom)
		}

		overrideModelConfigMap, err = modelConfigManager.ConfigValidate(c, appModel.TenantID, &args.ModelConfig)
		if err != nil {
			return err
		}

		overrideModelConfigMap.RetrieverResource.Enabled = true
	}

	appConfig, err := modelConfigManager.GetAppConfig(c, appModel, appModelConfig, conversationRecord, overrideModelConfigMap)

	acg.appConfig = appConfig

	if err != nil {
		return err
	}

	if conversationRecord != nil {
		conversationID = conversationRecord.ID
	}

	modelConverter := app_model_config.NewModelConfigConverter(acg.ProviderDomain)
	modelConf, err := modelConverter.Convert(c, appConfig.EasyUIBasedAppConfig, true)

	if err != nil {
		return err
	}

	if invokeFrom != biz_entity_app_generate.ServiceAPI {
		parentMessageID = args.ParentMessageId
	} else {
		parentMessageID = UUID_NIL
	}

	applicationGenerateEntity := &biz_entity_app_generate.AgentChatAppGenerateEntity{
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
		AgentEntity:     appConfig.AgentEntity,
		ConversationID:  conversationID,
		ParentMessageID: parentMessageID,
	}

	conversationRecord, messageRecord, err = acg.InitGenerateRecords(c, applicationGenerateEntity, conversationRecord)

	if err != nil {
		return err
	}

	queueManager := biz_entity_agent_generator.NewAgentStreamGenerateQueue(
		uuid.NewString(),
		applicationGenerateEntity.EasyUIBasedAppGenerateEntity.UserID,
		conversationRecord.ID,
		messageRecord.ID,
		po_entity.AppMode("agent-chat"),
		string(invokeFrom))

	taskScheduler := app_agent_chat_runner.NewAgentChatAppTaskScheduler(applicationGenerateEntity, acg.chatDomain.MessageRepo, messageRecord, acg.chatDomain.AnnotationRepo, nil)

	flusher := task_pipeline.NewAgentChatFlusher(applicationGenerateEntity, acg.agentDomain.AgentRepo, messageRecord)

	flusher.InitFlusher(c)

	go acg.ListenQueue(queueManager)

	acg.generateGoRoutine(c, applicationGenerateEntity, conversationRecord.ID, messageRecord.ID, queueManager, taskScheduler, flusher)

	return nil
}

func (g *AgentChatGenerator) generateGoRoutine(ctx context.Context, applicationGenerateEntity *biz_entity_app_generate.AgentChatAppGenerateEntity, conversationID string, messageID string, queueManager biz_entity_base_stream_generator.IStreamGenerateQueue, taskPipeline app_agent_chat_runner.IAgentChatAppTaskScheduler, flusher biz_entity.AgentFlusher) {

	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovered from generateGoRoutine panic: %+v", r)
			log.Errorf("Stack trace: %s", debug.Stack())
		}
	}()

	lunaConfig, err := config.GetLunaRuntimeConfig()

	if err != nil {
		queueManager.PushErr(err)
		return
	}

	appRunner := app_agent_chat_runner.NewAppAgentChatRunner(app_agent_chat_runner.NewAppBaseAgentChatRunner(app_chat_runner.NewAppBaseChatRunner()), g.AppDomain, g.chatDomain, g.ProviderDomain, g.DatasetDomain, g.redis, g.agentDomain, applicationGenerateEntity.AppConfig.TenantID, applicationGenerateEntity.UserID, g.appConfig, applicationGenerateEntity, lunaConfig.SystemOptions.SecretKey, lunaConfig.SystemOptions.FileBaseUrl, lunaConfig.MinioOptions.Bucket)

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

	appRunner.Run(ctx, applicationGenerateEntity, message, conversation, queueManager, taskPipeline, flusher, g.appConfig)
}

func (g *AgentChatGenerator) ListenQueue(queueManager biz_entity_base_stream_generator.IStreamGenerateQueue) {
	queueManager.Listen()
}

func (g *AgentChatGenerator) InitGenerateRecords(ctx context.Context, chatAppGenerateEntity *biz_entity_app_generate.AgentChatAppGenerateEntity, conversation *po_entity_chat.Conversation) (*po_entity_chat.Conversation, *po_entity_chat.Message, error) {

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

	if chatAppGenerateEntity.EasyUIBasedAppGenerateEntity.InvokeFrom == biz_entity_app_generate.WebApp || chatAppGenerateEntity.EasyUIBasedAppGenerateEntity.InvokeFrom == biz_entity_app_generate.ServiceAPI {
		fromSource = "api"
		endUserID = chatAppGenerateEntity.EasyUIBasedAppGenerateEntity.UserID
	} else {
		fromSource = "console"
		accountID = chatAppGenerateEntity.EasyUIBasedAppGenerateEntity.UserID
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
			Inputs:                  chatAppGenerateEntity.EasyUIBasedAppGenerateEntity.Inputs,
			Introduction:            "",
			SystemInstruction:       "",
			SystemInstructionTokens: 0,
			Status:                  "normal",
			InvokeFrom:              string(chatAppGenerateEntity.EasyUIBasedAppGenerateEntity.InvokeFrom),
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
		Inputs:                  chatAppGenerateEntity.EasyUIBasedAppGenerateEntity.Inputs,
		Query:                   chatAppGenerateEntity.EasyUIBasedAppGenerateEntity.Query,
		Message:                 make([]any, 0),
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
		InvokeFrom:              string(chatAppGenerateEntity.EasyUIBasedAppGenerateEntity.InvokeFrom),
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
