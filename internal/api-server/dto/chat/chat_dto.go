// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

import biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/provider_app_config"

type Speech2TextResp struct {
	Text string `json:"text"`
}

type TextToAudioRequest struct {
	MessageID string `json:"message_id"`
	Streaming bool   `json:"streaming"  validate:"required"`
	Text      string `json:"text"       validate:"required"`
	Voice     string `json:"voice"`
}

// FeedbackStats
type FeedBackStats struct {
	Like    int `json:"like"`
	Dislike int `json:"dislike"`
}

func NewFeedBackStats() *FeedBackStats {
	return &FeedBackStats{
		Like:    0,
		Dislike: 0,
	}
}

// ChatCreateMessage Dto
type CreateChatMessageUri struct {
	AppID string `uri:"appID" validate:"required"`
}

// ChatCreateMessage Dto
type DetailConversationUri struct {
	AppID          string `uri:"appID" validate:"required"`
	ConversationID string `uri:"conversationID" validate:"required"`
}

type PromptMessage struct {
	Role    string `json:"role"`
	Content string `json:"text"`
	Name    string `json:"name"`
}

type ListChatMessageItem struct {
	ID                      string                 `json:"id"`
	ConversationID          string                 `json:"conversation_id"`
	Inputs                  map[string]interface{} `json:"inputs"`
	Query                   string                 `json:"query"`
	Message                 []*PromptMessage       `json:"message"`
	MessageTokens           int64                  `json:"message_tokens"`
	MessageUnitPrice        float64                `json:"message_unit_price"`
	Answer                  string                 `json:"answer"`
	AnswerTokens            int64                  `json:"answer_tokens"`
	ProviderResponseLatency float64                `json:"provider_response_latency"`
	TotalPrice              float64                `json:"total_price"`
	Currency                string                 `json:"currency"`
	FromSource              string                 `json:"from_source"`
	FromEndUserID           string                 `json:"from_end_user_id"`
	FromAccountID           string                 `json:"from_account_id"`
	CreatedAt               int64                  `json:"created_at"`
	MessagePriceUnit        float64                `json:"message_price_unit"`
	AnswerPriceUnit         float64                `json:"answer_price_unit"`
	WorkflowRunID           string                 `json:"workflow_run_id"`
	Status                  string                 `json:"status"`
	Error                   string                 `json:"error"`
	MessageMetadata         map[string]interface{} `json:"message_metadata"`
	InvokeFrom              string                 `json:"invoke_from"`
	FeedBacks               []string               `json:"feedbacks"`
	AgentThoughts           []string               `json:"agent_thoughts"`
	MessageFiles            []string               `json:"message_files"`
	Metadata                map[string]interface{} `json:"metadata"`
	ParentMessageID         string                 `json:"parent_message_id"`
}

type ListChatMessagesResponse struct {
	Limit   int                    `json:"limit"`
	HasMore bool                   `json:"has_more"`
	Data    []*ListChatMessageItem `json:"data"`
	Count   int64                  `json:"count"`
}

// ChatCreateMessage Dto
type ListChatMessageQuery struct {
	ConversationID string `form:"conversation_id" validate:"required"`
	FirstID        string `form:"first_id"`
	Limit          int    `form:"limit"`
}

func NewListChatMessageQuery() *ListChatMessageQuery {
	return &ListChatMessageQuery{
		Limit: 20,
	}
}

type ListChatConversationQuery struct {
	Keyword         string `json:"keyword" form:"keyword"`
	Start           string `json:"start" form:"start"`
	End             string `json:"end" form:"end"`
	MessageCountGte int    `json:"message_count_gte" form:"message_count_gte"`
	Page            int    `json:"page" form:"page"`
	Limit           int    `json:"limit" form:"limit"`
	SortBy          string `json:"sort_by" form:"sort_by"`
}

func NewListChatConversationQuery() *ListChatConversationQuery {
	return &ListChatConversationQuery{
		Limit:  10,
		Page:   1,
		SortBy: "-created_at",
	}
}

