// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model_provider

import (
	"reflect"

	"github.com/lunarianss/Luna/internal/api-server/entities/base"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
)

type QuotaUnit string

const (
	TIMES   QuotaUnit = "times"
	TOKENS  QuotaUnit = "tokens"
	CREDITS QuotaUnit = "credits"
)

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

type RestrictModels struct {
	Model         string
	BaseModelName string
	ModelType     string
}

type QuotaConfiguration struct {
	QuotaType      model.ProviderQuotaType
	QuotaUnit      QuotaUnit
	QuotaLimit     int
	QuotaUsed      int
	IsValid        bool
	RestrictModels []*RestrictModels
}

type SystemConfiguration struct {
	Enabled             bool
	CurrentQuotaType    model.ProviderQuotaType
	QuotaConfigurations []*QuotaConfiguration
	Credentials         interface{}
}

type CustomProviderConfiguration struct {
	Credentials interface{}
}

type CustomConfiguration struct {
	Provider *CustomProviderConfiguration
	Models   []*CustomModelConfiguration
}

type CustomModelConfiguration struct {
	Model       string
	ModelType   string
	Credentials map[string]interface{}
}

type ModelSettings struct {
	Model     string
	ModelType base.ModelType
	Enabled   bool
}

type ProviderHelpEntity struct {
	Title *base.I18nObject `json:"title" yaml:"title"`
	Url   *base.I18nObject `json:"url"   yaml:"url"`
}

type FormShowOnObject struct {
	Variable string `json:"variable"`
	Value    string `json:"value"`
}

type CredentialFormSchema struct {
	Variable     string              `json:"variable"   yaml:"variable"`   // Variable name
	Label        *base.I18nObject    `json:"label"      yaml:"label"`      // Field label in i18n format
	Type         FormType            `json:"type"       yaml:"type"`       // Field type
	Required     bool                `json:"required"   yaml:"required"`   // Whether the field is required
	DefaultValue string              `json:"default"    yaml:"default"`    // Default value
	MaxLength    int                 `json:"max_length" yaml:"max_length"` // Maximum length
	ShowOn       []*FormShowOnObject `json:"show_on"    yaml:"show_on"`    // Conditions to show the field
}

type FieldModelSchema struct {
	Label       *base.I18nObject `json:"label"`
	PlaceHolder *base.I18nObject `json:"place_holder"`
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
	Label                    *base.I18nObject          `json:"label"                      yaml:"label"`                      // Label in i18n format
	Description              *base.I18nObject          `json:"description"                yaml:"description"`                // Description in i18n format
	IconSmall                *base.I18nObject          `json:"icon_small"                 yaml:"icon_small"`                 // Small icon in i18n format
	IconLarge                *base.I18nObject          `json:"icon_large"                 yaml:"icon_large"`                 // Large icon in i18n format
	Background               string                    `json:"background"                 yaml:"background"`                 // Background color or image
	Help                     *ProviderHelpEntity       `json:"help"                       yaml:"help"`                       // Help information
	SupportedModelTypes      []base.ModelType          `json:"supported_model_types"      yaml:"supported_model_types"`      // Supported model types
	ConfigurationMethods     []ConfigurationMethod     `json:"configuration_methods"      yaml:"configuration_methods"`      // Configuration methods
	Models                   []ProviderModel           `json:"models"                     yaml:"models"`                     // Models offered by the provider
	ProviderCredentialSchema *ProviderCredentialSchema `json:"provider_credential_schema" yaml:"provider_credential_schema"` // Schema for provider credentials
	ModelCredentialSchema    *ModelCredentialSchema    `json:"model_credential_schema"    yaml:"model_credential_schema"`    // Schema for model credentials
	Position                 int                       `json:"position" yaml:"position"`
}

type RecursiveObject func(obj interface{})

func (pe *ProviderEntity) PatchI18nObject() {
	var recursiveObject RecursiveObject
	recursiveObject = func(obj interface{}) {

		if obj == nil {
			return
		}

		var v reflect.Value

		objValue := reflect.ValueOf(obj)

		if objValue.Kind() == reflect.Ptr {
			v = objValue.Elem()
		} else {
			v = objValue
		}

		if v.Kind() == reflect.Struct {
			for i := 0; i < v.NumField(); i++ {
				field := v.Field(i)

				if field.Kind() == reflect.Ptr {
					if field.IsNil() {
						continue
					}

					if field.Type() == reflect.TypeOf(&base.I18nObject{}) {
						method := field.MethodByName(PATCH_FUNCTION_NAME)
						if method.IsValid() && method.Type().NumIn() == 0 {
							method.Call(nil)
						}
					} else {
						recursiveObject(field.Interface())
					}
				} else if field.Kind() == reflect.Slice {
					for j := 0; j < field.Len(); j++ {
						recursiveObject(field.Index(j).Interface())
					}
				} else if field.Kind() == reflect.Interface {
					if !field.IsNil() {
						recursiveObject(field.Interface())
					}
				}
			}
		}
	}

	recursiveObject(pe)
}
