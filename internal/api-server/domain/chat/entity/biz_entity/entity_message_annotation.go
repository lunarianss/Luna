package biz_entity

import (
	"github.com/lunarianss/Luna/internal/api-server/domain/account/entity/po_entity"
	po_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
)

type BizMessageAnnotation struct {
	ID             string             `json:"id"`
	AppID          string             `json:"app_id"`
	ConversationID string             `json:"conversation_id"`
	MessageID      string             `json:"message_id"`
	Question       string             `json:"question"`
	Content        string             `json:"content"`
	HitCount       int                `json:"hit_count"`
	AccountID      string             `json:"account_id"`
	CreatedAt      int64              `json:"created_at"`
	UpdatedAt      int64              `json:"updated_at"`
	Account        *po_entity.Account `json:"account"`
}

func ConvertToBizMessageAnnotation(annotation *po_chat.MessageAnnotation, account *po_entity.Account) *BizMessageAnnotation {
	return &BizMessageAnnotation{
		ID:             annotation.ID,
		AppID:          annotation.AppID,
		ConversationID: annotation.ConversationID,
		AccountID:      annotation.AccountID,
		MessageID:      annotation.MessageID,
		Question:       annotation.Question,
		Content:        annotation.Content,
		HitCount:       annotation.HitCount,
		CreatedAt:      annotation.CreatedAt,
		UpdatedAt:      annotation.UpdatedAt,
		Account:        account,
	}
}
