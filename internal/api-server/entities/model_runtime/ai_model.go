// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model_runtime

import "github.com/lunarianss/Luna/internal/api-server/entities/base"

type ParameterType string

const (
	FLOAT   ParameterType = "float"
	INT     ParameterType = "int"
	STRING  ParameterType = "string"
	BOOLEAN ParameterType = "boolean"
	TEXT    ParameterType = "text"
)

type ParameterRule struct {
	Name        string           `json:"name"`
	UseTemplate string           `json:"use_template"`
	Label       *base.I18nObject `json:"label"`
	Type        ParameterType    `json:"type"`
	Help        *base.I18nObject `json:"help"`
	Required    bool             `json:"required"`
	Default     any              `json:"default"`
	Min         float64          `json:"min"`
	Max         float64          `json:"max"`
	Precision   int              `json:"precision"`
	Options     string           `json:"options"`
}

type PriceConfig struct {
	Input    float64 `json:"input"`
	Output   float64 `json:"output"`
	Unit     float64 `json:"unit"`
	Currency string  `json:"currency"`
}

type AIModelEntity struct {
	ParameterRules []*ParameterRule `json:"parameter_rules"`
	Pricing        PriceConfig      `json:"pricing"`
}

type AIModel struct {
	ModelType    base.ModelType  `json:"model_type"`
	ModelSchemas []AIModelEntity `json:"model_schemas"`
	StartedAt    float64         `json:"started_at"`
}
