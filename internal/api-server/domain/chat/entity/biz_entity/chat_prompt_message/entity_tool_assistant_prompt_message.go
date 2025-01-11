package biz_entity

import (
	"github.com/lunarianss/Luna/infrastructure/errors"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity/openai_standard_response"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
)

// AssistantPromptMessage With ToolCall
type AssistantPromptMessage struct {
	*PromptMessage
	ToolCalls []*biz_entity.ToolCall `json:"tool_calls,omitempty"`
}

func NewAssistantToolPromptMessage(content interface{}) *AssistantPromptMessage {
	return &AssistantPromptMessage{
		PromptMessage: &PromptMessage{
			Content: content,
			Role:    ASSISTANT,
		},
	}
}

func (msg *AssistantPromptMessage) ConvertToRequestData() (map[string]interface{}, error) {
	var requestData = make(map[string]interface{})

	switch content := msg.Content.(type) {
	case string:
		requestData["role"] = "assistant"
		if content != "" {
			requestData["content"] = content
		}

		if len(msg.ToolCalls) > 0 {
			requestData["tool_calls"] = msg.ToolCalls
		}
	default:
		return nil, errors.WithCode(code.ErrTypeOfPromptMessage, "value %T is not string type", msg.Content)
	}

	return requestData, nil
}

func (ap *AssistantPromptMessage) GetRole() string {
	return string(ap.Role)
}

func (pm *AssistantPromptMessage) GetContent() string {
	return pm.Content.(string)
}
func (pm *AssistantPromptMessage) GetName() string {
	return pm.Name
}
