package entities

import (
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config"
	biz_entity_model "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider/model_provider"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_configuration"
)

// Define types for your enums
type InvokeFrom string

const (
	ServiceAPI InvokeFrom = "service-api"
	WebApp     InvokeFrom = "web-app"
	Explore    InvokeFrom = "explore"
	Debugger   InvokeFrom = "debugger"
)

// ModelConfigWithCredentialsEntity struct
type ModelConfigWithCredentialsEntity struct {
	Provider            string                                       `json:"provider"`
	Model               string                                       `json:"model"`
	ModelSchema         *biz_entity_model.AIModelStaticConfiguration `json:"model_schema"`
	Mode                string                                       `json:"mode"`
	ProviderModelBundle *biz_entity.ProviderModelBundleRuntime       `json:"provider_model_bundle"`
	Credentials         interface{}                                  `json:"credentials"`
	Parameters          map[string]interface{}                       `json:"parameters"`
	Stop                []string                                     `json:"stop"`
}

// AppGenerateEntity struct
type AppGenerateEntity struct {
	TaskID     string                 `json:"task_id"`
	AppConfig  *AppConfig             `json:"app_config"`
	Inputs     map[string]interface{} `json:"inputs"`
	UserID     string                 `json:"user_id"`
	Stream     bool                   `json:"stream"`
	InvokeFrom InvokeFrom             `json:"invoke_from"`
	CallDepth  int                    `json:"call_depth"`
	Extras     map[string]interface{} `json:"extras"`
}

// EasyUIBasedAppGenerateEntity struct
type EasyUIBasedAppGenerateEntity struct {
	*AppGenerateEntity
	AppConfig *app_config.EasyUIBasedAppConfig  `json:"app_config"`
	ModelConf *ModelConfigWithCredentialsEntity `json:"model_conf"`
	Query     string                            `json:"query"`
}

// ConversationAppGenerateEntity struct
type ConversationAppGenerateEntity struct {
	*AppGenerateEntity
	ConversationID  *string `json:"conversation_id"`
	ParentMessageID *string `json:"parent_message_id"`
}

// ChatAppGenerateEntity struct
type ChatAppGenerateEntity struct {
	*EasyUIBasedAppGenerateEntity
	ConversationID  *string `json:"conversation_id"`
	ParentMessageID *string `json:"parent_message_id"`
}

type CreatedByRole string

