package biz_entity

import (
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
)

type CreatedByRole string

const (
	CreatedByRoleAccount CreatedByRole = "account"
	CreatedByRoleEndUser CreatedByRole = "end_user"
)

type EasyUIBasedAppModelConfigFrom string

const (
	Args                       EasyUIBasedAppModelConfigFrom = "args"
	AppLatestConfig            EasyUIBasedAppModelConfigFrom = "app-latest-config"
	ConversationSpecificConfig EasyUIBasedAppModelConfigFrom = "conversation-specific-config"
)

type PromptType string

const (
	SIMPLE   PromptType = "simple"
	ADVANCED PromptType = "advanced"
)

type UserFrom string

const (
	UserFromAccount UserFrom = "account"
	UserFromEndUser UserFrom = "end-user"
)

type WorkflowRunTriggeredFrom string

const (
	WorkflowRunTriggeredFromDebugging WorkflowRunTriggeredFrom = "debugging"
	WorkflowRunTriggeredFromAppRun    WorkflowRunTriggeredFrom = "app-run"
)

type VariableEntityType string

const (
	TextInput        VariableEntityType = "text-input"
	Select           VariableEntityType = "select"
	Paragraph        VariableEntityType = "paragraph"
	Number           VariableEntityType = "number"
	ExternalDataTool VariableEntityType = "external_data_tool"
	File             VariableEntityType = "file"
	FileList         VariableEntityType = "file-list"
)

type AppModelConfigEnable struct {
	Enabled bool `json:"enabled"`
}

type ModelInfo struct {
	Provider         string                 `json:"provider"`
	Name             string                 `json:"name"`
	Mode             string                 `json:"mode"`
	CompletionParams map[string]interface{} `json:"completion_params"`
}

type UserInput struct {
	Label     string   `json:"label"`
	Variable  string   `json:"variable"`
	Required  bool     `json:"required"`
	MaxLength int      `json:"max_length"`
	Default   string   `json:"default"`
	Options   []string `json:"options"`
}

type AgentTools struct {
	Enabled        bool           `json:"enabled"`
	ProviderID     string         `json:"provider_id"`
	ProviderName   string         `json:"provider_name"`
	ProviderType   string         `json:"provider_type"`
	ToolLabel      string         `json:"tool_label"`
	ToolName       string         `json:"tool_name"`
	ToolParameters map[string]any `json:"tool_parameters"`
}

type AgentMode struct {
	Enabled        bool          `json:"enabled"`
	MaxInteraction int           `json:"max_interaction"`
	Prompt         string        `json:"prompt"`
	Strategy       string        `json:"strategy"`
	Tools          []*AgentTools `json:"tools"`
}

type UserInputForm map[string]*UserInput

type AppModelConfig struct {
	AppID                         string                 `json:"app_id" gorm:"column:app_id"`
	Provider                      string                 `json:"provider" gorm:"column:provider"`
	ModelID                       string                 `json:"model_id" gorm:"column:model_id"`
	Configs                       map[string]interface{} `json:"configs" gorm:"column:configs;serializer:json"`
	CreatedAt                     int64                  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt                     int64                  `json:"updated_at" gorm:"column:updated_at"`
	OpeningStatement              string                 `json:"opening_statement" gorm:"column:opening_statement;serializer:json"`
	SuggestedQuestions            []string               `json:"suggested_questions" gorm:"column:suggested_questions;serializer:json"`
	SuggestedQuestionsAfterAnswer AppModelConfigEnable   `json:"suggested_questions_after_answer" gorm:"column:suggested_questions_after_answer;serializer:json"`
	MoreLikeThis                  AppModelConfigEnable   `json:"more_like_this" gorm:"column:more_like_this;serializer:json"`
	Model                         ModelInfo              `json:"model" gorm:"column:model;serializer:json"`
	UserInputForm                 []UserInputForm        `json:"user_input_form" gorm:"column:user_input_form;serializer:json"`
	PrePrompt                     string                 `json:"pre_prompt" gorm:"column:pre_prompt;serializer:json"`
	AgentMode                     *AgentMode             `json:"agent_mode" gorm:"column:agent_mode;serializer:json"`
	SpeechToText                  AppModelConfigEnable   `json:"speech_to_text" gorm:"column:speech_to_text;serializer:json"`
	SensitiveWordAvoidance        map[string]interface{} `json:"sensitive_word_avoidance" gorm:"column:sensitive_word_avoidance;serializer:json"`
	RetrieverResource             AppModelConfigEnable   `json:"retriever_resource" gorm:"column:retriever_resource;serializer:json"`
	DatasetQueryVariable          string                 `json:"dataset_query_variable" gorm:"column:dataset_query_variable;serializer:json"`
	PromptType                    string                 `json:"prompt_type" gorm:"column:prompt_type"`
	ChatPromptConfig              map[string]interface{} `json:"chat_prompt_config" gorm:"column:chat_prompt_config;serializer:json"`
	CompletionPromptConfig        map[string]interface{} `json:"completion_prompt_config" gorm:"column:completion_prompt_config;serializer:json"`
	DatasetConfigs                map[string]interface{} `json:"dataset_configs" gorm:"column:dataset_configs;serializer:json"`
	ExternalDataTools             []string               `json:"external_data_tools" gorm:"column:external_data_tools;serializer:json"`
	FileUpload                    map[string]interface{} `json:"file_upload" gorm:"column:file_upload;serializer:json"`
	TextToSpeech                  AppModelConfigEnable   `json:"text_to_speech" gorm:"column:text_to_speech;serializer:json"`
	AppAnnotationReply            *AppAnnotationReply    `json:"annotation_reply"`
}

