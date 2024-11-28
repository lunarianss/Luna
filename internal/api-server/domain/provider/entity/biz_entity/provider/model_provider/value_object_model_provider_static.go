package biz_entity

import (
	common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"
	"github.com/lunarianss/Luna/internal/pkg/field"
)

type ParameterType string

const (
	FLOAT   ParameterType = "float"
	INT     ParameterType = "int"
	STRING  ParameterType = "string"
	BOOLEAN ParameterType = "boolean"
	TEXT    ParameterType = "text"
)

type PriceConfig struct {
	Input    field.Float64 `json:"input" yaml:"input"`
	Output   field.Float64 `json:"output" yaml:"output"`
	Unit     field.Float64 `json:"unit" yaml:"unit"`
	Currency string        `json:"currency" yaml:"currency"`
}

type ParameterRule struct {
	Name        string             `json:"name" yaml:"name"`
	UseTemplate string             `json:"use_template" yaml:"use_template"`
	Label       *common.I18nObject `json:"label" yaml:"label"`
	Type        ParameterType      `json:"type" yaml:"type"`
	Help        *common.I18nObject `json:"help" yaml:"help"`
	Required    bool               `json:"required" yaml:"required"`
	Default     any                `json:"default" yaml:"default"`
	Min         float64            `json:"min" yaml:"min"`
	Max         float64            `json:"max" yaml:"max"`
	Precision   int                `json:"precision" yaml:"precision"`
	Options     []string           `json:"options" yaml:"options"`
}
