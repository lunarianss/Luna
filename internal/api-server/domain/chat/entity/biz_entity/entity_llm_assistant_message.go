package biz_entity

import "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"

type AssistantPromptMessage struct {
	*po_entity.PromptMessage
}

func NewEmptyAssistantPromptMessage() *AssistantPromptMessage {
	return &AssistantPromptMessage{
		PromptMessage: &po_entity.PromptMessage{
			Content: "",
		},
	}
}

func NewAssistantPromptMessage(role po_entity.PromptMessageRole, content interface{}) *AssistantPromptMessage {
	return &AssistantPromptMessage{
		PromptMessage: &po_entity.PromptMessage{
			Content: content,
			Role:    role,
		},
	}
}
