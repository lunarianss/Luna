// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

import (
	providerEntities "github.com/lunarianss/Luna/internal/api-server/entities/provider"
	"github.com/lunarianss/Luna/internal/api-server/model-runtime/entities"
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
	Enabled             bool                                 `json:"enabled"`
	CurrentQuotaType    model.ProviderQuotaType              `json:"current_quota_type"`
	QuotaConfigurations *providerEntities.QuotaConfiguration `json:"quota_configurations"`
}

type ProviderResponse struct {
	Provider                 string                             `json:"provider"`                   // Provider name
	Label                    *entities.I18nObject               `json:"label"`                      // Label in i18n format
	Description              *entities.I18nObject               `json:"description"`                // Description in i18n format
	IconSmall                *entities.I18nObject               `json:"icon_small"`                 // Small icon in i18n format
	IconLarge                *entities.I18nObject               `json:"icon_large"`                 // Large icon in i18n format
	Background               string                             `json:"background"`                 // Background color or image
	Help                     *entities.ProviderHelpEntity       `json:"help"`                       // Help information
	SupportedModelTypes      []entities.ModelType               `json:"supported_model_types"`      // Supported model types
	ConfigurationMethods     []entities.ConfigurationMethod     `json:"configuration_methods"`      // Configuration methods                    // Models offered by the provider
	ProviderCredentialSchema *entities.ProviderCredentialSchema `json:"provider_credential_schema"` // Schema for provider credentials
	ModelCredentialSchema    *entities.ModelCredentialSchema    `json:"model_credential_schema"`    // Schema for model credentials
	PreferredProviderType    model.ProviderType                 `json:"preferred_provider_type"`    //
	CustomConfiguration      *CustomConfigurationResponse       `json:"custom_configuration"`
	SystemConfiguration      *SystemConfigurationResponse       `json:"system_configuration"`
	Position                 int                                `json:"position"`
}

// --
// --- List icon
// --
type ListIconRequest struct {
	IconType string `json:"icon_type" uri:"iconType" validate:"required"`
	Lang     string `json:"lang" uri:"lang" validate:"required"`
	Provider string `json:"provider" uri:"provider" validate:"required"`
}