type ListChatConversationItem struct {
	ID                   string                     `json:"id"`
	Status               string                     `json:"status"`
	FromSource           string                     `json:"from_source"`
	FromEndUserID        string                     `json:"from_end_user_id"`
	FromEndUserSessionID string                     `json:"from_end_user_session_id"`
	FromAccountID        string                     `json:"from_account_id"`
	FromAccountName      string                     `json:"from_account_name"`
	Name                 string                     `json:"name"`
	Summary              string                     `json:"summary"`
	ReadAt               int64                      `json:"read_at"`
	CreatedAt            int64                      `json:"created_at"`
	UpdatedAt            int64                      `json:"updated_at"`
	Annotated            bool                       `json:"annotated"`
	ModelConfig          *biz_entity.AppModelConfig `json:"model_config"`
	MessageCount         int64                      `json:"message_count"`
	UserFeedbackStats    *FeedBackStats             `json:"user_feedback_stats"`
	AdminFeedbackStats   *FeedBackStats             `json:"admin_feedback_stats"`
}

type ListChatConversationResponse struct {
	HasMore bool                        `json:"has_more"`
	Data    []*ListChatConversationItem `json:"data"`
	Page    int                         `json:"page"`
	Limit   int                         `json:"limit"`
	Total   int64                       `json:"total"`
}

type AppModelConfigDtoEnable struct {
	Enabled bool `json:"enabled"`
}

// Model holds the model-specific configuration.
type ModelDto struct {
	Provider         string                 `json:"provider"`
	Name             string                 `json:"name"`
	Mode             string                 `json:"mode"`
	CompletionParams map[string]interface{} `json:"completion_params"`
}

type UserInput struct {
	Label     string   `json:"label"`
	Variable  string   `json:"variable"`
	Required  bool     `json:"required"`
	MaxLength int      `json:"max_length"`
	Default   string   `json:"default"`
	Options   []string `json:"options"`
}

type UserInputForm map[string]*UserInput

type AppModelConfigDto struct {
	AppID                         string                  `json:"appId"`
	ModelID                       string                  `json:"model_id"`
	OpeningStatement              string                  `json:"opening_statement"`
	SuggestedQuestions            []string                `json:"suggested_questions"`
	SuggestedQuestionsAfterAnswer AppModelConfigDtoEnable `json:"suggested_questions_after_answer"`
	MoreLikeThis                  AppModelConfigDtoEnable `json:"more_like_this"`
	Model                         ModelDto                `json:"model"`
	UserInputForm                 []UserInputForm         `json:"user_input_form"`
	PrePrompt                     string                  `json:"pre_prompt"`
	AgentMode                     map[string]interface{}  `json:"agent_mode"`
	SpeechToText                  AppModelConfigDtoEnable `json:"speech_to_text"`
	SensitiveWordAvoidance        map[string]interface{}  `json:"sensitive_word_avoidance"`
	RetrieverResource             AppModelConfigDtoEnable `json:"retriever_resource"`
	DatasetQueryVariable          string                  `json:"dataset_query_variable"`
	PromptType                    string                  `json:"prompt_type"`
	ChatPromptConfig              map[string]interface{}  `json:"chat_prompt_config"`
	CompletionPromptConfig        map[string]interface{}  `json:"completion_prompt_config"`
	DatasetConfigs                map[string]interface{}  `json:"dataset_configs"`
	FileUpload                    map[string]interface{}  `json:"file_upload"`
	TextToSpeech                  AppModelConfigDtoEnable `json:"text_to_speech"`
	ExternalDataTools             []string                `json:"external_data_tools" `
	Configs                       map[string]interface{}  `json:"configs"`
}

type CreateChatMessageBody struct {
	ResponseMode                 string                 `json:"response_mode" validate:"required"`
	ConversationID               string                 `json:"conversation_id"`
	Query                        string                 `json:"query" validate:"required"`
	Files                        []string               `json:"files"`
	Inputs                       map[string]interface{} `json:"inputs" `
	ModelConfig                  AppModelConfigDto      `json:"model_config"`
	ParentMessageId              string                 `json:"parent_message_id"`
	AutoGenerateConversationName bool                   `json:"auto_generate_conversation_name"`
}

type ServiceCreateChatMessageBody struct {
	ResponseMode                 string                 `json:"response_mode" validate:"required"`
	ConversationID               string                 `json:"conversation_id"`
	Query                        string                 `json:"query" validate:"required"`
	Files                        []string               `json:"files"`
	Inputs                       map[string]interface{} `json:"inputs" `
	ModelConfig                  AppModelConfigDto      `json:"model_config"`
	ParentMessageId              string                 `json:"parent_message_id"`
	User                         string                 `json:"user"`
	AutoGenerateConversationName bool                   `json:"auto_generate_conversation_name"`
}
