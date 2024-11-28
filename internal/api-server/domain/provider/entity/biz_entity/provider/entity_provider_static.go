// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package biz_entity

import (
	"reflect"

	common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"
)

type ProviderStaticConfiguration struct {
	Provider                 string                    `json:"provider"                   yaml:"provider"`                   // Provider name
	Label                    *common.I18nObject        `json:"label"                      yaml:"label"`                      // Label in i18n format
	Description              *common.I18nObject        `json:"description"                yaml:"description"`                // Description in i18n format
	IconSmall                *common.I18nObject        `json:"icon_small"                 yaml:"icon_small"`                 // Small icon in i18n format
	IconLarge                *common.I18nObject        `json:"icon_large"                 yaml:"icon_large"`                 // Large icon in i18n format
	Background               string                    `json:"background"                 yaml:"background"`                 // Background color or image
	Help                     *ProviderHelpEntity       `json:"help"                       yaml:"help"`                       // Help information
	SupportedModelTypes      []common.ModelType        `json:"supported_model_types"      yaml:"supported_model_types"`      // Supported model types
	ConfigurationMethods     []ConfigurationMethod     `json:"configurate_methods"      yaml:"configurate_methods"`          // Configuration methods
	Models                   []*common.ProviderModel   `json:"models"                     yaml:"models"`                     // Models offered by the provider
	ProviderCredentialSchema *ProviderCredentialSchema `json:"provider_credential_schema" yaml:"provider_credential_schema"` // Schema for provider credentials
	ModelCredentialSchema    *ModelCredentialSchema    `json:"model_credential_schema"    yaml:"model_credential_schema"`    // Schema for model credentials
	Position                 int                       `json:"position" yaml:"position"`
}

type RecursiveObject func(obj interface{})

func (pe *ProviderStaticConfiguration) PatchI18nObject() {
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

					if field.Type() == reflect.TypeOf(&common.I18nObject{}) {
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
