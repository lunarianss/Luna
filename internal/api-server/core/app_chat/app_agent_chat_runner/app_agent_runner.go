// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package app_agent_chat_runner

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/core/app_chat/token_buffer_memory"
	"github.com/lunarianss/Luna/internal/api-server/core/app_feature"
	"github.com/lunarianss/Luna/internal/infrastructure/util"

	agentDomain "github.com/lunarianss/Luna/internal/api-server/domain/agent/domain_service"
	biz_entity_agent "github.com/lunarianss/Luna/internal/api-server/domain/agent/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	biz_entity_app_config "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/provider_app_config"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	po_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"

	"github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_registry"
	datasetDomain "github.com/lunarianss/Luna/internal/api-server/domain/dataset/domain_service"
	providerDomain "github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	biz_entity_app_generate "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_generate"
	"github.com/redis/go-redis/v9"
)

type appAgentChatRunner struct {
	*AppBaseAgentChatRunner
	AppDomain                   *domain_service.AppDomain
	agentDomain                 *agentDomain.AgentDomain
	ChatDomain                  *chatDomain.ChatDomain
	ProviderDomain              *providerDomain.ProviderDomain
	DatasetDomain               *datasetDomain.DatasetDomain
	redis                       *redis.Client
	tenantID                    string
	userID                      string
	appConfig                   *biz_entity_app_config.AgentChatAppConfig
	application_generate_entity *biz_entity_app_generate.AgentChatAppGenerateEntity
}

func NewAppAgentChatRunner(appBaseChatRunner *AppBaseAgentChatRunner, appDomain *domain_service.AppDomain, chatDomain *chatDomain.ChatDomain, providerDomain *providerDomain.ProviderDomain, datasetDomain *datasetDomain.DatasetDomain, redis *redis.Client, agentDomain *agentDomain.AgentDomain, tenantID string, userID string, appConfig *biz_entity_app_config.AgentChatAppConfig, application_generate_entity *biz_entity_app_generate.AgentChatAppGenerateEntity) *appAgentChatRunner {
	return &appAgentChatRunner{
		AppBaseAgentChatRunner:      appBaseChatRunner,
		AppDomain:                   appDomain,
		ChatDomain:                  chatDomain,
		ProviderDomain:              providerDomain,
		DatasetDomain:               datasetDomain,
		redis:                       redis,
		agentDomain:                 agentDomain,
		tenantID:                    tenantID,
		userID:                      userID,
		appConfig:                   appConfig,
		application_generate_entity: application_generate_entity,
	}
}

func (r *appAgentChatRunner) baseRun(ctx context.Context, applicationGenerateEntity *biz_entity_app_generate.AgentChatAppGenerateEntity, conversation *po_entity_chat.Conversation) (model_registry.IModelRegistryCall, []*po_entity_chat.PromptMessage, []string, *po_entity.App, error) {

	var (
		memory token_buffer_memory.ITokenBufferMemory
	)

	appRecord, err := r.AppDomain.AppRepo.GetAppByID(ctx, applicationGenerateEntity.AppConfig.AppID)

	if err != nil {
		return nil, nil, nil, nil, err
	}

	credentials, err := applicationGenerateEntity.ModelConf.ProviderModelBundle.Configuration.GetCurrentCredentials(applicationGenerateEntity.ModelConf.ProviderModelBundle.ModelTypeInstance.ModelType, applicationGenerateEntity.AppConfig.Model.Model)

	if err != nil {
		return nil, nil, nil, nil, err
	}

	if applicationGenerateEntity.ConversationID != "" {
		modelCaller := model_registry.NewModelRegisterCaller(applicationGenerateEntity.AppConfig.Model.Model, string(applicationGenerateEntity.ModelConf.ProviderModelBundle.ModelTypeInstance.ModelType), applicationGenerateEntity.ModelConf.ProviderModelBundle.Configuration.Provider.Provider, credentials, applicationGenerateEntity.ModelConf.ProviderModelBundle.ModelTypeInstance)

		memory = token_buffer_memory.NewTokenBufferMemory(conversation, modelCaller, r.ChatDomain)
	}

	promptMessages, stop, err := r.OrganizePromptMessage(ctx, appRecord, applicationGenerateEntity.ModelConf, applicationGenerateEntity.AppConfig.PromptTemplate, applicationGenerateEntity.EasyUIBasedAppGenerateEntity.Inputs, nil, applicationGenerateEntity.Query, "", memory)

	if err != nil {
		return nil, nil, nil, nil, err
	}

	modelInstance := model_registry.NewModelRegisterCaller(applicationGenerateEntity.AppConfig.Model.Model, string(applicationGenerateEntity.ModelConf.ProviderModelBundle.ModelTypeInstance.ModelType), applicationGenerateEntity.ModelConf.ProviderModelBundle.Configuration.Provider.Provider, credentials, applicationGenerateEntity.ModelConf.ProviderModelBundle.ModelTypeInstance)

	return modelInstance, promptMessages, stop, appRecord, nil
}

