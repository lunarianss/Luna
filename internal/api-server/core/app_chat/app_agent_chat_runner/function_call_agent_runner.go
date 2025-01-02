package app_agent_chat_runner

import (
	"github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_registry"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	biz_entity_app_generate "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_generate"
)

type FunctionCallAgentRunner struct {
	applicationGenerateEntity *biz_entity_app_generate.AgentChatAppGenerateEntity

	conversation *po_entity.Conversation

	modelInstance model_registry.IModelRegistryCall
}

func NewFunctionCallAgentRunner(applicationGenerateEntity *biz_entity_app_generate.AgentChatAppGenerateEntity, conversation *po_entity.Conversation, modelInstance model_registry.IModelRegistryCall) *FunctionCallAgentRunner {
	return &FunctionCallAgentRunner{
		applicationGenerateEntity: applicationGenerateEntity,
		conversation:              conversation,
		modelInstance:             modelInstance,
	}
}

func (fca *FunctionCallAgentRunner) Run(message *po_entity.Message, query string) {
	

}
