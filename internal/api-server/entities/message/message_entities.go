package message

import (
	"fmt"

	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
)

type PromptMessageRole string

const (
	SYSTEM    PromptMessageRole = "system"
	USER      PromptMessageRole = "user"
	ASSISTANT PromptMessageRole = "assistant"
	TOOL      PromptMessageRole = "tool"
)

type PromptMessageContentType string

const (
	TEXT  PromptMessageContentType = "text"
	IMAGE PromptMessageContentType = "image"
	AUDIO PromptMessageContentType = "audio"
)

type PromptMessageContent struct {
	Type PromptMessageContentType `json:"type"`
	Data string                   `json:"data"`
}

type PromptMessage struct {
	Role    PromptMessageRole `json:"role"`
	Content interface{}       `json:"content"`
	Name    string            `json:"name"`
}

type AssistantPromptMessage struct {
	*PromptMessage
}

func (msg *PromptMessage) IsEmpty() bool {
	return msg.Content == ""
}

func (msg *PromptMessage) ConvertToRequestData() (map[string]interface{}, error) {
	var requestData = make(map[string]interface{})

	if msg.Role == USER {
		switch v := msg.Content.(type) {
		case string:
			requestData["role"] = "user"
			requestData["content"] = v
		case []*PromptMessageContent:
			var subMessage []map[string]interface{}
			for _, messageContent := range v {
				if messageContent.Type == TEXT {
					subMessageItem := map[string]interface{}{
						"type": "text",
						"text": messageContent.Data,
					}
					subMessage = append(subMessage, subMessageItem)
				}
			}
			requestData["role"] = "user"
			requestData["content"] = subMessage
		default:
			return nil, errors.WithCode(code.ErrTypeOfPromptMessage, fmt.Sprintf("the type %T is neither string or []*promptMessageContent", v))
		}
	}

	return requestData, nil

}
