// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

import (
	"fmt"

	"github.com/lunarianss/Luna/internal/api-server/config"
	"github.com/lunarianss/Luna/internal/api-server/entities/base"
	"github.com/lunarianss/Luna/internal/api-server/entities/model_provider"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
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
	Enabled             bool                               `json:"enabled"`
	CurrentQuotaType    model.ProviderQuotaType            `json:"current_quota_type"`
	QuotaConfigurations *model_provider.QuotaConfiguration `json:"quota_configurations"`
}

type ProviderResponse struct {
	Provider                 string                                   `json:"provider"`                   // Provider name
	Label                    *base.I18nObject                         `json:"label"`                      // Label in i18n format
	Description              *base.I18nObject                         `json:"description"`                // Description in i18n format
	IconSmall                *base.I18nObject                         `json:"icon_small"`                 // Small icon in i18n format
	IconLarge                *base.I18nObject                         `json:"icon_large"`                 // Large icon in i18n format
	Background               string                                   `json:"background"`                 // Background color or image
	Help                     *model_provider.ProviderHelpEntity       `json:"help"`                       // Help information
	SupportedModelTypes      []base.ModelType                         `json:"supported_model_types"`      // Supported model types
	ConfigurationMethods     []model_provider.ConfigurationMethod     `json:"configuration_methods"`      // Configuration methods                    // Models offered by the provider
	ProviderCredentialSchema *model_provider.ProviderCredentialSchema `json:"provider_credential_schema"` // Schema for provider credentials
	ModelCredentialSchema    *model_provider.ModelCredentialSchema    `json:"model_credential_schema"`    // Schema for model credentials
	PreferredProviderType    model.ProviderType                       `json:"preferred_provider_type"`    //
	CustomConfiguration      *CustomConfigurationResponse             `json:"custom_configuration"`
	SystemConfiguration      *SystemConfigurationResponse             `json:"system_configuration"`
	Position                 int                                      `json:"position"`
}

func (pr *ProviderResponse) PatchIcon() error {
	runtimeConfig, err := config.GetLunaRuntimeConfig()

	provider := pr.Provider

	if err != nil {
		return err
	}

	insecureAddress := fmt.Sprintf("%s:%d", runtimeConfig.InsecureServing.BindAddress, runtimeConfig.InsecureServing.BindPort)

	urlPrefix := fmt.Sprintf("http://%s/%s/%s", insecureAddress, "v1/console/workspace/current/model-providers", provider)

	if pr.IconLarge != nil {
		pr.IconLarge = &base.I18nObject{
			Zh_Hans: fmt.Sprintf("%s/%s", urlPrefix, "icon_large/zh_Hans"),
			En_US:   fmt.Sprintf("%s/%s", urlPrefix, "icon_large/en_US"),
		}
	} else if pr.IconSmall != nil {
		pr.IconSmall = &base.I18nObject{
			Zh_Hans: fmt.Sprintf("%s/%s", urlPrefix, "icon_small/zh_Hans"),
			En_US:   fmt.Sprintf("%s/%s", urlPrefix, "icon_small/en_US"),
		}
	}

	return nil
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
	ModelType string `uri:"modelType" validate:"required"`
}

type ProviderWithModelsResponse struct {
	Provider  string                                          `json:"provider"`
	Label     *base.I18nObject                                `json:"label"`
	IconSmall *base.I18nObject                                `json:"icon_small"`
	IconLarge *base.I18nObject                                `json:"icon_large"`
	Status    CustomConfigurationStatus                       `json:"status"`
	Models    []*model_provider.ProviderModelWithStatusEntity `json:"models"`
}
