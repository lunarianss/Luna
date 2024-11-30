// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package po_entity

import (
	"github.com/google/uuid"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/infrastructure/field"
	"gorm.io/gorm"
)

type AppMode string

const (
	COMPLETION    AppMode = "completion"
	WORKFLOW      AppMode = "workflow"
	CHAT          AppMode = "chat"
	ADVANCED_CHAT AppMode = "advanced-chat"
	AGENT_CHAT    AppMode = "agent-chat"
	CHANNEL       AppMode = "channel"
)

// App represents the app table in the database
type App struct {
	ID                  string        `gorm:"column:id" json:"id"`
	TenantID            string        `gorm:"column:tenant_id" json:"tenant_id"`
	Name                string        `gorm:"column:name" json:"name"`
	Mode                string        `gorm:"column:mode" json:"mode"`
	Icon                string        `gorm:"column:icon" json:"icon"`
	IconBackground      string        `gorm:"column:icon_background" json:"icon_background"`
	AppModelConfigID    string        `gorm:"column:app_model_config_id" json:"app_model_config_id"`
	Status              string        `gorm:"column:status;default:normal" json:"status"`
	EnableSite          field.BitBool `gorm:"column:enable_site" json:"enable_site"`
	EnableAPI           field.BitBool `gorm:"column:enable_api" json:"enable_api"`
	APIRpm              int           `gorm:"column:api_rpm" json:"api_rpm"`
	APIRph              int           `gorm:"column:api_rph" json:"api_rph"`
	IsDemo              field.BitBool `gorm:"column:is_demo" json:"is_demo"`
	IsPublic            field.BitBool `gorm:"column:is_public" json:"is_public"`
	CreatedAt           int           `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           int           `gorm:"column:updated_at" json:"updated_at"`
	IsUniversal         field.BitBool `gorm:"column:is_universal" json:"is_universal"`
	WorkflowID          string        `gorm:"column:workflow_id" json:"workflow_id"`
	Description         string        `gorm:"column:description" json:"description"`
	Tracing             string        `gorm:"column:tracing" json:"tracing"`
	MaxActiveRequests   int           `gorm:"column:max_active_requests" json:"max_active_requests"`
	IconType            string        `gorm:"column:icon_type" json:"icon_type"`
	CreatedBy           string        `gorm:"column:created_by" json:"created_by"`
	UpdatedBy           string        `gorm:"column:updated_by" json:"updated_by"`
	UseIconAsAnswerIcon field.BitBool `gorm:"column:use_icon_as_answer_icon" json:"use_icon_as_answer_icon"`
}

func (a *App) TableName() string {
	return "apps"
}

func (a *App) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.NewString()
	return
}

type AppModelConfigEnable struct {
	Enable bool `json:"enable"`
}

type AppModelConfig struct {
	ID                            string                              `json:"id" gorm:"column:id"`
	AppID                         string                              `json:"app_id" gorm:"column:app_id"`
	Provider                      string                              `json:"provider" gorm:"column:provider"`
	ModelID                       string                              `json:"model_id" gorm:"column:model_id"`
	Configs                       map[string]interface{}              `json:"configs" gorm:"column:configs;serializer:json"`
	CreatedAt                     int64                               `json:"created_at" gorm:"column:created_at"`
	UpdatedAt                     int64                               `json:"updated_at" gorm:"column:updated_at"`
	OpeningStatement              map[string]interface{}              `json:"opening_statement" gorm:"column:opening_statement;serializer:json"`
	SuggestedQuestions            []string                            `json:"suggested_questions" gorm:"column:suggested_questions;serializer:json"`
	SuggestedQuestionsAfterAnswer AppModelConfigEnable                `json:"suggested_questions_after_answer" gorm:"column:suggested_questions_after_answer;serializer:json"`
	MoreLikeThis                  AppModelConfigEnable                `json:"more_like_this" gorm:"column:more_like_this;serializer:json"`
	Model                         biz_entity.Model                    `json:"model" gorm:"column:model;serializer:json"`
	UserInputForm                 []map[string]map[string]interface{} `json:"user_input_form" gorm:"column:user_input_form;serializer:json"`
	PrePrompt                     string                              `json:"pre_prompt" gorm:"column:pre_prompt;serializer:json"`
	AgentMode                     map[string]interface{}              `json:"agent_mode" gorm:"column:agent_mode;serializer:json"`
	SpeechToText                  AppModelConfigEnable                `json:"speech_to_text" gorm:"column:speech_to_text;serializer:json"`
	SensitiveWordAvoidance        map[string]interface{}              `json:"sensitive_word_avoidance" gorm:"column:sensitive_word_avoidance;serializer:json"`
	RetrieverResource             AppModelConfigEnable                `json:"retriever_resource" gorm:"column:retriever_resource;serializer:json"`
	DatasetQueryVariable          string                              `json:"dataset_query_variable" gorm:"column:dataset_query_variable;serializer:json"`
	PromptType                    string                              `json:"prompt_type" gorm:"column:prompt_type"`
	ChatPromptConfig              map[string]interface{}              `json:"chat_prompt_config" gorm:"column:chat_prompt_config;serializer:json"`
	CompletionPromptConfig        map[string]interface{}              `json:"completion_prompt_config" gorm:"column:completion_prompt_config;serializer:json"`
	DatasetConfigs                map[string]interface{}              `json:"dataset_configs" gorm:"column:dataset_configs;serializer:json"`
	ExternalDataTools             []string                            `json:"external_data_tools" gorm:"column:external_data_tools;serializer:json"`
	FileUpload                    map[string]map[string]interface{}   `json:"file_upload" gorm:"column:file_upload;serializer:json"`
	TextToSpeech                  AppModelConfigEnable                `json:"text_to_speech" gorm:"column:text_to_speech;serializer:json"`
	CreatedBy                     string                              `json:"created_by" gorm:"column:created_by"`
	UpdatedBy                     string                              `json:"updated_by" gorm:"column:updated_by"`
}

func (a *AppModelConfig) TableName() string {
	return "app_model_configs"
}

func (a *AppModelConfig) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.NewString()
	return
}
