package dto

import (
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
)

// Create App Input Dto
type CreateAppRequest struct {
	Name           string `json:"name" validate:"required"`
	Mode           string `json:"mode" validate:"required"`
	Icon           string `json:"icon" validate:"required"`
	Description    string `json:"description"`
	IconType       string `json:"icon_type"`
	IconBackground string `json:"icon_background"`
	ApiRph         int    `json:"api_rph"`
	ApiRpm         int    `json:"api_rpm"`
}

// Create App Response Dto
type CreateAppResponse struct {
	*model.App
	ModelConfig *model.AppModelConfig `json:"model_config"`
}

type ListAppRequest struct {
	Page     int `form:"page" validate:"required,min=1"`
	PageSize int `form:"limit" validate:"required,min=1,max=100"`
}
type ListAppsResponse struct {
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
	Total    int64          `json:"total"`
	Data     []*ListAppItem `json:"data"`
	HasMore  int            `json:"has_more"`
}

type ListAppItem struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Mode                string `json:"mode"`
	Icon                string `json:"icon"`
	IconBackground      string `json:"icon_background"`
	CreatedAt           int    `json:"created_at"`
	UpdatedAt           int    `json:"updated_at"`
	WorkflowID          string `json:"workflow_id"`
	Description         string `json:"description"`
	MaxActiveRequests   int    `json:"max_active_requests"`
	IconType            string `json:"icon_type"`
	CreatedBy           string `json:"created_by"`
	UpdatedBy           string `json:"updated_by"`
	UseIconAsAnswerIcon int    `json:"use_icon_as_answer_icon"`
}

func ListAppRecordToItem(app *model.App) *ListAppItem {
	return &ListAppItem{
		ID:                  app.ID,
		Name:                app.Name,
		Mode:                app.Mode,
		Icon:                app.Icon,
		IconBackground:      app.IconBackground,
		CreatedAt:           app.CreatedAt,
		UpdatedAt:           app.UpdatedAt,
		WorkflowID:          app.WorkflowID,
		Description:         app.Description,
		MaxActiveRequests:   app.MaxActiveRequests,
		IconType:            app.IconType,
		CreatedBy:           app.CreatedBy,
		UpdatedBy:           app.UpdatedBy,
		UseIconAsAnswerIcon: int(app.UseIconAsAnswerIcon),
	}

}

type AppDetailRequest struct {
	AppID string `uri:"appID" validate:"required"`
}

type AppDetail struct {
	ID                  string                 `json:"id"`
	Name                string                 `json:"name"`
	Description         string                 `json:"description"`
	Mode                string                 `json:"mode"`
	Icon                string                 `json:"icon"`
	IconType            string                 `json:"icon_type"`
	IconBackground      string                 `json:"icon_background"`
	EnableSite          int                    `json:"enable_site"`
	EnableApi           int                    `json:"enable_api"`
	ModelConfig         *model.AppModelConfig  `json:"model_config"`
	Workflow            map[string]interface{} `json:"workflow"`
	Site                map[string]interface{} `json:"site"`
	UseIconAsAnswerIcon int                    `json:"use_icon_as_answer_icon"`
	APIBaseUrl          string                 `json:"api_base_url"`
	CreatedAt           int                    `json:"created_at"`
	UpdatedAt           int                    `json:"updated_at"`
	CreatedBy           string                 `json:"created_by"`
	UpdatedBy           string                 `json:"updated_by"`
	DeletedTools        []interface{}          `json:"deleted_tools"`
}

func AppRecordToDetail(app *model.App, modelConfig *model.AppModelConfig) *AppDetail {
	return &AppDetail{
		ID:                  app.ID,
		Name:                app.Name,
		Mode:                app.Mode,
		Icon:                app.Icon,
		IconBackground:      app.IconBackground,
		CreatedAt:           app.CreatedAt,
		UpdatedAt:           app.UpdatedAt,
		Description:         app.Description,
		IconType:            app.IconType,
		CreatedBy:           app.CreatedBy,
		UpdatedBy:           app.UpdatedBy,
		UseIconAsAnswerIcon: int(app.UseIconAsAnswerIcon),
		ModelConfig:         modelConfig,
	}
}
