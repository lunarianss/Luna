// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

// ChatCreateMessage Dto
type CreateChatMessageUri struct {
	AppID string `uri:"appID" validate:"required"`
}

type AppModelConfigDtoEnable struct {
	Enable bool `json:"enable"`
}

// Model holds the model-specific configuration.
type ModelDto struct {
	Provider         string                 `json:"provider"`
	Name             string                 `json:"name"`
	Mode             string                 `json:"mode"`
	CompletionParams map[string]interface{} `json:"completion_params"`
}

type UserInputForm struct {
	TextInput *BaseTextUserInput `json:"text-input"`
}

type BaseTextUserInput struct {
	Label     string `json:"label"`
	Variable  string `json:"variable"`
	Required  bool   `json:"required"`
	MaxLength int    `json:"max_length"`
	Default   string `json:"default"`
}

type AppModelConfigDto struct {
	AppID                         string                  `json:"appId"`
	ModelID                       string                  `json:"model_id"`
	OpeningStatement              string                  `json:"opening_statement"`
	SuggestedQuestions            []string                `json:"suggested_questions"`
	SuggestedQuestionsAfterAnswer AppModelConfigDtoEnable `json:"suggested_questions_after_answer"`
	MoreLikeThis                  AppModelConfigDtoEnable `json:"more_like_this"`
	Model                         ModelDto                `json:"model"`
	UserInputForm                 []*UserInputForm        `json:"user_input_form"`
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
