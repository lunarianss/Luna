// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package biz_entity

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/lunarianss/Luna/infrastructure/errors"
	common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/shopspring/decimal"
	"gopkg.in/yaml.v3"
)

type IAIModelRuntime interface {
	GetPrice(model string, credentials any, priceType PriceType, tokens int64) (*PriceInfo, error)
	GetModelSchema(modelName string, credentials any) (*AIModelStaticConfiguration, error)
	PredefinedModels() ([]*AIModelStaticConfiguration, error)
	GetModelPositionMap() (map[string]int, error)
}

type AIModelRuntime struct {
	ModelType     common.ModelType              `json:"model_type" yaml:"model_type"`
	ModelSchemas  []*AIModelStaticConfiguration `json:"model_schemas" yaml:"model_schemas"`
	StartedAt     float64                       `json:"started_at" yaml:"started_at"`
	ModelConfPath string                        `json:"model_conf_path" yaml:"model_conf_path"`
}

func (a *AIModelRuntime) GetModelPositionMap() (map[string]int, error) {
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

func (a *AIModelRuntime) PredefinedModels() ([]*AIModelStaticConfiguration, error) {
	var (
		modelSchemaYamlPath []string
		AIModelEntities     []*AIModelStaticConfiguration
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
		AIModelEntity := &AIModelStaticConfiguration{ProviderModel: &common.ProviderModel{}}
		AIModelEntity.FetchFrom = common.PREDEFINED_MODEL_FROM
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

		for _, parameterRule := range AIModelEntity.ParameterRules {
			if parameterRule.UseTemplate != "" {
				prt, ok := ParameterRuleTemplates[parameterRule.UseTemplate]

				if ok {
					parameterRule.Type = ParameterType(prt.Type)
					parameterRule.Label = prt.Label
					parameterRule.Help = prt.Help
					parameterRule.Required = prt.Required
					parameterRule.Min = prt.Min
					parameterRule.Max = prt.Max
					parameterRule.Default = prt.Default
					parameterRule.Precision = prt.Precision
					parameterRule.Options = prt.Options
				}

			}
		}

		AIModelEntities = append(AIModelEntities, AIModelEntity)
	}

	sort.Slice(AIModelEntities, func(i, j int) bool {
		return AIModelEntities[i].Position < AIModelEntities[j].Position
	})

	return AIModelEntities, nil
}

func (a *AIModelRuntime) GetModelSchema(modelName string, credentials any) (*AIModelStaticConfiguration, error) {

	models, err := a.PredefinedModels()

	if err != nil {
		return nil, err
	}

	for _, model := range models {
		if model.Model == modelName {
			return model, nil
		}
	}

	return nil, errors.WithCode(code.ErrModelSchemaNotFound, fmt.Sprintf("model schema %s is not found", modelName))
}

func (a *AIModelRuntime) GetPrice(model string, credentials any, priceType PriceType, tokens int64) (*PriceInfo, error) {
	var (
		priceConfig *PriceConfig
		unitPrice   float64
	)
	modelSchema, err := a.GetModelSchema(model, credentials)

	if err != nil {
		return nil, err
	}

	if modelSchema.Pricing == nil {
		return NewFreePriceInfo(), errors.WithCode(code.ErrModelNotHavePrice, fmt.Sprintf("model %s not have price info", model))
	}

	priceConfig = modelSchema.Pricing

	if priceType == INPUT {
		unitPrice = float64(priceConfig.Input)
	} else if priceType == OUTPUT {
		unitPrice = float64(priceConfig.Output)
	}

	if unitPrice == 0 {
		return NewFreePriceInfo(), nil
	}

	tokens_decimal := decimal.NewFromInt(tokens)
	unitPrice_decimal := decimal.NewFromFloat(unitPrice)
	unit_decimal := decimal.NewFromFloat(float64(priceConfig.Unit))

	totalAmount := tokens_decimal.Mul(unitPrice_decimal).Mul(unit_decimal)
	totalAmount = totalAmount.Round(7)

	return &PriceInfo{
		TotalAmount: totalAmount.InexactFloat64(),
		UnitPrice:   unitPrice,
		Unit:        float64(priceConfig.Unit),
		Currency:    priceConfig.Currency,
	}, nil
}
