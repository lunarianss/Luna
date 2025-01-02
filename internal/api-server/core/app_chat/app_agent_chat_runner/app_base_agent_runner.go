package app_agent_chat_runner

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/core/app_chat/app_chat_runner"
	"github.com/lunarianss/Luna/internal/api-server/core/app_chat/token_buffer_memory"
	prompt "github.com/lunarianss/Luna/internal/api-server/core/app_prompt"

	biz_entity_app_config "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/provider_app_config"
	po_entity_app "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	po_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	biz_entity_provider_config "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_configuration"
)

type AppBaseAgentChatRunner struct {
	*app_chat_runner.AppBaseChatRunner
}

func NewAppBaseAgentChatRunner(appBaseRunner *app_chat_runner.AppBaseChatRunner) *AppBaseAgentChatRunner {
	return &AppBaseAgentChatRunner{
		AppBaseChatRunner: appBaseRunner,
	}
}

func (r *AppBaseAgentChatRunner) OrganizePromptMessage(ctx context.Context, appRecord *po_entity_app.App, modelConfig *biz_entity_provider_config.ModelConfigWithCredentialsEntity, promptTemplateEntity *biz_entity_app_config.PromptTemplateEntity, inputs map[string]interface{}, files []string, query string, context string, memory token_buffer_memory.ITokenBufferMemory) ([]*po_entity_chat.PromptMessage, []string, error) {

	var (
		promptMessages []*po_entity_chat.PromptMessage
		stop           []string
		err            error
	)
	if promptTemplateEntity.PromptType == string(biz_entity_app_config.SIMPLE) {
		simplePrompt := prompt.SimplePromptTransform{}

		promptMessages, stop, err = simplePrompt.GetPrompt(po_entity_app.AppMode(appRecord.Mode), promptTemplateEntity, inputs, query, files, context, memory, modelConfig)

		if err != nil {
			return nil, nil, err
		}
	}

	return promptMessages, stop, err
}
