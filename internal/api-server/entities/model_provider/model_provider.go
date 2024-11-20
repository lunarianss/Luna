// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model_provider

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"gopkg.in/yaml.v3"

	"github.com/lunarianss/Luna/internal/api-server/entities/base"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
)

type IModelProviderRepo interface {
	ValidateProviderCredentials() error
}

type ModelProvider struct {
	ProviderSchema   ProviderEntity
	ModelConfPath    string
	ModelInstanceMap map[string]*AIModel
}

func (mp *ModelProvider) Models(modelType base.ModelType) ([]*AIModelEntity, error) {

	providerEntity, err := mp.GetProviderSchema()

	if !slices.Contains(providerEntity.SupportedModelTypes, modelType) {
		return nil, nil
	}

	AIModelInstance := mp.GetModelInstance(modelType)

	if err != nil {
		return nil, err
	}

	return AIModelInstance.PredefinedModels()
}

func (mp *ModelProvider) GetModelInstance(modelType base.ModelType) *AIModel {
	providerName := filepath.Base(mp.ModelConfPath)
	modelSchemaPath := fmt.Sprintf("%s/%s", mp.ModelConfPath, modelType)
	mp.ModelInstanceMap = make(map[string]*AIModel)

	if _, ok := mp.ModelInstanceMap[fmt.Sprintf("%s.%s", providerName, modelType)]; ok {
		return mp.ModelInstanceMap[fmt.Sprintf("%s.%s", providerName, modelType)]
	}

	AIModel := &AIModel{
		ModelType:     modelType,
		ModelConfPath: modelSchemaPath,
	}

	mp.ModelInstanceMap[fmt.Sprintf("%s.%s", providerName, modelType)] = AIModel

	return AIModel
}

func (mp *ModelProvider) GetProviderSchema() (*ProviderEntity, error) {
	providerName := filepath.Base(mp.ModelConfPath)
	providerSchemaPath := fmt.Sprintf("%s/%s.yaml", mp.ModelConfPath, providerName)
	providerContent, err := os.ReadFile(providerSchemaPath)

	if err != nil {
		return nil, errors.WithCode(code.ErrRunTimeCaller, err.Error())
	}

	provider := &ProviderEntity{}
	err = yaml.Unmarshal(providerContent, provider)

	if err != nil {
		return nil, errors.WithCode(code.ErrRunTimeCaller, err.Error())
	}

	provider.PatchI18nObject()
	return provider, nil
}

type DefaultModelProviderEntity struct {
	Provider            string           `json:"provider"`
	Label               *base.I18nObject `json:"label"`
	IconSmall           *base.I18nObject `json:"icon_small"`
	IconLarge           *base.I18nObject `json:"icon_large"`
	SupportedModelTypes []base.ModelType `json:"supported_model_types"`
}

type DefaultModelEntity struct {
	Model     string                      `json:"model"`
	ModelType string                      `json:"model_type"`
	Provider  *DefaultModelProviderEntity `json:"provider"`
}

type SimpleModelProviderEntity struct {
	Provider            string           `json:"provider"`
	Label               *base.I18nObject `json:"label"`
	IconSmall           *base.I18nObject `json:"icon_small"`
	IconLarge           *base.I18nObject `json:"icon_large"`
	SupportedModelTypes []base.ModelType `json:"supported_model_types"`
}

type ProviderModelWithStatusEntity struct {
	Status         ModelStatus `json:"status"`
	*ProviderModel `json:"provider_model"`
}

type ModelWithProviderEntity struct {
	*ProviderModelWithStatusEntity
	Provider *SimpleModelProviderEntity `json:"provider"`
}

type CustomConfigurationStatus string

const (
	Custom_ACTIVE       = "active"
	Custom_NO_CONFIGURE = "no-configure"
)
