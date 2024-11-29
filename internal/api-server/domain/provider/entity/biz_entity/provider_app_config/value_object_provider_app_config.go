package biz_entity

type AppModelConfigEnable struct {
	Enable bool `json:"enable"`
}

// Model holds the model-specific configuration.
type ModelInfo struct {
	Provider         string                 `json:"provider"`
	Name             string                 `json:"name"`
	Mode             string                 `json:"mode"`
	CompletionParams map[string]interface{} `json:"completion_params"`
}

type AppModelConfig struct {
	AppID                         string                              `json:"app_id" gorm:"column:app_id"`
	Provider                      string                              `json:"provider" gorm:"column:provider"`
	ModelID                       string                              `json:"model_id" gorm:"column:model_id"`
	Configs                       map[string]interface{}              `json:"configs" gorm:"column:configs;serializer:json"`
	CreatedAt                     int64                               `json:"created_at" gorm:"column:created_at"`
	UpdatedAt                     int64                               `json:"updated_at" gorm:"column:updated_at"`
	OpeningStatement              map[string]interface{}              `json:"opening_statement" gorm:"column:opening_statement;serializer:json"`
	SuggestedQuestions            []string                            `json:"suggested_questions" gorm:"column:suggested_questions;serializer:json"`
	SuggestedQuestionsAfterAnswer AppModelConfigEnable                `json:"suggested_questions_after_answer" gorm:"column:suggested_questions_after_answer;serializer:json"`
	MoreLikeThis                  AppModelConfigEnable                `json:"more_like_this" gorm:"column:more_like_this;serializer:json"`
	Model                         ModelInfo                           `json:"model" gorm:"column:model;serializer:json"`
	UserInputForm                 []map[string]map[string]interface{} `json:"user_input_form" gorm:"column:user_input_form;serializer:json"`
	PrePrompt                     string                              `json:"pre_prompt" gorm:"column:pre_prompt;serializer:json"`
	AgentMode                     map[string]interface{}              `json:"agent_mode" gorm:"column:agent_mode;serializer:json"`
	SpeechToText                  AppModelConfigEnable                `json:"speech_to_text" gorm:"column:speech_to_text;serializer:json"`
	SensitiveWordAvoidance        map[string]interface{}              `json:"sensitive_word_avoidance" gorm:"column:sensitive_word_avoidance;serializer:json"`
	RetrieverResource             AppModelConfigEnable                `json:"retriever_resource" gorm:"column:retriever_resource;serializer:json"`
	DatasetQueryVariable          map[string]interface{}              `json:"dataset_query_variable" gorm:"column:dataset_query_variable;serializer:json"`
	PromptType                    string                              `json:"prompt_type" gorm:"column:prompt_type"`
	ChatPromptConfig              map[string]interface{}              `json:"chat_prompt_config" gorm:"column:chat_prompt_config;serializer:json"`
	CompletionPromptConfig        map[string]interface{}              `json:"completion_prompt_config" gorm:"column:completion_prompt_config;serializer:json"`
	DatasetConfigs                map[string]interface{}              `json:"dataset_configs" gorm:"column:dataset_configs;serializer:json"`
	ExternalDataTools             []string                            `json:"external_data_tools" gorm:"column:external_data_tools;serializer:json"`
	FileUpload                    map[string]map[string]interface{}   `json:"file_upload" gorm:"column:file_upload;serializer:json"`
	TextToSpeech                  AppModelConfigEnable                `json:"text_to_speech" gorm:"column:text_to_speech;serializer:json"`
}
