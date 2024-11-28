// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package po_entity

import (
	"github.com/google/uuid"
	"github.com/lunarianss/Luna/internal/api-server/entities/message"
	"github.com/lunarianss/Luna/internal/pkg/field"
	"gorm.io/gorm"
)

type Conversation struct {
	ID                      string                 `gorm:"column:id" json:"id"`
	AppID                   string                 `gorm:"column:app_id" json:"app_id"`
	AppModelConfigID        string                 `gorm:"column:app_model_config_id" json:"app_model_config_id"`
	ModelProvider           string                 `gorm:"column:model_provider" json:"model_provider"`
	OverrideModelConfigs    map[string]interface{} `gorm:"column:override_model_configs;serializer:json" json:"override_model_configs"`
	ModelID                 string                 `gorm:"column:model_id" json:"model_id"`
	Mode                    string                 `gorm:"column:mode" json:"mode"`
	Name                    string                 `gorm:"column:name" json:"name"`
	Summary                 string                 `gorm:"column:summary" json:"summary"`
	Inputs                  map[string]interface{} `gorm:"column:inputs;serializer:json" json:"inputs"`
	Introduction            string                 `gorm:"column:introduction" json:"introduction"`
	SystemInstruction       string                 `gorm:"column:system_instruction" json:"system_instruction"`
	SystemInstructionTokens int                    `gorm:"column:system_instruction_tokens" json:"system_instruction_tokens"`
	Status                  string                 `gorm:"column:status" json:"status"`
	InvokeFrom              string                 `gorm:"column:invoke_from" json:"invoke_from"`
	FromSource              string                 `gorm:"column:from_source" json:"from_source"`
	FromEndUserID           string                 `gorm:"column:from_end_user_id" json:"from_end_user_id"`
	FromAccountID           string                 `gorm:"column:from_account_id" json:"from_account_id"`
	ReadAt                  int64                  `gorm:"column:read_at" json:"read_at"`
	ReadAccountID           string                 `gorm:"column:read_account_id" json:"read_account_id"`
	DialogueCount           int                    `gorm:"column:dialogue_count" json:"dialogue_count"`
	CreatedAt               int64                  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt               int64                  `gorm:"column:updated_at" json:"updated_at"`
	IsDeleted               field.BitBool          `gorm:"column:is_deleted" json:"is_deleted"`
}

func (a *Conversation) TableName() string {
	return "conversations"
}

func (a *Conversation) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.NewString()
	return
}

type Message struct {
	ID                      string                   `gorm:"column:id" json:"id"`
	AppID                   string                   `gorm:"column:app_id" json:"app_id"`
	ModelProvider           string                   `gorm:"column:model_provider" json:"model_provider"`
	ModelID                 string                   `gorm:"column:model_id" json:"model_id"`
	OverrideModelConfigs    map[string]interface{}   `gorm:"column:override_model_configs;serializer:json" json:"override_model_configs"`
	ConversationID          string                   `gorm:"column:conversation_id" json:"conversation_id"`
	Inputs                  map[string]interface{}   `gorm:"column:inputs;serializer:json" json:"inputs"`
	Query                   string                   `gorm:"column:query" json:"query"`
	Message                 []*message.PromptMessage `gorm:"column:message;serializer:json" json:"message"`
	MessageTokens           int                      `gorm:"column:message_tokens" json:"message_tokens"`
	MessageUnitPrice        float64                  `gorm:"column:message_unit_price" json:"message_unit_price"`
	Answer                  string                   `gorm:"column:answer" json:"answer"`
	AnswerTokens            int                      `gorm:"column:answer_tokens" json:"answer_tokens"`
	AnswerUnitPrice         float64                  `gorm:"column:answer_unit_price" json:"answer_unit_price"`
	ProviderResponseLatency float64                  `gorm:"column:provider_response_latency" json:"provider_response_latency"`
	TotalPrice              float64                  `gorm:"column:total_price" json:"total_price"`
	Currency                string                   `gorm:"column:currency" json:"currency"`
	FromSource              string                   `gorm:"column:from_source" json:"from_source"`
	FromEndUserID           string                   `gorm:"column:from_end_user_id" json:"from_end_user_id"`
	FromAccountID           string                   `gorm:"column:from_account_id" json:"from_account_id"`
	CreatedAt               int64                    `gorm:"column:created_at" json:"created_at"`
	UpdatedAt               int64                    `gorm:"column:updated_at" json:"updated_at"`
	AgentBased              field.BitBool            `gorm:"column:agent_based" json:"agent_based"`
	MessagePriceUnit        float64                  `gorm:"column:message_price_unit" json:"message_price_unit"`
	AnswerPriceUnit         float64                  `gorm:"column:answer_price_unit" json:"answer_price_unit"`
	WorkflowRunID           string                   `gorm:"column:workflow_run_id" json:"workflow_run_id"`
	Status                  string                   `gorm:"column:status" json:"status"`
	Error                   string                   `gorm:"column:error" json:"error"`
	MessageMetadata         map[string]interface{}   `gorm:"column:message_metadata;serializer:json" json:"message_metadata"`
	InvokeFrom              string                   `gorm:"column:invoke_from" json:"invoke_from"`
	ParentMessageID         string                   `gorm:"column:parent_message_id" json:"parent_message_id"`
}

func (a *Message) TableName() string {
	return "messages"
}

func (a *Message) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.NewString()
	return
}

type PinnedConversation struct {
	ID             string `gorm:"column:id" json:"id"`
	AppID          string `gorm:"column:app_id" json:"app_id"`
	ConversationID string `gorm:"column:conversation_id" json:"conversation_id"`
	CreatedByRole  string `gorm:"column:created_by_role;default:end_user" json:"created_by_role"`
	CreatedBy      string `gorm:"column:created_by" json:"create_by"`
	CreatedAt      int64  `gorm:"column:created_at" json:"created_at"`
}

func (a *PinnedConversation) TableName() string {
	return "pinned_conversations"
}

func (a *PinnedConversation) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.NewString()
	return
}
