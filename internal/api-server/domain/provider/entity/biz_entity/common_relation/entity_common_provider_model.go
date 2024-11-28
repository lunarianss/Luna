// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package common

type ProviderModel struct {
	Model           string                           `json:"model"            yaml:"model"`            // Model identifier
	Label           *I18nObject                      `json:"label"            yaml:"label"`            // Model label in i18n format
	ModelType       ModelType                        `json:"model_type"       yaml:"model_type"`       // Type of the model
	Features        []ModelFeature                   `json:"features"         yaml:"features"`         // List of model features
	FetchFrom       FetchFrom                        `json:"fetch_from"       yaml:"fetch_from"`       // Source from which to fetch the model
	ModelProperties map[ModelPropertyKey]interface{} `json:"model_properties" yaml:"model_properties"` // Properties of the model
	Deprecated      bool                             `json:"deprecated"       yaml:"deprecated"`       // Deprecation status
}


