package biz_entity

import common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"

const (
	PATCH_FUNCTION_NAME = "PatchZh"
)

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

type ProviderHelpEntity struct {
	Title *common.I18nObject `json:"title" yaml:"title"`
	Url   *common.I18nObject `json:"url"   yaml:"url"`
}

type FormShowOnObject struct {
	Variable string `json:"variable" yaml:"variable"`
	Value    string `json:"value" yaml:"value"`
}

type FormOptions struct {
	Label  *common.I18nObject  `json:"label" yaml:"label"`
	Value  string              `json:"value" yaml:"value"`
	ShowOn []*FormShowOnObject `json:"show_on" yaml:"show_on"`
}

type CredentialFormSchema struct {
	Variable     string              `json:"variable"   yaml:"variable"`   // Variable name
	Label        *common.I18nObject  `json:"label"      yaml:"label"`      // Field label in i18n format
	Type         FormType            `json:"type"       yaml:"type"`       // Field type
	Required     bool                `json:"required"   yaml:"required"`   // Whether the field is required
	DefaultValue string              `json:"default"    yaml:"default"`    // Default value
	MaxLength    int                 `json:"max_length" yaml:"max_length"` // Maximum length
	ShowOn       []*FormShowOnObject `json:"show_on"    yaml:"show_on"`    // Conditions to show the field
	Options      []*FormOptions      `json:"options" yaml:"options"`
}

type FieldModelSchema struct {
	Label       *common.I18nObject `json:"label" yaml:"label"`
	PlaceHolder *common.I18nObject `json:"place_holder" yaml:"place_holder"`
}

type ModelCredentialSchema struct {
	Model                 FieldModelSchema        `json:"model" yaml:"model"`
	CredentialFormSchemas []*CredentialFormSchema `json:"credential_form_schemas" yaml:"credential_form_schemas"`
}

type ProviderCredentialSchema struct {
	CredentialFormSchemas []*CredentialFormSchema `json:"credential_form_schemas" yaml:"credential_form_schemas"`
}