func (r *appAgentChatRunner) Run(ctx context.Context, applicationGenerateEntity *biz_entity_app_generate.AgentChatAppGenerateEntity, message *po_entity_chat.Message, conversation *po_entity_chat.Conversation, queueManager *biz_entity.StreamGenerateQueue, taskScheduler IAgentChatAppTaskScheduler, flusher biz_entity_agent.AgentFlusher) {

	modelCaller, promptMessages, stop, app, err := r.baseRun(ctx, applicationGenerateEntity, conversation)

	if err != nil {
		queueManager.PushErr(err)
		return
	}

	if applicationGenerateEntity.Query != "" {
		annotation, err := r.QueryAppAnnotationToReply(ctx, app, message, applicationGenerateEntity.Query, applicationGenerateEntity.EasyUIBasedAppGenerateEntity.UserID, string(applicationGenerateEntity.EasyUIBasedAppGenerateEntity.InvokeFrom))

		if err != nil {
			queueManager.PushErr(err)
			return
		}

		if annotation != nil {
			queueEvent := biz_entity.NewAppQueueEvent(biz_entity.AnnotationReply)
			queueManager.Push(&biz_entity.QueueAnnotationReplyEvent{
				AppQueueEvent:       queueEvent,
				MessageAnnotationID: annotation.ID,
			})

			r.DirectOutStream(applicationGenerateEntity, message, conversation, queueManager, annotation.Content, promptMessages)
			return
		}
	}

	if applicationGenerateEntity.Strategy == biz_entity_app_config.FUNCTION_CALLING {
		toolRuntimeMap, promptToolMessage, err := r.InitPromptTools()

		if err != nil {
			queueManager.PushErr(err)
			return
		}

		go modelCaller.InvokeLLM(ctx, util.ConvertToInterfaceSlice(promptMessages, func(pm *po_entity_chat.PromptMessage) po_entity_chat.IPromptMessage {
			return pm
		}), queueManager, applicationGenerateEntity.ModelConf.Parameters, promptToolMessage, stop, applicationGenerateEntity.UserID, nil)

		agentRunner := NewFunctionCallAgentRunner(app.TenantID, applicationGenerateEntity, conversation, r.agentDomain, queueManager, flusher, promptToolMessage, util.ConvertToInterfaceSlice(promptMessages, func(pm *po_entity_chat.PromptMessage) po_entity_chat.IPromptMessage {
			return pm
		}), toolRuntimeMap, modelCaller, "builtin")

		taskScheduler.SetFunctionCallRunner(agentRunner)

		taskScheduler.Process(ctx)
	}
}

func (r *appAgentChatRunner) QueryAppAnnotationToReply(ctx context.Context, appRecord *po_entity.App, message *po_entity_chat.Message, query, accountID, invokeFrom string) (*po_entity_chat.MessageAnnotation, error) {
	return app_feature.NewAnnotationReplyFeature(r.ChatDomain, r.DatasetDomain, r.ProviderDomain, r.redis).Query(ctx, appRecord, message, query, accountID, invokeFrom)
}

func (r *appAgentChatRunner) InitPromptTools() (map[string]*biz_entity_agent.ToolRuntimeConfiguration, []*biz_entity.PromptMessageTool, error) {
	var (
		promptMessageTools []*biz_entity.PromptMessageTool
		toolInstance       = make(map[string]*biz_entity_agent.ToolRuntimeConfiguration)
	)

	for _, tool := range r.appConfig.AgentEntity.Tools {
		promptTool, toolRuntime, err := r.convertToolToPromptMessageTool(tool)

		if err != nil {
			return nil, nil, err
		}

		toolInstance[tool.ToolName] = toolRuntime
		promptMessageTools = append(promptMessageTools, promptTool)
	}

	return toolInstance, promptMessageTools, nil
}

func (r *appAgentChatRunner) convertToolToPromptMessageTool(tool *biz_entity_app_config.AgentToolEntity) (*biz_entity.PromptMessageTool, *biz_entity_agent.ToolRuntimeConfiguration, error) {
	toolRuntime, err := r.agentDomain.GetAgentToolRuntime(r.tenantID, r.appConfig.AppID, tool, "builtin")

	if err != nil {
		return nil, nil, err
	}

	promptMessageTool := &biz_entity.PromptMessageTool{
		Name:        tool.ToolName,
		Description: toolRuntime.Description.LLM,
		Parameters:  biz_entity.NewPromptMessageToolParameter(),
	}

	parameters := toolRuntime.GetAllRuntimeParameters()

	for _, parameter := range parameters {
		if parameter.Form != biz_entity_agent.LLMForm {
			continue
		}

		if parameter.Type == biz_entity_agent.FileType || parameter.Type == biz_entity_agent.FilesType || parameter.Type == biz_entity_agent.SystemFilesType {
			continue
		}

		var enum []string

		if parameter.Type == biz_entity_agent.SelectType {
			for _, opt := range parameter.Options {
				enum = append(enum, opt.Value)
			}
		}

		promptMessageTool.Parameters.Properties[parameter.Name] = &biz_entity.PromptMessageToolProperty{
			Type:        parameter.Type.AsNormalType(),
			Description: parameter.LLMDescription,
		}

		if len(enum) > 0 {
			promptMessageTool.Parameters.Properties[parameter.Name].Enum = enum
		}

		if parameter.Required {
			promptMessageTool.Parameters.Required = append(promptMessageTool.Parameters.Required, parameter.Name)
		}
	}

	return promptMessageTool, toolRuntime, nil
}
