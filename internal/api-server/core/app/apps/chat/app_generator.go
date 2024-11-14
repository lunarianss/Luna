package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lunarianss/Luna/internal/api-server/core/app"
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config"
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config/entities"
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config/model_config"
	"github.com/lunarianss/Luna/internal/api-server/core/app/apps"
	appEntities "github.com/lunarianss/Luna/internal/api-server/core/app/apps/entities"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/provider"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
	"github.com/lunarianss/Luna/pkg/log"
)

type ChatAppGenerator struct {
	AppDomain      *appDomain.AppDomain
	ProviderDomain *domain.ModelProviderDomain
}

func NewChatAppGenerator(appDomain *appDomain.AppDomain, providerDomain *domain.ModelProviderDomain) *ChatAppGenerator {

	return &ChatAppGenerator{
		AppDomain:      appDomain,
		ProviderDomain: providerDomain,
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

func (g *ChatAppGenerator) Generate(c context.Context, appModel *model.App, user interface{}, args *dto.CreateChatMessageBody, invokeFrom appEntities.InvokeForm, stream bool) error {

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
	if args.AutoGenerateConversationName == nil {
		extras["auto_generate_conversation_name"] = true
	} else {
		extras["auto_generate_conversation_name"] = *args.AutoGenerateConversationName
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

		overrideModelConfigMap, err := modelConfigManager.ConfigValidate(c, appModel.TenantID, args.ModelConfig)
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

	var conversationID *string

	if conversationRecord != nil {
		if conversationRecord.ID != "" {
			conversationID = &conversationRecord.ID
		}
	}

	modelConverter := model_config.NewModelConfigConverter(g.ProviderDomain)
	modelConf, err := modelConverter.Convert(c, appConfig.EasyUIBasedAppConfig, true)

	if err != nil {
		return err
	}

	applicationGenerateEntity := &app.ChatAppGenerateEntity{
		ConversationID:  conversationID,
		ParentMessageID: &args.ParentMessageId,
		EasyUIBasedAppGenerateEntity: &app.EasyUIBasedAppGenerateEntity{
			AppConfig: appConfig.EasyUIBasedAppConfig,
			ModelConf: modelConf,
			AppGenerateEntity: &app.AppGenerateEntity{
				TaskID:     uuid.NewString(),
				AppConfig:  appConfig.AppConfig,
				Stream:     stream,
				Inputs:     inputs,
				UserID:     "",
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

	g.handleMessageQueueEvent(c, streamResultChunkQueue, streamFinalChunkQueue)
	return nil
}

func (g *ChatAppGenerator) ListenQueue(queueManager *model_runtime.StreamGenerateQueue) {
	queueManager.Listen()
}

func (g *ChatAppGenerator) handleMessageQueueEvent(c context.Context, streamResultChunkQueue chan *entities.MessageQueueMessage, streamFinalChunkQueue chan *entities.MessageQueueMessage) {
	// 确保 Gin 使用 HTTP 流式传输
	c.(*gin.Context).Writer.Header().Set("Content-Type", "text/event-stream")
	c.(*gin.Context).Writer.Header().Set("Cache-Control", "no-cache")
	c.(*gin.Context).Writer.Header().Set("Connection", "keep-alive")

	// 确保 c.Writer 实现了 http.Flusher 接口
	flusher, ok := c.(*gin.Context).Writer.(http.Flusher)
	if !ok {
		c.(*gin.Context).String(http.StatusInternalServerError, "Streaming unsupported!")
		return
	}

	for v := range streamResultChunkQueue {
		if cm, ok := v.Event.(*entities.QueueLLMChunkEvent); ok {
			// 将事件格式化为 SSE 格式发送给客户端
			fmt.Fprintf(c.(*gin.Context).Writer, "data: {answer: %s}\n\n", cm.Chunk.Delta.Message.Content)
			flusher.Flush() // 确保数据立即发送到客户端
		}
	}

	for v := range streamFinalChunkQueue {
		if mc, ok := v.Event.(*entities.QueueLLMChunkEvent); ok {
			chunkByte, _ := json.Marshal(mc.Chunk)
			log.Infof("Event type: %s, Answer: %s", mc.Event, string(chunkByte))
			// 将事件格式化为 SSE 格式发送给客户端
			fmt.Fprintf(c.(*gin.Context).Writer, "data: %s\n\n", mc.Chunk.Delta.Message.Content)
			flusher.Flush() // 确保数据立即发送到客户端
		} else if mc, ok := v.Event.(*entities.QueueMessageEndEvent); ok {
			log.Infof("Event type: %s, End LLM Result %+v", mc.Event, mc.LLMResult)
		} else if mc, ok := v.Event.(*entities.QueueErrorEvent); ok {
			log.Errorf("Event type: %s, Err: %s", mc.Event, mc.Err.Error())
		}
	}
}

func (g *ChatAppGenerator) generateGoRoutine(ctx context.Context, applicationGenerateEntity *app.ChatAppGenerateEntity, conversationID string, messageID string, queueManager *model_runtime.StreamGenerateQueue) {

	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovered from generateGoRoutine panic: %v", r)
		}
	}()

	appRunner := &apps.AppRunner{}

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
		Message:                 make([]map[string]interface{}, 0),
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
