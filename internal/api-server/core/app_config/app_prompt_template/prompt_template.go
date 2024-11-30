// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package app_prompt_template

import (
	"github.com/lunarianss/Luna/infrastructure/errors"
	biz_entity_app_config "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/provider_app_config"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
)

type promptTemplateConfigManager struct{}

func NewPromptTemplateConfigManager() *promptTemplateConfigManager {
	return &promptTemplateConfigManager{}
}

func (m *promptTemplateConfigManager) Convert(config *dto.AppModelConfigDto) (*biz_entity_app_config.PromptTemplateEntity, error) {
	promptType := config.PromptType

	if promptType == "" {
		return nil, errors.WithCode(code.ErrRequiredPromptType, "prompt_type is required")
	}

	if promptType == string(biz_entity_app_config.SIMPLE) {
		simplePromptMessage := config.PrePrompt
		return &biz_entity_app_config.PromptTemplateEntity{
			SimplePromptTemplate: simplePromptMessage,
			PromptType:           promptType,
		}, nil
	}

	return nil, nil
}
