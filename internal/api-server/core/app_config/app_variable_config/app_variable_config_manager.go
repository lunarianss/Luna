package app_variable_config

import (
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/provider_app_config"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
)

type BasicVariablesConfigManager struct{}

func NewBasicVariablesConfigManager() *BasicVariablesConfigManager {
	return &BasicVariablesConfigManager{}
}

func (*BasicVariablesConfigManager) Convert(appModelConfig *dto.AppModelConfigDto) []*biz_entity.VariableEntity {

	var variables []*biz_entity.VariableEntity
	userInputForms := appModelConfig.UserInputForm

	for _, userInputForm := range userInputForms {
		if userInputForm.TextInput != nil {
			variables = append(variables, &biz_entity.VariableEntity{
				Type:      "text-input",
				Variable:  userInputForm.TextInput.Variable,
				Label:     userInputForm.TextInput.Label,
				Required:  userInputForm.TextInput.Required,
				MaxLength: userInputForm.TextInput.MaxLength,
			})
		}
	}
	return variables
}
