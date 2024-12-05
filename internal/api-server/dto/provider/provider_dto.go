// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

import (
	"fmt"

	"github.com/lunarianss/Luna/internal/api-server/config"
	common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider"
	biz_entity_provider_config "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_configuration"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/po_entity"
)

// --
// --- List model providers
// --
type CustomConfigurationStatus string

const (
	ACTIVE       CustomConfigurationStatus = "active"
	NO_CONFIGURE CustomConfigurationStatus = "no-configure"
)

type CustomConfigurationResponse struct {
	Status CustomConfigurationStatus `json:"status"`
}

type SystemConfigurationResponse struct {
	Enabled             bool                                             `json:"enabled"`
	CurrentQuotaType    po_entity.ProviderQuotaType                      `json:"current_quota_type"`
	QuotaConfigurations []*biz_entity_provider_config.QuotaConfiguration `json:"quota_configurations"`
}

type ProviderResponse struct {
	Provider                 string                               `json:"provider"`                   // Provider name
	Label                    *common.I18nObject                   `json:"label"`                      // Label in i18n format
	Description              *common.I18nObject                   `json:"description"`                // Description in i18n format
	IconSmall                *common.I18nObject                   `json:"icon_small"`                 // Small icon in i18n format
	IconLarge                *common.I18nObject                   `json:"icon_large"`                 // Large icon in i18n format
	Background               string                               `json:"background"`                 // Background color or image
	Help                     *biz_entity.ProviderHelpEntity       `json:"help"`                       // Help information
	SupportedModelTypes      []common.ModelType                   `json:"supported_model_types"`      // Supported model types
	ConfigurationMethods     []biz_entity.ConfigurationMethod     `json:"configurate_methods"`        // Configuration methods                    // Models offered by the provider
	ProviderCredentialSchema *biz_entity.ProviderCredentialSchema `json:"provider_credential_schema"` // Schema for provider credentials
	ModelCredentialSchema    *biz_entity.ModelCredentialSchema    `json:"model_credential_schema"`    // Schema for model credentials
	PreferredProviderType    po_entity.ProviderType               `json:"preferred_provider_type"`    //
	CustomConfiguration      *CustomConfigurationResponse         `json:"custom_configuration"`
	SystemConfiguration      *SystemConfigurationResponse         `json:"system_configuration"`
	Position                 int                                  `json:"position"`
}

func (pr *ProviderResponse) PatchIcon(runtimeConfig *config.Config) {

	provider := pr.Provider

	urlPrefix := fmt.Sprintf("%s/%s/%s", runtimeConfig.SystemOptions.IconBaseUrl, "v1/console/api/workspaces/current/model-providers", provider)

	if pr.IconLarge != nil {
		pr.IconLarge = &common.I18nObject{
			Zh_Hans: fmt.Sprintf("%s/%s", urlPrefix, "icon_large/zh_Hans"),
			En_US:   fmt.Sprintf("%s/%s", urlPrefix, "icon_large/en_US"),
		}
	}

	if pr.IconSmall != nil {
		pr.IconSmall = &common.I18nObject{
			Zh_Hans: fmt.Sprintf("%s/%s", urlPrefix, "icon_small/zh_Hans"),
			En_US:   fmt.Sprintf("%s/%s", urlPrefix, "icon_small/en_US"),
		}
	}

}

// --
// --- List icon
// --
type ListIconRequest struct {
	IconType string `json:"icon_type" uri:"iconType" validate:"required"`
	Lang     string `json:"lang" uri:"lang" validate:"required"`
	Provider string `json:"provider" uri:"provider" validate:"required"`
}

// --
// --- Create provider credentials
// --

type CreateProviderCredentialUri struct {
	Provider string `uri:"provider"  validate:"required"`
}

type CreateProviderCredentialBody struct {
	ConfigFrom  string                 `json:"config_from"  validate:"required"`
	Credentials map[string]interface{} `json:"credentials"  validate:"required"`
}

// --
// --- Create  model credentials
// --
type CreateModelCredentialUri struct {
	Provider string `uri:"provider"  validate:"required"`
}

type CreateModelCredentialBody struct {
	Model       string                 `json:"model"  validate:"required"`
	ModelType   string                 `json:"model_type"  validate:"required"`
	Credentials map[string]interface{} `json:"credentials"  validate:"required"`
}

type GetAccountAvailableModelsRequest struct {
	ModelType string `uri:"modelType" validate:"required,valid_model_type" json:"model_type"`
}

type ProviderWithModelsResponse struct {
	Provider  string                                                `json:"provider"`
	Label     *common.I18nObject                                    `json:"label"`
	IconSmall *common.I18nObject                                    `json:"icon_small"`
	IconLarge *common.I18nObject                                    `json:"icon_large"`
	Status    CustomConfigurationStatus                             `json:"status"`
	Models    []*biz_entity_provider_config.ProviderModelWithStatus `json:"models"`
}

func (pr *ProviderWithModelsResponse) PatchIcon(runtimeConfig *config.Config) {

	provider := pr.Provider

	urlPrefix := fmt.Sprintf("%s/%s/%s", runtimeConfig.SystemOptions.IconBaseUrl, "v1/console/api/workspaces/current/model-providers", provider)
	if pr.IconLarge != nil {
		pr.IconLarge = &common.I18nObject{
			Zh_Hans: fmt.Sprintf("%s/%s", urlPrefix, "icon_large/zh_Hans"),
			En_US:   fmt.Sprintf("%s/%s", urlPrefix, "icon_large/en_US"),
		}
	}

	if pr.IconSmall != nil {
		pr.IconSmall = &common.I18nObject{
			Zh_Hans: fmt.Sprintf("%s/%s", urlPrefix, "icon_small/zh_Hans"),
			En_US:   fmt.Sprintf("%s/%s", urlPrefix, "icon_small/en_US"),
		}
	}
}

type ParameterRulesQuery struct {
	Model string `form:"model" validate:"required"`
}
type DefaultModelByTypeQuery struct {
	ModelType string `form:"model_type" validate:"required"`
}

// DefaultModelResponse represents the default model entity.
type DefaultModelResponse struct {
	Model     string                                          `json:"model"`
	ModelType string                                          `json:"model_type"`
	Provider  *biz_entity_provider_config.SimpleModelProvider `json:"provider"`
}
