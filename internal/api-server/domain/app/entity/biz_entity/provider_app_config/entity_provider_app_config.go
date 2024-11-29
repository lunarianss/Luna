// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package biz_entity

type ModelConfigEntity struct {
	Provider   string                 `json:"provider"`
	Model      string                 `json:"model"`
	Mode       string                 `json:"mode"`
	Parameters map[string]interface{} `json:"parameters"`
	Stop       []string               `json:"stop"`
}

type AdvancedChatMessageEntity struct {
	Text string `json:"text"`
	Role string `json:"role"` // Assuming PromptMessageRole is defined as string
}

type AdvancedChatPromptTemplateEntity struct {
	Messages []*AdvancedChatMessageEntity `json:"messages"`
}

type RolePrefixEntity struct {
	User      string `json:"user"`
	Assistant string `json:"assistant"`
}

type AdvancedCompletionPromptTemplateEntity struct {
	Prompt     string            `json:"prompt"`
	RolePrefix *RolePrefixEntity `json:"role_prefix"`
}

type PromptTemplateEntity struct {
	PromptType                       string                                  `json:"prompt_type"`
	SimplePromptTemplate             string                                  `json:"simple_prompt_template"`
	AdvancedChatPromptTemplate       *AdvancedChatPromptTemplateEntity       `json:"advanced_chat_prompt_template"`
	AdvancedCompletionPromptTemplate *AdvancedCompletionPromptTemplateEntity `json:"advanced_completion_prompt_template"`
}

type VariableEntity struct {
	Variable                 string             `json:"variable"`
	Label                    string             `json:"label"`
	Description              string             `json:"description"`
	Type                     VariableEntityType `json:"type"`
	Required                 bool               `json:"required"`
	MaxLength                int                `json:"max_length"`
	Options                  []string           `json:"options"`
	AllowedFileTypes         []string           `json:"allowed_file_types"`
	AllowedFileExtensions    []string           `json:"allowed_file_extensions"`
	AllowedFileUploadMethods []string           `json:"allowed_file_upload_methods"`
}

type ExternalDataVariableEntity struct {
	Variable string                 `json:"variable"`
	Type     string                 `json:"type"`
	Config   map[string]interface{} `json:"config"`
}

type RetrieveStrategy string

const (
	Single   RetrieveStrategy = "single"
	Multiple RetrieveStrategy = "multiple"
)

type SensitiveWordAvoidanceEntity struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

type TextToSpeechEntity struct {
	Enabled  bool   `json:"enabled"`
	Voice    string `json:"voice"`
	Language string `json:"language"`
}

type TracingConfigEntity struct {
	Enabled         bool   `json:"enabled"`
	TracingProvider string `json:"tracing_provider"`
}

type AppAdditionalFeatures struct {
	FileUpload                    string               `json:"file_upload"`
	OpeningStatement              string               `json:"opening_statement"`
	SuggestedQuestions            []string             `json:"suggested_questions"`
	SuggestedQuestionsAfterAnswer bool                 `json:"suggested_questions_after_answer"`
	ShowRetrieveSource            bool                 `json:"show_retrieve_source"`
	MoreLikeThis                  bool                 `json:"more_like_this"`
	SpeechToText                  bool                 `json:"speech_to_text"`
	TextToSpeech                  *TextToSpeechEntity  `json:"text_to_speech"`
	TraceConfig                   *TracingConfigEntity `json:"trace_config"`
}

type AppConfig struct {
	TenantID               string                        `json:"tenant_id"`
	AppID                  string                        `json:"app_id"`
	AppMode                string                        `json:"app_mode"` // Assuming AppMode is defined as string
	AdditionalFeatures     *AppAdditionalFeatures        `json:"additional_features"`
	Variables              []*VariableEntity             `json:"variables"`
	SensitiveWordAvoidance *SensitiveWordAvoidanceEntity `json:"sensitive_word_avoidance"`
}

type EasyUIBasedAppConfig struct {
	*AppConfig
	AppModelConfigFrom    EasyUIBasedAppModelConfigFrom `json:"app_model_config_from"`
	AppModelConfigID      string                        `json:"app_model_config_id"`
	AppModelConfig        *AppModelConfig               `json:"app_model_config_dict"`
	Model                 *ModelConfigEntity            `json:"model"`
	PromptTemplate        *PromptTemplateEntity         `json:"prompt_template"`
	ExternalDataVariables []ExternalDataVariableEntity  `json:"external_data_variables"`
}

type WorkflowUIBasedAppConfig struct {
	*AppConfig
	WorkflowID string `json:"workflow_id"`
}

type ChatAppConfig struct {
	*EasyUIBasedAppConfig
}
