// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

import (
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
)

type ListConversationQuery struct {
	Limit  int    `json:"limit" form:"limit" validate:"required"`
	Pinned *bool  `json:"pinned" form:"pinned" validate:"required"`
	LastID string `json:"last_id" form:"last_id"`
	SortBy string `json:"sort_by" form:"sort_by"`
}

func NewListConversationQuery() *ListConversationQuery {
	return &ListConversationQuery{
		SortBy: "-updated_at",
	}
}

type WebConversationDetail struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Inputs       map[string]interface{} `json:"inputs"`
	Introduction string                 `json:"introduction"`
	Status       string                 `json:"status"`
	CreatedAt    int64                  `json:"created_at"`
	UpdatedAt    int64                  `json:"updated_at"`
}

type WebMessageDetail struct {
	ID                 string                 `json:"id"`
	ConversationID     string                 `json:"conversation_id"`
	ParentMessageID    string                 `json:"parent_message_id"`
	Inputs             map[string]interface{} `json:"inputs"`
	Query              string                 `json:"query"`
	Answer             string                 `json:"answer"`
	Status             string                 `json:"status"`
	Error              string                 `json:"error"`
	MessageFiles       []string               `json:"message_files"`
	FeedBack           map[string]interface{} `json:"feedback"`
	RetrieverResources []any                  `json:"retriever_resources"`
	AgentThoughts      []any                  `json:"agent_thoughts"`
	CreatedAt          int64                  `json:"created_at"`
	UpdatedAt          int64                  `json:"updated_at"`
}

func MessageRecordToDetail(c *po_entity.Message) *WebMessageDetail {
	return &WebMessageDetail{
		ID:                 c.ID,
		ConversationID:     c.ConversationID,
		ParentMessageID:    c.ParentMessageID,
		Inputs:             c.Inputs,
		Query:              c.Query,
		Answer:             c.Answer,
		Status:             c.Status,
		AgentThoughts:      make([]any, 0),
		RetrieverResources: make([]any, 0),
		FeedBack:           make(map[string]any, 0),
		Error:              c.Error,
		CreatedAt:          c.CreatedAt,
		UpdatedAt:          c.UpdatedAt,
	}
}

func ConversationRecordToDetail(c *po_entity.Conversation) *WebConversationDetail {
	return &WebConversationDetail{
		ID:           c.ID,
		Name:         c.Name,
		Inputs:       c.Inputs,
		Introduction: c.Introduction,
		Status:       c.Status,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

type ListConversationResponse struct {
	Data    []*WebConversationDetail `json:"data"`
	Limit   int                      `json:"limit"`
	HasMore int                      `json:"has_more"`
	Count   int64                    `json:"count"`
}

type ListMessageResponse struct {
	Data    []*WebMessageDetail `json:"data"`
	Limit   int                 `json:"limit"`
	HasMore int                 `json:"has_more"`
	Count   int64               `json:"count"`
}

type ListMessageQuery struct {
	ConversationID string `json:"conversation_id" form:"conversation_id" validate:"required"`
	Limit          int    `json:"limit" form:"limit" validate:"required"`
	LastID         string `json:"last_id" form:"last_id"`
	FirstID        string `json:"first_id" form:"first_id"`
}

type ConversationIDUrl struct {
	ConversationID string `json:"conversationID" form:"conversationID" uri:"conversationID" validate:"required"`
}

type RenameConversationRequest struct {
	Name         string `json:"name" form:"name"`
	AutoGenerate bool   `json:"auto_generate" form:"auto_generate"`
}

func NewRenameConversationRequest() *RenameConversationRequest {
	return &RenameConversationRequest{
		AutoGenerate: false,
	}
}
