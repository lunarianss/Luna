package biz_entity

import (
	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
)

type ToolPromptMessage struct {
	*PromptMessage
	ToolCallID string `json:"tool_call_id"`
}

func (pm *ToolPromptMessage) GetRole() string {
	return string(pm.Role)
}

func (pm *ToolPromptMessage) GetContent() string {
	return pm.Content.(string)
}
func (pm *ToolPromptMessage) GetName() string {
	return pm.Name
}

func (msg *ToolPromptMessage) ConvertToRequestData() (map[string]interface{}, error) {
	var requestData = make(map[string]interface{})

	switch content := msg.Content.(type) {
	case string:
		requestData["role"] = "tool"
		requestData["content"] = content
		requestData["tool_call_id"] = msg.ToolCallID
	default:
		return nil, errors.WithCode(code.ErrTypeOfPromptMessage, "value %T is not string type", msg.Content)
	}

	return requestData, nil
}
