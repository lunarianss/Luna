// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

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
	Content any               `json:"content"`
	Name    string            `json:"name"`
}

func NewSystemMessage(content any) *PromptMessage {
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

func (msg *PromptMessage) ConvertToRequestData() (map[string]interface{}, error) {
	var requestData = make(map[string]interface{})

	if msg.Role == USER {
		switch content := msg.Content.(type) {
		case string:
			requestData["role"] = "user"
			requestData["content"] = content
		case []*PromptMessageContent:
			var subMessage []map[string]string
			for _, messageContent := range content {
				if messageContent.Type == TEXT {
					subMessageItem := map[string]string{
						"type": "text",
						"text": messageContent.Data,
					}
					subMessage = append(subMessage, subMessageItem)
				}
			}
			requestData["role"] = "user"
			requestData["content"] = subMessage
		default:
			return nil, errors.WithCode(code.ErrTypeOfPromptMessage, fmt.Sprintf("value %T is not string or []*promptMessageContent type", msg.Content))
		}
	}
	return requestData, nil
}

type AssistantPromptMessage struct {
	*PromptMessage
}

func NewEmptyAssistantPromptMessage() *AssistantPromptMessage {
	return &AssistantPromptMessage{
		PromptMessage: &PromptMessage{
			Content: "",
		},
	}
}

func NewAssistantPromptMessage(role PromptMessageRole, content interface{}) *AssistantPromptMessage {
	return &AssistantPromptMessage{
		PromptMessage: &PromptMessage{
			Content: content,
			Role:    role,
		},
	}
}