func (a *AppModelConfig) ConvertToAppConfigPoEntity() *po_entity.AppModelConfig {
	return &po_entity.AppModelConfig{
		AppID:                         a.AppID,
		Provider:                      a.Provider,
		ModelID:                       a.ModelID,
		Configs:                       a.Configs,
		CreatedAt:                     a.CreatedAt,
		UpdatedAt:                     a.UpdatedAt,
		OpeningStatement:              a.OpeningStatement,
		SuggestedQuestions:            a.SuggestedQuestions,
		SuggestedQuestionsAfterAnswer: po_entity.AppModelConfigEnable(a.SuggestedQuestionsAfterAnswer),
		MoreLikeThis:                  po_entity.AppModelConfigEnable(a.MoreLikeThis),
		Model:                         ConvertToModelPoEntity(a.Model),
		UserInputForm:                 ConvertToUserInputPoEntity(a.UserInputForm),
		PrePrompt:                     a.PrePrompt,
		AgentMode:                     ConvertToAgentMode(a.AgentMode),
		SpeechToText:                  po_entity.AppModelConfigEnable(a.SpeechToText),
		SensitiveWordAvoidance:        a.SensitiveWordAvoidance,
		RetrieverResource:             po_entity.AppModelConfigEnable(a.RetrieverResource),
		DatasetQueryVariable:          a.DatasetQueryVariable,
		PromptType:                    a.PromptType,
		ChatPromptConfig:              a.ChatPromptConfig,
		CompletionPromptConfig:        a.CompletionPromptConfig,
		DatasetConfigs:                a.DatasetConfigs,
		ExternalDataTools:             a.ExternalDataTools,
		FileUpload:                    a.FileUpload,
		TextToSpeech:                  po_entity.AppModelConfigEnable(a.TextToSpeech),
	}
}

// ConvertToModelEntity converts a ModelDto to a biz_entity.Model.
func ConvertToModelPoEntity(entityModel ModelInfo) po_entity.ModelInfo {
	return po_entity.ModelInfo{
		Provider:         entityModel.Provider,
		Name:             entityModel.Name,
		Mode:             entityModel.Mode,
		CompletionParams: entityModel.CompletionParams,
	}
}

// ConvertToModelEntity converts a ModelDto to a biz_entity.Model.
func ConvertToModelBizEntity(entityModel po_entity.ModelInfo) ModelInfo {
	return ModelInfo{
		Provider:         entityModel.Provider,
		Name:             entityModel.Name,
		Mode:             entityModel.Mode,
		CompletionParams: entityModel.CompletionParams,
	}
}

