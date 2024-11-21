// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model_provider

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/lunarianss/Luna/internal/api-server/entities/base"

	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/internal/pkg/field"
	"github.com/lunarianss/Luna/pkg/errors"
	"gopkg.in/yaml.v3"
)

type ParameterType string

const (
	FLOAT   ParameterType = "float"
	INT     ParameterType = "int"
	STRING  ParameterType = "string"
	BOOLEAN ParameterType = "boolean"
	TEXT    ParameterType = "text"
)

type ParameterRule struct {
	Name        string          `json:"name" yaml:"name"`
	UseTemplate string          `json:"use_template" yaml:"use_template"`
	Label       base.I18nObject `json:"label" yaml:"label"`
	Type        ParameterType   `json:"type" yaml:"type"`
	Help        base.I18nObject `json:"help" yaml:"help"`
	Required    bool            `json:"required" yaml:"required"`
	Default     any             `json:"default" yaml:"default"`
	Min         float64         `json:"min" yaml:"min"`
	Max         float64         `json:"max" yaml:"max"`
	Precision   int             `json:"precision" yaml:"precision"`
	Options     []string        `json:"options" yaml:"options"`
}

type PriceConfig struct {
	Input    field.Float64 `json:"input" yaml:"input"`
	Output   field.Float64 `json:"output" yaml:"output"`
	Unit     field.Float64 `json:"unit" yaml:"unit"`
	Currency string        `json:"currency" yaml:"currency"`
}

type AIModelEntity struct {
	*ProviderModel `yaml:",inline"`
	ParameterRules []*ParameterRule `json:"parameter_rules" yaml:"parameter_rules"`
	Pricing        *PriceConfig     `json:"pricing" yaml:"pricing"`
	Position       int              `json:"position" yaml:"position"`
}

type AIModel struct {
	ModelType     base.ModelType   `json:"model_type" yaml:"model_type"`
	ModelSchemas  []*AIModelEntity `json:"model_schemas" yaml:"model_schemas"`
	StartedAt     float64          `json:"started_at" yaml:"started_at"`
	ModelConfPath string           `json:"model_conf_path" yaml:"model_conf_path"`
}

func (a *AIModel) GetModelPositionMap() (map[string]int, error) {
	var positionMap = make(map[string]int)

	var modelPosition []string
	modelConfDir := a.ModelConfPath
	positionFilePath := fmt.Sprintf("%s/_position.yaml", modelConfDir)

	positionContext, err := os.ReadFile(positionFilePath)

	if os.IsNotExist(err) {
		return positionMap, nil
	}

	if err != nil {
		return nil, errors.WithCode(code.ErrRunTimeCaller, err.Error())
	}

	if err := yaml.Unmarshal(positionContext, &modelPosition); err != nil {
		return nil, errors.WithCode(code.ErrDecodingJSON, err.Error())
	}

	for i, v := range modelPosition {
		positionMap[v] = i + 1
	}

	return positionMap, nil
}

func (a *AIModel) PredefinedModels() ([]*AIModelEntity, error) {
	var (
		modelSchemaYamlPath []string
		AIModelEntities     []*AIModelEntity
	)

	modelConfDir := a.ModelConfPath

	dirEntries, err := os.ReadDir(modelConfDir)

	if os.IsNotExist(err) {
		return nil, nil
	}

	if err != nil {
		return nil, errors.WithCode(code.ErrRunTimeCaller, err.Error())
	}
	modelPosition, err := a.GetModelPositionMap()

	if err != nil {
		return nil, err
	}

	for _, dirEntry := range dirEntries {
		dirOrFileName := dirEntry.Name()
		if !dirEntry.IsDir() && !strings.HasPrefix(dirOrFileName, "_") && !strings.HasPrefix(dirOrFileName, "__") && strings.HasSuffix(dirOrFileName, ".yaml") {
			modelSchemaYamlPath = append(modelSchemaYamlPath, fmt.Sprintf("%s/%s", modelConfDir, dirOrFileName))
		}
	}

	for _, modelSchemaYamlPath := range modelSchemaYamlPath {
		AIModelEntity := &AIModelEntity{ProviderModel: &ProviderModel{}}
		AIModelEntity.FetchFrom = base.PREDEFINED_MODEL_FROM
		AIModelEntityContent, err := os.ReadFile(modelSchemaYamlPath)
		if err != nil {
			return nil, errors.WithCode(code.ErrRunTimeCaller, err.Error())
		}

		if err := yaml.Unmarshal(AIModelEntityContent, AIModelEntity); err != nil {
			return nil, errors.WithCode(code.ErrDecodingYaml, fmt.Sprintf("when decoding model %s of path,  failed: %s", modelSchemaYamlPath, err.Error()))
		}

		if v, ok := modelPosition[AIModelEntity.Model]; ok {
			AIModelEntity.Position = v
		} else {
			AIModelEntity.Position = 999
		}
		AIModelEntities = append(AIModelEntities, AIModelEntity)
	}

	sort.Slice(AIModelEntities, func(i, j int) bool {
		return AIModelEntities[i].Position < AIModelEntities[j].Position
	})

	return AIModelEntities, nil
}
