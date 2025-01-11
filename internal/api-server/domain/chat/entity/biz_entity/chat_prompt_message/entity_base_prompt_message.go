package biz_entity

import (
	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
)

// Base PromptMessage Only contains role and content
type PromptMessage struct {
	Role    PromptMessageRole `json:"role"`
	Content any               `json:"content,omitempty"`
	Name    string            `json:"name,omitempty"`
}

func NewSystemMessage(content any) *PromptMessage {
	return &PromptMessage{
		Role:    SYSTEM,
		Content: content,
	}
}

func NewAssistantMessage(content any) *PromptMessage {
	return &PromptMessage{
		Role:    ASSISTANT,
		Content: content,
	}
}

func NewUserMessage(content any) *PromptMessage {
	return &PromptMessage{
		Role:    USER,
		Content: content,
	}
}

func (pm *PromptMessage) GetRole() string {
	return string(pm.Role)
}

func (pm *PromptMessage) GetContent() string {
	return pm.Content.(string)
}
func (pm *PromptMessage) GetName() string {
	return pm.Name
}

func (msg *PromptMessage) ConvertToRequestData() (map[string]interface{}, error) {
	var requestData = make(map[string]interface{})

	if msg.Role == USER {
		switch content := msg.Content.(type) {
		case string:
			requestData["role"] = "user"
			requestData["content"] = content
		case []*PromptMessageContent:
			var subMessage []map[string]any
			for _, messageContent := range content {
				if messageContent.Type == TEXT {
					subMessageItem := map[string]any{
						"type": "text",
						"text": messageContent.Data,
					}
					subMessage = append(subMessage, subMessageItem)
				}
			}
			requestData["role"] = "user"
			requestData["content"] = subMessage
		default:
			return nil, errors.WithCode(code.ErrTypeOfPromptMessage, "value %T is not string or []*promptMessageContent type", msg.Content)
		}
	} else if msg.Role == ASSISTANT {
		requestData["role"] = "assistant"
		requestData["content"] = msg.Content
	} else if msg.Role == SYSTEM {
		requestData["role"] = "system"
		requestData["content"] = msg.Content
	}

	return requestData, nil
}
