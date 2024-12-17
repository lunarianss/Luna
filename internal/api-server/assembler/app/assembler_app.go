package assembler

import (
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/provider_app_config"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	appDto "github.com/lunarianss/Luna/internal/api-server/dto/app"
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

func ConvertToUserInputEntity(userInputs []dto.UserInputForm) []biz_entity.UserInputForm {
	var returnUserInput []biz_entity.UserInputForm

	for _, userInputMap := range userInputs {
		userInputForm := biz_entity.UserInputForm{}
		for k, v := range userInputMap {
			userInputForm[k] = &biz_entity.UserInput{
				Label:     v.Label,
				Variable:  v.Variable,
				Required:  v.Required,
				MaxLength: v.MaxLength,
				Default:   v.Default,
				Options:   v.Options,
			}
		}
		returnUserInput = append(returnUserInput, userInputForm)
	}

	return returnUserInput
}

func ConvertToServiceTokens(serviceTokens []*po_entity.ApiToken) []*appDto.GenerateServiceToken {
	var dtoServiceTokens []*appDto.GenerateServiceToken

	for _, serviceToken := range serviceTokens {

		dtoServiceTokens = append(dtoServiceTokens, &appDto.GenerateServiceToken{
			ID:        serviceToken.ID,
			Type:      serviceToken.Type,
			Token:     serviceToken.Token,
			CreatedAt: serviceToken.CreatedAt,
		})
	}
	return dtoServiceTokens
}
