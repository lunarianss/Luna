// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package entities

import (
	"github.com/lunarianss/Luna/internal/api-server/model-runtime/entities"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
)

type QuotaUnit string

const (
	TIMES   QuotaUnit = "times"
	TOKENS  QuotaUnit = "tokens"
	CREDITS QuotaUnit = "credits"
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

type CustomConfigurationStatus string

type CustomProviderConfiguration struct {
	Credentials interface{}
}

type CustomConfiguration struct {
	Provider CustomProviderConfiguration
	Models   []*CustomModelConfiguration
}

type CustomModelConfiguration struct {
	Model       string
	ModelType   string
	Credentials interface{}
}

type ModelSettings struct {
	Model     string
	ModelType entities.ModelType
	Enabled   bool
}
