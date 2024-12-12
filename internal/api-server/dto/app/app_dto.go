// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

import (
	"github.com/lunarianss/Luna/internal/api-server/config"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	po_entity_web_app "github.com/lunarianss/Luna/internal/api-server/domain/web_app/entity/po_entity"
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
	*po_entity.App
	ModelConfig *po_entity.AppModelConfig `json:"model_config"`
}

type ListAppRequest struct {
	Page     int `form:"page" validate:"required,min=1"`
	PageSize int `form:"limit" validate:"required,min=1,max=100"`
}
type ListAppsResponse struct {
	Page     int            `json:"page"`
	PageSize int            `json:"limit"`
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
	Tags                []any  `json:"tags"`
}

func ListAppRecordToItem(app *po_entity.App) *ListAppItem {
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
		Tags:                make([]any, 0),
	}
}

type AppDetailRequest struct {
	AppID string `uri:"appID" validate:"required"`
}

type AppDetail struct {
	ID                  string                   `json:"id"`
	Name                string                   `json:"name"`
	Description         string                   `json:"description"`
	Mode                string                   `json:"mode"`
	Icon                string                   `json:"icon"`
	IconType            string                   `json:"icon_type"`
	IconBackground      string                   `json:"icon_background"`
	EnableSite          int                      `json:"enable_site"`
	EnableApi           int                      `json:"enable_api"`
	ModelConfig         po_entity.AppModelConfig `json:"model_config"`
	Workflow            map[string]interface{}   `json:"workflow"`
	UseIconAsAnswerIcon int                      `json:"use_icon_as_answer_icon"`
	APIBaseUrl          string                   `json:"api_base_url"`
	CreatedAt           int                      `json:"created_at"`
	UpdatedAt           int                      `json:"updated_at"`
	CreatedBy           string                   `json:"created_by"`
	UpdatedBy           string                   `json:"updated_by"`
	DeletedTools        []interface{}            `json:"deleted_tools"`
	SiteDetail          *SiteDetail              `json:"site"`
}

func AppRecordToDetail(app *po_entity.App, config *config.Config, modelConfig *po_entity.AppModelConfig, siteRecord *po_entity_web_app.Site) *AppDetail {

	appDetail := &AppDetail{
		ID:                  app.ID,
		Name:                app.Name,
		Mode:                app.Mode,
		Icon:                app.Icon,
		IconBackground:      app.IconBackground,
		CreatedAt:           app.CreatedAt,
		UpdatedAt:           app.UpdatedAt,
		Description:         app.Description,
		EnableSite:          int(app.EnableSite),
		EnableApi:           int(app.EnableAPI),
		IconType:            app.IconType,
		CreatedBy:           app.CreatedBy,
		UpdatedBy:           app.UpdatedBy,
		UseIconAsAnswerIcon: int(app.UseIconAsAnswerIcon),
		ModelConfig:         *modelConfig,
		SiteDetail:          SiteRecordToSiteDetail(siteRecord, config),
		APIBaseUrl:          config.SystemOptions.ApiBaseUrl,
	}

	if !appDetail.ModelConfig.RetrieverResource.Enabled {
		appDetail.ModelConfig.RetrieverResource.Enabled = true
	}

	if appDetail.ModelConfig.SensitiveWordAvoidance == nil {
		appDetail.ModelConfig.SensitiveWordAvoidance = map[string]any{
			"enabled": false,
			"type":    "",
			"configs": []any{},
		}
	}

	if appDetail.ModelConfig.AgentMode == nil {
		appDetail.ModelConfig.AgentMode = map[string]any{
			"enabled":  false,
			"strategy": nil,
			"tools":    []any{},
			"prompt":   nil,
		}
	}

	if appDetail.ModelConfig.ChatPromptConfig == nil {
		appDetail.ModelConfig.ChatPromptConfig = map[string]any{}
	}

	if appDetail.ModelConfig.CompletionPromptConfig == nil {
		appDetail.ModelConfig.CompletionPromptConfig = map[string]any{}
	}

	if appDetail.ModelConfig.ExternalDataTools == nil {
		appDetail.ModelConfig.ExternalDataTools = []string{}
	}

	if appDetail.ModelConfig.DatasetConfigs == nil {
		appDetail.ModelConfig.DatasetConfigs = map[string]any{
			"retrieval_model": "multiple",
		}
	}

	if appDetail.ModelConfig.FileUpload == nil {
		appDetail.ModelConfig.FileUpload = map[string]interface{}{
			"image": map[string]interface{}{
				"enabled":          false,
				"number_limits":    3,
				"detail":           "high",
				"transfer_methods": []string{"remote_url", "local_file"},
			},
		}
	}

	if appDetail.ModelConfig.UserInputForm == nil {
		appDetail.ModelConfig.UserInputForm = make([]*po_entity.UserInputForm, 0)
	}

	return appDetail
}

type SiteDetail struct {
	*po_entity_web_app.Site
	AccessToken string `json:"access_token"`
	AppBaseUrl  string `json:"app_base_url"`
}

func SiteRecordToSiteDetail(sm *po_entity_web_app.Site, config *config.Config) *SiteDetail {
	return &SiteDetail{
		Site:        sm,
		AppBaseUrl:  config.SystemOptions.AppWebUrl,
		AccessToken: sm.Code,
	}
}

type EnableSiteRequest struct {
	EnableSite bool `json:"enable_site"`
}

type APIUrlParameter struct {
	AppID string `json:"appID" uri:"appID"`
}
