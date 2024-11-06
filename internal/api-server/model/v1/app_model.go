package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// App represents the app table in the database
type App struct {
	ID                  string `gorm:"column:id" json:"id"`
	TenantID            string `gorm:"column:tenant_id" json:"tenant_id"`
	Name                string `gorm:"column:name" json:"name"`
	Mode                string `gorm:"column:mode" json:"mode"`
	Icon                string `gorm:"column:icon" json:"icon"`
	IconBackground      string `gorm:"column:icon_background" json:"icon_background"`
	AppModelConfigID    string `gorm:"column:app_model_config_id" json:"app_model_config_id"`
	Status              string `gorm:"column:status" json:"status"`
	EnableSite          bool   `gorm:"column:enable_site" json:"enable_site"`
	EnableAPI           bool   `gorm:"column:enable_api" json:"enable_api"`
	APIRpm              int    `gorm:"column:api_rpm" json:"api_rpm"`
	APIRph              int    `gorm:"column:api_rph" json:"api_rph"`
	IsDemo              bool   `gorm:"column:is_demo" json:"is_demo"`
	IsPublic            bool   `gorm:"column:is_public" json:"is_public"`
	CreatedAt           int    `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           int    `gorm:"column:updated_at" json:"updated_at"`
	IsUniversal         bool   `gorm:"column:is_universal" json:"is_universal"`
	WorkflowID          string `gorm:"column:workflow_id" json:"workflow_id"`
	Description         string `gorm:"column:description" json:"description"`
	Tracing             string `gorm:"column:tracing" json:"tracing"`
	MaxActiveRequests   int    `gorm:"column:max_active_requests" json:"max_active_requests"`
	IconType            string `gorm:"column:icon_type" json:"icon_type"`
	CreatedBy           string `gorm:"column:created_by" json:"created_by"`
	UpdatedBy           string `gorm:"column:updated_by" json:"updated_by"`
	UseIconAsAnswerIcon bool   `gorm:"column:use_icon_as_answer_icon" json:"use_icon_as_answer_icon"`
}

func (a *App) TableName() string {
	return "apps"
}

func (a *App) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.NewString()
	return
}

type AppMode string

const (
	COMPLETION    AppMode = "completion"
	WORKFLOW      AppMode = "workflow"
	CHAT          AppMode = "chat"
	ADVANCED_CHAT AppMode = "advanced-chat"
	AGENT_CHAT    AppMode = "agent-chat"
	CHANNEL       AppMode = "channel"
)
