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
		AgentMode:                     ConvertToBizAgentMode(dtoAppConfig.AgentMode),
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

func ConvertToBizAgentTools(agentTools []*dto.AgentTools) []*biz_entity.AgentTools {
	poAgentTools := make([]*biz_entity.AgentTools, 0, len(agentTools))

	for _, agentTool := range agentTools {
		poAgentTools = append(poAgentTools, &biz_entity.AgentTools{
			Enabled:        agentTool.Enabled,
			ProviderID:     agentTool.ProviderID,
			ProviderName:   agentTool.ProviderName,
			ProviderType:   agentTool.ProviderType,
			ToolLabel:      agentTool.ToolLabel,
			ToolName:       agentTool.ToolName,
			ToolParameters: agentTool.ToolParameters,
		})
	}
	return poAgentTools
}

func ConvertToBizAgentMode(agentMode *dto.AgentMode) *biz_entity.AgentMode {
	if agentMode == nil {
		return &biz_entity.AgentMode{
			Tools: make([]*biz_entity.AgentTools, 0),
		}
	}
	return &biz_entity.AgentMode{
		Enabled:        agentMode.Enabled,
		MaxInteraction: agentMode.MaxInteraction,
		Prompt:         agentMode.Prompt,
		Strategy:       agentMode.Strategy,
		Tools:          ConvertToBizAgentTools(agentMode.Tools),
	}
}
