// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package app_prompt_template

import (
	biz_entity_app_config "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_config"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
)

type PromptTemplateConfigManager struct{}

func (m *PromptTemplateConfigManager) Convert(config map[string]interface{}) (*biz_entity_app_config.PromptTemplateEntity, error) {
	promptType, ok := config["prompt_type"]

	if !ok {
		return nil, errors.WithCode(code.ErrRequiredPromptType, "prompt_type is required")
	}

	promptStr, ok := promptType.(string)
	if !ok {
		return nil, errors.WithCode(code.ErrRequiredPromptType, "prompt_type must be string")
	}

	if promptStr == string(biz_entity_app_config.SIMPLE) {
		simplePromptMessage, _ := config["pre_prompt"].(string)
		return &biz_entity_app_config.PromptTemplateEntity{
			SimplePromptTemplate: simplePromptMessage,
			PromptType:           promptStr,
		}, nil
	}

	return nil, nil
}
