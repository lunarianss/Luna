// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package biz_entity

import (
	biz_entity_app_config "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/provider_app_config"
	biz_entity_provider_config "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_configuration"
)

type BasedAppGenerateEntity interface {
	GetModel() string
	GetTaskID() string
	GetConversationID() string
	GetQuery() string
}

// Define types for your enums
type InvokeFrom string

const (
	ServiceAPI InvokeFrom = "service-api"
	WebApp     InvokeFrom = "web-app"
	Explore    InvokeFrom = "explore"
	Debugger   InvokeFrom = "debugger"
)

type CreatedByRole string

const (
	CreatedByRoleAccount CreatedByRole = "account"
	CreatedByRoleEndUser CreatedByRole = "end_user"
)

type PromptType string

const (
	SIMPLE   PromptType = "simple"
	ADVANCED PromptType = "advanced"
)

type UserFrom string

const (
	UserFromAccount UserFrom = "account"
	UserFromEndUser UserFrom = "end-user"
)

type WorkflowRunTriggeredFrom string

const (
	WorkflowRunTriggeredFromDebugging WorkflowRunTriggeredFrom = "debugging"
	WorkflowRunTriggeredFromAppRun    WorkflowRunTriggeredFrom = "app-run"
)

// AppGenerateEntity struct
type AppGenerateEntity struct {
	TaskID     string                           `json:"task_id"`
	AppConfig  *biz_entity_app_config.AppConfig `json:"app_config"`
	Inputs     map[string]interface{}           `json:"inputs"`
	UserID     string                           `json:"user_id"`
	Stream     bool                             `json:"stream"`
	InvokeFrom InvokeFrom                       `json:"invoke_from"`
	CallDepth  int                              `json:"call_depth"`
	Extras     map[string]interface{}           `json:"extras"`
}

// EasyUIBasedAppGenerateEntity struct
type EasyUIBasedAppGenerateEntity struct {
	*AppGenerateEntity
	AppConfig *biz_entity_app_config.EasyUIBasedAppConfig                  `json:"app_config"`
	ModelConf *biz_entity_provider_config.ModelConfigWithCredentialsEntity `json:"model_conf"`
	Query     string                                                       `json:"query"`
}

// ConversationAppGenerateEntity struct
type ConversationAppGenerateEntity struct {
	*AppGenerateEntity
	ConversationID  string `json:"conversation_id"`
	ParentMessageID string `json:"parent_message_id"`
}

// ChatAppGenerateEntity struct
type ChatAppGenerateEntity struct {
	*EasyUIBasedAppGenerateEntity
	ConversationID  string `json:"conversation_id"`
	ParentMessageID string `json:"parent_message_id"`
}

func (cag *ChatAppGenerateEntity) GetModel() string {
	return cag.ModelConf.Model
}

func (cag *ChatAppGenerateEntity) GetTaskID() string {
	return cag.EasyUIBasedAppGenerateEntity.TaskID
}

func (cag *ChatAppGenerateEntity) GetConversationID() string {
	return cag.ConversationID
}

func (cag *ChatAppGenerateEntity) GetQuery() string {
	return cag.Query
}

type AgentChatAppGenerateEntity struct {
	*EasyUIBasedAppGenerateEntity
	ConversationID  string `json:"conversation_id"`
	ParentMessageID string `json:"parent_message_id"`
	*biz_entity_app_config.AgentEntity
}

func (cag *AgentChatAppGenerateEntity) GetModel() string {
	return cag.EasyUIBasedAppGenerateEntity.ModelConf.Model
}

func (cag *AgentChatAppGenerateEntity) GetTaskID() string {
	return cag.EasyUIBasedAppGenerateEntity.TaskID
}

func (cag *AgentChatAppGenerateEntity) GetConversationID() string {
	return cag.ConversationID
}

func (cag *AgentChatAppGenerateEntity) GetQuery() string {
	return cag.Query
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

type VariableEntityType string

const (
	TextInput        VariableEntityType = "text-input"
	Select           VariableEntityType = "select"
	Paragraph        VariableEntityType = "paragraph"
	Number           VariableEntityType = "number"
	ExternalDataTool VariableEntityType = "external_data_tool"
	File             VariableEntityType = "file"
	FileList         VariableEntityType = "file-list"
)

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
