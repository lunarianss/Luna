package dto

import "github.com/lunarianss/Luna/internal/api-server/model/v1"

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
	ID           string                 `gorm:"column:id" json:"id"`
	Name         string                 `gorm:"column:name" json:"name"`
	Inputs       map[string]interface{} `gorm:"column:inputs;serializer:json" json:"inputs"`
	Introduction string                 `gorm:"column:introduction" json:"introduction"`
	Status       string                 `gorm:"column:status" json:"status"`
	CreatedAt    int64                  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    int64                  `gorm:"column:updated_at" json:"updated_at"`
}

func ConversationRecordToDetail(c *model.Conversation) *WebConversationDetail {
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

type ListMessageUrl struct {
	ConversationID string `json:"conversation_id" form:"conversation_id" validate:"required"`
	Limit          int    `json:"limit" form:"limit" validate:"required"`
}
