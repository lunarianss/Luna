// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model_provider

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/lunarianss/Luna/internal/api-server/entities/base"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
)

type IModelProviderRepo interface {
	ValidateProviderCredentials() error
}

type ModelProvider struct {
	ProviderSchema ProviderEntity
	ModelConfPath  string
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
	Model     string                     `json:"model"`
	ModelType string                     `json:"model_type"`
	Provider  DefaultModelProviderEntity `json:"provider"`
}

type SimpleModelProviderEntity struct {
}

type ProviderModelWithStatusEntity struct {
	Status ModelStatus `json:"status"`
	*ProviderEntity
}

type ModelWithProviderEntity struct {
	*ProviderModelWithStatusEntity
	Provider *SimpleModelProviderEntity
}
