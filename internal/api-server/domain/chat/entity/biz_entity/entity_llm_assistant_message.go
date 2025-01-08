package biz_entity

import (
	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
)

type AssistantPromptMessage struct {
	*po_entity.PromptMessage
	ToolCalls []*ToolCall `json:"tool_calls"`
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

func NewEmptyAssistantPromptMessage() *AssistantPromptMessage {
	return &AssistantPromptMessage{
		PromptMessage: &po_entity.PromptMessage{
			Content: "",
		},
	}
}

func NewAssistantPromptMessage(content interface{}) *AssistantPromptMessage {
	return &AssistantPromptMessage{
		PromptMessage: &po_entity.PromptMessage{
			Content: content,
			Role:    po_entity.ASSISTANT,
		},
	}
}

type PromptMessageToolProperty struct {
	Type        string   `json:"string,omitempty"`
	Description string   `json:"description,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

type PromptMessageToolProperties map[string]*PromptMessageToolProperty

type PromptMessageToolParameter struct {
	Type       string                      `json:"type"`
	Properties PromptMessageToolProperties `json:"properties"`
	Required   []string                    `json:"required,omitempty"`
}

func NewPromptMessageToolParameter() *PromptMessageToolParameter {
	return &PromptMessageToolParameter{
		Required:   make([]string, 0),
		Type:       "object",
		Properties: make(PromptMessageToolProperties, 0),
	}
}

type PromptMessageTool struct {
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	Parameters  *PromptMessageToolParameter `json:"parameters"`
}

type PromptMessageFunction struct {
	Type     string             `json:"type"`
	Function *PromptMessageTool `json:"function"`
}

func NewFunctionTools(function *PromptMessageTool) *PromptMessageFunction {
	return &PromptMessageFunction{
		Type:     "function",
		Function: function,
	}
}