func ConvertToAgentTools(agentTools []*AgentTools) []*po_entity.AgentTools {
	poAgentTools := make([]*po_entity.AgentTools, 0, len(agentTools))

	for _, agentTool := range agentTools {
		poAgentTools = append(poAgentTools, &po_entity.AgentTools{
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

func ConvertToAgentMode(agentMode *AgentMode) *po_entity.AgentMode {
	if agentMode == nil {
		return &po_entity.AgentMode{
			Tools: make([]*po_entity.AgentTools, 0),
		}
	}

	return &po_entity.AgentMode{
		Enabled:        agentMode.Enabled,
		MaxInteraction: agentMode.MaxInteraction,
		Prompt:         agentMode.Prompt,
		Strategy:       agentMode.Strategy,
		Tools:          ConvertToAgentTools(agentMode.Tools),
	}
}

func ConvertToBizAgentTools(agentTools []*po_entity.AgentTools) []*AgentTools {

	poAgentTools := make([]*AgentTools, 0, len(agentTools))

	for _, agentTool := range agentTools {
		poAgentTools = append(poAgentTools, &AgentTools{
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

func ConvertToBizAgentMode(agentMode *po_entity.AgentMode) *AgentMode {
	if agentMode == nil {
		return &AgentMode{
			Tools: make([]*AgentTools, 0),
		}
	}

	return &AgentMode{
		Enabled:        agentMode.Enabled,
		MaxInteraction: agentMode.MaxInteraction,
		Prompt:         agentMode.Prompt,
		Strategy:       agentMode.Strategy,
		Tools:          ConvertToBizAgentTools(agentMode.Tools),
	}
}

func ConvertToUserInputPoEntity(userInputs []UserInputForm) []po_entity.UserInputForm {
	var returnUserInput []po_entity.UserInputForm

	for _, userInputMap := range userInputs {
		userInputForm := po_entity.UserInputForm{}
		for k, v := range userInputMap {
			userInputForm[k] = &po_entity.UserInput{
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

func ConvertToUserInputBizEntity(userInputs []po_entity.UserInputForm) []UserInputForm {
	var returnUserInput []UserInputForm

	for _, userInputMap := range userInputs {
		userInputForm := UserInputForm{}
		for k, v := range userInputMap {
			userInputForm[k] = &UserInput{
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

type AppAnnotationEmbeddingModel struct {
	EmbeddingProviderName string `json:"embedding_provider_name,omitempty"`
	EmbeddingModelName    string `json:"embedding_model_name,omitempty"`
}

type AppAnnotationReply struct {
	ID             string                       `json:"id,omitempty"`
	Enabled        bool                         `json:"enabled"`
	ScoreThreshold float32                      `json:"score_threshold,omitempty"`
	EmbeddingModel *AppAnnotationEmbeddingModel `json:"embedding_model,omitempty"`
}

func ConvertToAppConfigBizEntity(a *po_entity.AppModelConfig, annotation *AppAnnotationReply) *AppModelConfig {
	return &AppModelConfig{
		AppID:                         a.AppID,
		Provider:                      a.Provider,
		ModelID:                       a.ModelID,
		Configs:                       a.Configs,
		CreatedAt:                     a.CreatedAt,
		UpdatedAt:                     a.UpdatedAt,
		OpeningStatement:              a.OpeningStatement,
		SuggestedQuestions:            a.SuggestedQuestions,
		SuggestedQuestionsAfterAnswer: AppModelConfigEnable(a.SuggestedQuestionsAfterAnswer), // 注意类型转换
		MoreLikeThis:                  AppModelConfigEnable(a.MoreLikeThis),                  // 注意类型转换
		Model:                         ConvertToModelBizEntity(a.Model),                      // 假设 Model 是直接可以赋值的，如果不是需要进行类型转换
		UserInputForm:                 ConvertToUserInputBizEntity(a.UserInputForm),
		PrePrompt:                     a.PrePrompt,
		AgentMode:                     ConvertToBizAgentMode(a.AgentMode),
		SpeechToText:                  AppModelConfigEnable(a.SpeechToText), // 注意类型转换
		SensitiveWordAvoidance:        a.SensitiveWordAvoidance,
		RetrieverResource:             AppModelConfigEnable(a.RetrieverResource), // 注意类型转换
		DatasetQueryVariable:          a.DatasetQueryVariable,
		PromptType:                    a.PromptType,
		ChatPromptConfig:              a.ChatPromptConfig,
		CompletionPromptConfig:        a.CompletionPromptConfig,
		DatasetConfigs:                a.DatasetConfigs,
		ExternalDataTools:             a.ExternalDataTools,
		FileUpload:                    a.FileUpload,
		TextToSpeech:                  AppModelConfigEnable(a.TextToSpeech), // 注意类型转换
		AppAnnotationReply:            annotation,
	}
}
