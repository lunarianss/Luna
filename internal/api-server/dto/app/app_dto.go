package dto

import "github.com/lunarianss/Luna/internal/api-server/model/v1"

// Create App Input Dto
type CreateAppRequest struct {
	Name           string `json:"name" validate:"required"`
	Mode           string `json:"mode" validate:"required"`
	Icon           string `json:"icon" validate:"required"`
	Description    string `json:"description"`
	IconType       string `json:"icon_type"`
	IconBackground string `json:"icon_background"`
	ApiRph         int    `json:"api_rph"`
	ApiRpm         int    `json:"api_rpm"`
}

// Create App Response Dto
type CreateAppResponse struct {
	*model.App
	ModelConfig *model.AppModelConfig `json:"model_config"`
}

// ChatCreateMessage Dto
type CreateChatMessageUri struct {
	AppID string `uri:"appID" validate:"required"`
}

type CreateChatMessageBody struct {
	ResponseMode    string                `json:"response_mode" validate:"required"`
	ConversationID  string                `json:"conversation_id" validate:"required"`
	Query           string                `json:"query" validate:"required"`
	Files           []string              `json:"files"`
	Inputs          []interface{}         `json:"inputs" `
	ModelConfig     *model.AppModelConfig `json:"model_config"`
	ParentMessageId string                `json:"parent_message_id"`
}
