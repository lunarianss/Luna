package prompt_template

import (
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
)

type PromptTemplateConfigManager struct{}

func (m *PromptTemplateConfigManager) Convert(config map[string]interface{}) (*app_config.PromptTemplateEntity, error) {
	promptType, ok := config["prompt_type"]

	if !ok {
		return nil, errors.WithCode(code.ErrRequiredPromptType, "prompt_type is required")
	}

	promptStr, ok := promptType.(string)
	if !ok {
		return nil, errors.WithCode(code.ErrRequiredPromptType, "prompt_type must be string")
	}

	if promptStr == string(app_config.SIMPLE) {
		simplePromptMessage, _ := config["pre_prompt"].(string)
		return &app_config.PromptTemplateEntity{
			SimplePromptTemplate: simplePromptMessage,
			PromptType:           promptStr,
		}, nil
	}

	return nil, nil
}
