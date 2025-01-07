package app_agent_chat_runner

import (
	"github.com/lunarianss/Luna/internal/api-server/core/app_chat/app_chat_runner"
)

type AppBaseAgentChatRunner struct {
	*app_chat_runner.AppBaseChatRunner
}

func NewAppBaseAgentChatRunner(appBaseRunner *app_chat_runner.AppBaseChatRunner) *AppBaseAgentChatRunner {
	return &AppBaseAgentChatRunner{
		AppBaseChatRunner: appBaseRunner,
	}
}
