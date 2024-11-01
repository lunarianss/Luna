// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package entities

type ConfigurationMethod string

const (
	PREDEFINED_MODEL   ConfigurationMethod = "predefined-model"
	CUSTOMIZABLE_MODEL ConfigurationMethod = "customizable-model"
)

type FormType string

const (
	TEXT_INPUT   FormType = "text-input"
	SECRET_INPUT FormType = "secret-input"
	SELECT       FormType = "select"
	RADIO        FormType = "radio"
	SWITCH       FormType = "switch"
)

type I18nObject struct {
	Zh_Hans string `json:"zh_Hans" yaml:"zh_Hans"`
	En_US   string `json:"en_US"   yaml:"en_US"`
}

type ProviderHelpEntity struct {
	Title I18nObject `json:"title" yaml:"title"`
	Url   I18nObject `json:"url"   yaml:"url"`
}

type FormShowOnObject struct {
	Variable string `json:"variable"`
	Value    string `json:"value"`
}

type CredentialFormSchema struct {
	Variable     string             `json:"variable"   yaml:"variable"`   // Variable name
	Label        I18nObject         `json:"label"      yaml:"label"`      // Field label in i18n format
	Type         FormType           `json:"type"       yaml:"type"`       // Field type
	Required     bool               `json:"required"   yaml:"required"`   // Whether the field is required
	DefaultValue string             `json:"default"    yaml:"default"`    // Default value
	MaxLength    int                `json:"max_length" yaml:"max_length"` // Maximum length
	ShowOn       []FormShowOnObject `json:"show_on"    yaml:"show_on"`    // Conditions to show the field
}

type FieldModelSchema struct {
	Label       I18nObject `json:"label"`
	PlaceHolder I18nObject `json:"place_holder"`
}

type ModelCredentialSchema struct {
	Model                 FieldModelSchema        `json:"model"`
	CredentialFormSchemas []*CredentialFormSchema `json:"credential_form"`
}

type ProviderCredentialSchema struct {
	CredentialFormSchemas []*CredentialFormSchema `json:"credential_form_schemas" yaml:"credential_form_schemas"`
}

type ProviderEntity struct {
	Provider                 string                    `json:"provider"                   yaml:"provider"`                   // Provider name
	Label                    *I18nObject               `json:"label"                      yaml:"label"`                      // Label in i18n format
	Description              *I18nObject               `json:"description"                yaml:"description"`                // Description in i18n format
	IconSmall                *I18nObject               `json:"icon_small"                 yaml:"icon_small"`                 // Small icon in i18n format
	IconLarge                *I18nObject               `json:"icon_large"                 yaml:"icon_large"`                 // Large icon in i18n format
	Background               string                    `json:"background"                 yaml:"background"`                 // Background color or image
	Help                     *ProviderHelpEntity       `json:"help"                       yaml:"help"`                       // Help information
	SupportedModelTypes      []ModelType               `json:"supported_model_types"      yaml:"supported_model_types"`      // Supported model types
	ConfigurationMethods     []ConfigurationMethod     `json:"configuration_methods"      yaml:"configuration_methods"`      // Configuration methods
	Models                   []*ProviderModel          `json:"models"                     yaml:"models"`                     // Models offered by the provider
	ProviderCredentialSchema *ProviderCredentialSchema `json:"provider_credential_schema" yaml:"provider_credential_schema"` // Schema for provider credentials
	ModelCredentialSchema    *ModelCredentialSchema    `json:"model_credential_schema"    yaml:"model_credential_schema"`    // Schema for model credentials
	Position                 int                       `json:"position" yaml:"position"`
}