const (
	CreatedByRoleAccount CreatedByRole = "account"
	CreatedByRoleEndUser CreatedByRole = "end_user"
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

type ModelConfigEntity struct {
	Provider   string                 `json:"provider"`
	Model      string                 `json:"model"`
	Mode       string                 `json:"mode"`
	Parameters map[string]interface{} `json:"parameters"`
	Stop       []string               `json:"stop"`
}

type AdvancedChatMessageEntity struct {
	Text string `json:"text"`
	Role string `json:"role"` // Assuming PromptMessageRole is defined as string
}

type AdvancedChatPromptTemplateEntity struct {
	Messages []*AdvancedChatMessageEntity `json:"messages"`
}

type RolePrefixEntity struct {
	User      string `json:"user"`
	Assistant string `json:"assistant"`
}

type AdvancedCompletionPromptTemplateEntity struct {
	Prompt     string            `json:"prompt"`
	RolePrefix *RolePrefixEntity `json:"role_prefix"`
}

type PromptTemplateEntity struct {
	PromptType                       string                                  `json:"prompt_type"`
	SimplePromptTemplate             string                                  `json:"simple_prompt_template"`
	AdvancedChatPromptTemplate       *AdvancedChatPromptTemplateEntity       `json:"advanced_chat_prompt_template"`
	AdvancedCompletionPromptTemplate *AdvancedCompletionPromptTemplateEntity `json:"advanced_completion_prompt_template"`
}

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

type VariableEntity struct {
	Variable                 string             `json:"variable"`
	Label                    string             `json:"label"`
	Description              string             `json:"description"`
	Type                     VariableEntityType `json:"type"`
	Required                 bool               `json:"required"`
	MaxLength                int                `json:"max_length"`
	Options                  []string           `json:"options"`
	AllowedFileTypes         []string           `json:"allowed_file_types"`
	AllowedFileExtensions    []string           `json:"allowed_file_extensions"`
	AllowedFileUploadMethods []string           `json:"allowed_file_upload_methods"`
}

type ExternalDataVariableEntity struct {
	Variable string                 `json:"variable"`
	Type     string                 `json:"type"`
	Config   map[string]interface{} `json:"config"`
}

type RetrieveStrategy string

const (
	Single   RetrieveStrategy = "single"
	Multiple RetrieveStrategy = "multiple"
)

type DatasetRetrieveConfigEntity struct {
	QueryVariable    string                 `json:"query_variable"`
	RetrieveStrategy RetrieveStrategy       `json:"retrieve_strategy"`
	TopK             int                    `json:"top_k"`
	ScoreThreshold   float64                `json:"score_threshold"`
	RerankMode       string                 `json:"rerank_mode"`
	RerankingModel   map[string]interface{} `json:"reranking_model"`
	Weights          map[string]interface{} `json:"weights"`
	RerankingEnabled bool                   `json:"reranking_enabled"`
}

type DatasetEntity struct {
	DatasetIDs     []string                    `json:"dataset_ids"`
	RetrieveConfig DatasetRetrieveConfigEntity `json:"retrieve_config"`
}

type SensitiveWordAvoidanceEntity struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

type TextToSpeechEntity struct {
	Enabled  bool   `json:"enabled"`
	Voice    string `json:"voice"`
	Language string `json:"language"`
}

type TracingConfigEntity struct {
	Enabled         bool   `json:"enabled"`
	TracingProvider string `json:"tracing_provider"`
}

type AppAdditionalFeatures struct {
	FileUpload                    string               `json:"file_upload"`
	OpeningStatement              string               `json:"opening_statement"`
	SuggestedQuestions            []string             `json:"suggested_questions"`
	SuggestedQuestionsAfterAnswer bool                 `json:"suggested_questions_after_answer"`
	ShowRetrieveSource            bool                 `json:"show_retrieve_source"`
	MoreLikeThis                  bool                 `json:"more_like_this"`
	SpeechToText                  bool                 `json:"speech_to_text"`
	TextToSpeech                  *TextToSpeechEntity  `json:"text_to_speech"`
	TraceConfig                   *TracingConfigEntity `json:"trace_config"`
}

type AppConfig struct {
	TenantID               string                        `json:"tenant_id"`
	AppID                  string                        `json:"app_id"`
	AppMode                string                        `json:"app_mode"` // Assuming AppMode is defined as string
	AdditionalFeatures     *AppAdditionalFeatures        `json:"additional_features"`
	Variables              []*VariableEntity             `json:"variables"`
	SensitiveWordAvoidance *SensitiveWordAvoidanceEntity `json:"sensitive_word_avoidance"`
}

type EasyUIBasedAppModelConfigFrom string

const (
	Args                       EasyUIBasedAppModelConfigFrom = "args"
	AppLatestConfig            EasyUIBasedAppModelConfigFrom = "app-latest-config"
	ConversationSpecificConfig EasyUIBasedAppModelConfigFrom = "conversation-specific-config"
)

type EasyUIBasedAppConfig struct {
	*AppConfig
	AppModelConfigFrom    EasyUIBasedAppModelConfigFrom `json:"app_model_config_from"`
	AppModelConfigID      string                        `json:"app_model_config_id"`
	AppModelConfigDict    map[string]interface{}        `json:"app_model_config_dict"`
	Model                 *ModelConfigEntity            `json:"model"`
	PromptTemplate        *PromptTemplateEntity         `json:"prompt_template"`
	Dataset               *DatasetEntity                `json:"dataset"`
	ExternalDataVariables []ExternalDataVariableEntity  `json:"external_data_variables"`
}

type WorkflowUIBasedAppConfig struct {
	*AppConfig
	WorkflowID string `json:"workflow_id"`
}

type ChatAppConfig struct {
	*EasyUIBasedAppConfig
}
