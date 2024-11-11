package dto

import "github.com/lunarianss/Luna/internal/api-server/model/v1"

// ChatCreateMessage Dto
type CreateChatMessageUri struct {
	AppID string `uri:"appID" validate:"required"`
}

type CreateChatMessageBody struct {
	ResponseMode    string                 `json:"response_mode" validate:"required"`
	ConversationID  string                 `json:"conversation_id"`
	Query           string                 `json:"query" validate:"required"`
	Files           []string               `json:"files"`
	Inputs          map[string]interface{} `json:"inputs" `
	ModelConfig     *model.AppModelConfig  `json:"model_config"`
	ParentMessageId string                 `json:"parent_message_id"`
}
