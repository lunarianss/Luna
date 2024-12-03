package assembler

import (
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/provider_app_config"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
)

// ConvertToConfigEntity converts a DTO (Data Transfer Object) to a business entity.
func ConvertToConfigEntity(dtoAppConfig *dto.AppModelConfigDto) *biz_entity.AppModelConfig {
	return &biz_entity.AppModelConfig{
		AppID:                         dtoAppConfig.AppID,
		ModelID:                       dtoAppConfig.ModelID,
		Configs:                       dtoAppConfig.Configs,
		OpeningStatement:              dtoAppConfig.OpeningStatement,
		SuggestedQuestions:            dtoAppConfig.SuggestedQuestions,
		SuggestedQuestionsAfterAnswer: biz_entity.AppModelConfigEnable(dtoAppConfig.SuggestedQuestionsAfterAnswer),
		MoreLikeThis:                  biz_entity.AppModelConfigEnable(dtoAppConfig.MoreLikeThis),
		Model:                         ConvertToModelEntity(dtoAppConfig.Model),
		UserInputForm:                 ConvertToUserInputEntity(dtoAppConfig.UserInputForm),
		PrePrompt:                     dtoAppConfig.PrePrompt,
		AgentMode:                     dtoAppConfig.AgentMode,
		SpeechToText:                  biz_entity.AppModelConfigEnable(dtoAppConfig.SpeechToText),
		SensitiveWordAvoidance:        dtoAppConfig.SensitiveWordAvoidance,
		RetrieverResource:             biz_entity.AppModelConfigEnable(dtoAppConfig.RetrieverResource),
		DatasetQueryVariable:          dtoAppConfig.DatasetQueryVariable,
		PromptType:                    dtoAppConfig.PromptType,
		ChatPromptConfig:              dtoAppConfig.ChatPromptConfig,
		CompletionPromptConfig:        dtoAppConfig.CompletionPromptConfig,
		DatasetConfigs:                dtoAppConfig.DatasetConfigs,
		ExternalDataTools:             dtoAppConfig.ExternalDataTools,
		FileUpload:                    dtoAppConfig.FileUpload,
		TextToSpeech:                  biz_entity.AppModelConfigEnable(dtoAppConfig.TextToSpeech),
	}
}
// ConvertToModelEntity converts a ModelDto to a biz_entity.Model.
func ConvertToModelEntity(dtoModel dto.ModelDto) biz_entity.ModelInfo {
	return biz_entity.ModelInfo{
		Provider:         dtoModel.Provider,
		Name:             dtoModel.Name,
		Mode:             dtoModel.Mode,
		CompletionParams: dtoModel.CompletionParams,
	}
}

func ConvertToUserInputEntity(dtoModels []*dto.UserInputForm) []*biz_entity.UserInputForm {
	var returnUserInput []*biz_entity.UserInputForm
	var baseUserTextInput *biz_entity.BaseTextUserInput
	var userInputForm *biz_entity.UserInputForm

	for _, dtoModel := range dtoModels {
		userInputForm = &biz_entity.UserInputForm{}
		if dtoModel.TextInput != nil {
			baseUserTextInput = &biz_entity.BaseTextUserInput{
				Label:     dtoModel.TextInput.Label,
				Variable:  dtoModel.TextInput.Variable,
				Required:  dtoModel.TextInput.Required,
				MaxLength: dtoModel.TextInput.MaxLength,
				Default:   dtoModel.TextInput.Default,
			}
		}
		userInputForm.TextInput = baseUserTextInput
		returnUserInput = append(returnUserInput, userInputForm)
	}
	return returnUserInput
}
