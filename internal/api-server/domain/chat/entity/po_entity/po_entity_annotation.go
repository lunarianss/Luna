package po_entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageAnnotation struct {
	ID             string `json:"id" gorm:"column:id"`
	AppID          string `json:"app_id" gorm:"column:app_id"`
	ConversationID string `json:"conversation_id" gorm:"column:conversation_id"`
	MessageID      string `json:"message_id" gorm:"column:message_id"`
	Question       string `json:"question" gorm:"column:question"`
	Content        string `json:"content" gorm:"column:content"`
	HitCount       int    `json:"hit_count" gorm:"column:hit_count"`
	AccountID      string `json:"account_id" gorm:"column:account_id"`
	CreatedAt      int64  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt      int64  `json:"updated_at" gorm:"column:updated_at"`
}

func (a *MessageAnnotation) TableName() string {
	return "message_annotations"
}

func (a *MessageAnnotation) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.NewString()
	return
}

type AppAnnotationSetting struct {
	ID                  string  `json:"id" gorm:"column:id"`
	AppID               string  `json:"app_id" gorm:"column:app_id"`
	ScoreThreshold      float64 `json:"score_threshold" gorm:"column:score_threshold"`
	CollectionBindingID string  `json:"collection_binding_id" gorm:"column:collection_binding_id"`
	CreatedUserID       string  `json:"created_user_id" gorm:"column:created_user_id"`
	CreatedAt           int64   `json:"created_at" gorm:"column:created_at"`
	UpdatedUserID       string  `json:"updated_user_id" gorm:"column:updated_user_id"`
	UpdatedAt           int64   `json:"updated_at" gorm:"column:updated_at"`
}

func (a *AppAnnotationSetting) TableName() string {
	return "app_annotation_settings"
}

func (a *AppAnnotationSetting) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.NewString()
	return
}
