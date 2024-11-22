package model

import (
	"github.com/google/uuid"
	"github.com/lunarianss/Luna/internal/pkg/field"
	"gorm.io/gorm"
)

type EndUser struct {
	ID             string        `json:"id" gorm:"column:id"`
	TenantID       string        `json:"tenant_id" gorm:"column:tenant_id"`
	AppID          string        `json:"app_id" gorm:"column:app_id"`
	Type           string        `json:"type" gorm:"column:type"`
	ExternalUserID string        `json:"external_user_id" gorm:"column:external_user_id"`
	Name           string        `json:"name" gorm:"column:name"`
	IsAnonymous    field.BitBool `json:"is_anonymous" gorm:"column:is_anonymous"`
	SessionID      string        `json:"session_id" gorm:"column:session_id"`
	CreatedAt      int64         `json:"created_at" gorm:"column:created_at"`
	UpdatedAt      int64         `json:"updated_at" gorm:"column:updated_at"`
}

func (a *EndUser) TableName() string {
	return "end_users"
}

func (a *EndUser) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.NewString()
	return
}

type InstalledApp struct {
	ID               string        `json:"id" gorm:"column:id"`
	TenantID         string        `json:"tenant_id" gorm:"column:tenant_id"`
	AppID            string        `json:"app_id" gorm:"column:app_id"`
	AppOwnerTenantID string        `json:"app_owner_tenant_id" gorm:"column:app_owner_tenant_id"`
	Position         int           `json:"position" gorm:"column:position"`
	IsPinned         field.BitBool `json:"is_pinned" gorm:"column:is_pinned"`
	LastUsedAt       int64         `json:"last_used_at" gorm:"column:last_used_at"`
	CreatedAt        int64         `json:"created_at" gorm:"column:created_at"`
}

func (a *InstalledApp) TableName() string {
	return "installed_apps"
}

func (a *InstalledApp) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.NewString()
	return
}

type Site struct {
	ID                     string        `json:"id" gorm:"column:id"`
	AppID                  string        `json:"app_id" gorm:"column:app_id"`
	Title                  string        `json:"title" gorm:"column:title"`
	IconType               string        `json:"icon_type" gorm:"column:icon_type"`
	Icon                   string        `json:"icon" gorm:"column:icon"`
	IconBackground         string        `json:"icon_background" gorm:"column:icon_background"`
	Description            string        `json:"description" gorm:"column:description"`
	DefaultLanguage        string        `json:"default_language" gorm:"column:default_language"`
	ChatColorTheme         string        `json:"chat_color_theme" gorm:"column:chat_color_theme"`
	ChatColorThemeInverted field.BitBool `json:"chat_color_theme_inverted" gorm:"column:chat_color_theme_inverted"`
	Copyright              string        `json:"copyright" gorm:"column:copyright"`
	PrivacyPolicy          string        `json:"privacy_policy" gorm:"column:privacy_policy"`
	ShowWorkflowSteps      field.BitBool `json:"show_workflow_steps" gorm:"column:show_workflow_steps"`
	UseIconAsAnswerIcon    field.BitBool `json:"use_icon_as_answer_icon" gorm:"column:use_icon_as_answer_icon"`
	CustomDisclaimer       string        `json:"custom_disclaimer" gorm:"column:custom_disclaimer"`
	CustomizeDomain        string        `json:"customize_domain" gorm:"column:customize_domain"`
	CustomizeTokenStrategy string        `json:"customize_token_strategy" gorm:"column:customize_token_strategy"`
	PromptPublic           field.BitBool `json:"prompt_public" gorm:"column:prompt_public"`
	Status                 string        `json:"status" gorm:"column:status"`
	CreatedBy              string        `json:"created_by" gorm:"column:created_by"`
	CreatedAt              int64         `json:"created_at" gorm:"column:created_at"`
	UpdatedBy              string        `json:"updated_by" gorm:"column:updated_by"`
	UpdatedAt              int64         `json:"updated_at" gorm:"column:updated_at"`
	Code                   string        `json:"code" gorm:"column:code"`
}

func (a *Site) TableName() string {
	return "sites"
}

func (a *Site) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.NewString()
	return
}
