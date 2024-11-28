// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

// ChatCreateMessage Dto
type CreateChatMessageUri struct {
	AppID string `uri:"appID" validate:"required"`
}

type CreateChatMessageBody struct {
	ResponseMode                 string                 `json:"response_mode" validate:"required"`
	ConversationID               string                 `json:"conversation_id"`
	Query                        string                 `json:"query" validate:"required"`
	Files                        []string               `json:"files"`
	Inputs                       map[string]interface{} `json:"inputs" `
	ModelConfig                  map[string]interface{} `json:"model_config"`
	ParentMessageId              string                 `json:"parent_message_id"`
	AutoGenerateConversationName bool                   `json:"auto_generate_conversation_name"`
}
