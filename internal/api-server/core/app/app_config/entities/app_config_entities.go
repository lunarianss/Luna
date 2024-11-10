package entities

import (
	"github.com/lunarianss/Luna/internal/api-server/core/app/file"
	"github.com/lunarianss/Luna/internal/api-server/entities/message"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
)

// ModelConfigEntity represents the model configuration
type ModelConfigEntity struct {
	Provider   string                 `json:"provider"`
	Model      string                 `json:"model"`
	Mode       string                 `json:"mode,omitempty"`
	Parameters map[string]interface{} `json:"parameters"`
	Stop       []string               `json:"stop"`
}

// AdvancedChatMessageEntity represents an advanced chat message
type AdvancedChatMessageEntity struct {
	Text string                     `json:"text"`
	Role *message.PromptMessageRole `json:"role"`
}

// AdvancedChatPromptTemplateEntity holds messages for a chat prompt template
type AdvancedChatPromptTemplateEntity struct {
	Messages []AdvancedChatMessageEntity `json:"messages"`
}

// AdvancedCompletionPromptTemplateEntity holds a prompt template and optional role prefix
type AdvancedCompletionPromptTemplateEntity struct {
	Prompt     string
	RolePrefix *RolePrefixEntity `json:"role_prefix,omitempty"`
}

// RolePrefixEntity represents user and assistant prefixes
type RolePrefixEntity struct {
	User      string `json:"user"`
	Assistant string `json:"assistant"`
}

// PromptTemplateEntity represents a prompt template with simple and advanced options
type PromptTemplateEntity struct {
	PromptType                       PromptType                              `json:"prompt_type"`
	SimplePromptTemplate             string                                  `json:"simple_prompt_template,omitempty"`
	AdvancedChatPromptTemplate       *AdvancedChatPromptTemplateEntity       `json:"advanced_chat_prompt_template,omitempty"`
	AdvancedCompletionPromptTemplate *AdvancedCompletionPromptTemplateEntity `json:"advanced_completion_prompt_template,omitempty"`
}

// PromptType enum values
type PromptType string

const (
	SimplePromptType   PromptType = "simple"
	AdvancedPromptType PromptType = "advanced"
)

// VariableEntityType enum values
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

// VariableEntity represents a variable configuration
type VariableEntity struct {
	Variable                 string                    `json:"variable"`
	Label                    string                    `json:"label"`
	Description              string                    `json:"description"`
	Type                     VariableEntityType        `json:"type"`
	Required                 bool                      `json:"required"`
	MaxLength                int                       `json:"max_length,omitempty"`
	Options                  []string                  `json:"options"`
	AllowedFileTypes         []file.FileType           `json:"allowed_file_types"`
	AllowedFileExtensions    []string                  `json:"allowed_file_extensions"`
	AllowedFileUploadMethods []file.FileTransferMethod `json:"allowed_file_upload_methods"`
}

// ExternalDataVariableEntity represents an external data variable
type ExternalDataVariableEntity struct {
	Variable string                 `json:"variable"`
	Type     string                 `json:"type"`
	Config   map[string]interface{} `json:"config"`
}

// DatasetRetrieveConfigEntity represents configuration for dataset retrieval
type DatasetRetrieveConfigEntity struct {
	QueryVariable    *string                 `json:"query_variable,omitempty"`
	RetrieveStrategy RetrieveStrategy        `json:"retrieve_strategy"`
	TopK             *int                    `json:"top_k,omitempty"`
	ScoreThreshold   *float64                `json:"score_threshold,omitempty"`
	RerankMode       *string                 `json:"rerank_mode,omitempty"`
	RerankingModel   *map[string]interface{} `json:"reranking_model,omitempty"`
	Weights          *map[string]interface{} `json:"weights,omitempty"`
	RerankingEnabled *bool                   `json:"reranking_enabled,omitempty"`
}

// RetrieveStrategy enum values
type RetrieveStrategy string

const (
	Single   RetrieveStrategy = "single"
	Multiple RetrieveStrategy = "multiple"
)

// DatasetEntity represents a dataset configuration
type DatasetEntity struct {
	DatasetIDs     []string                    `json:"dataset_ids"`
	RetrieveConfig DatasetRetrieveConfigEntity `json:"retrieve_config"`
}

// SensitiveWordAvoidanceEntity represents configuration for sensitive word avoidance
type SensitiveWordAvoidanceEntity struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

// TextToSpeechEntity represents configuration for text-to-speech features
type TextToSpeechEntity struct {
	Enabled  bool    `json:"enabled"`
	Voice    *string `json:"voice,omitempty"`
	Language *string `json:"language,omitempty"`
}

// TracingConfigEntity represents tracing configuration
type TracingConfigEntity struct {
	Enabled         bool   `json:"enabled"`
	TracingProvider string `json:"tracing_provider"`
}

// AppAdditionalFeatures represents additional application features
type AppAdditionalFeatures struct {
	FileUpload                    *file.FileExtraConfig `json:"file_upload,omitempty"`
	OpeningStatement              string                `json:"opening_statement,omitempty"`
	SuggestedQuestions            []string              `json:"suggested_questions"`
	SuggestedQuestionsAfterAnswer bool                  `json:"suggested_questions_after_answer"`
	ShowRetrieveSource            bool                  `json:"show_retrieve_source"`
	MoreLikeThis                  bool                  `json:"more_like_this"`
	SpeechToText                  bool                  `json:"speech_to_text"`
	TextToSpeech                  *TextToSpeechEntity   `json:"text_to_speech,omitempty"`
	TraceConfig                   *TracingConfigEntity  `json:"trace_config,omitempty"`
}

// AppConfig represents the main configuration for an application
type AppConfig struct {
	TenantID               string                        `json:"tenant_id"`
	AppID                  string                        `json:"app_id"`
	AppMode                model.AppMode                 `json:"app_mode"`
	AdditionalFeatures     *AppAdditionalFeatures        `json:"additional_features"`
	Variables              []*VariableEntity             `json:"variables"`
	SensitiveWordAvoidance *SensitiveWordAvoidanceEntity `json:"sensitive_word_avoidance,omitempty"`
}
