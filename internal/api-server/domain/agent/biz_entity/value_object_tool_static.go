package biz_entity

import common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"

type ToolParameterForm string

const (
	SchemaForm ToolParameterForm = "schema" // should be set while adding tool
	FormForm   ToolParameterForm = "form"   // should be set before invoking tool
	LLMForm    ToolParameterForm = "llm"    // will be set by LLM
)

type ToolParameterType string

const (
	StringType      ToolParameterType = "string"
	NumberType      ToolParameterType = "number"
	BooleanType     ToolParameterType = "boolean"
	SelectType      ToolParameterType = "select"
	SecretInputType ToolParameterType = "secret-input"
	FileType        ToolParameterType = "file"
	FilesType       ToolParameterType = "files"
	SystemFilesType ToolParameterType = "systme-files" // Deprecated
)

type ToolInvokeFrom string

const (
	WorkflowInvoke ToolInvokeFrom = "workflow"
	AgentInvoke    ToolInvokeFrom = "agent"
)

type ToolRuntimeVariableType string

const (
	TextType  ToolRuntimeVariableType = "text"
	ImageType ToolRuntimeVariableType = "image"
)

type InvokeFrom string

const (
	ServiceAPIInvoke InvokeFrom = "service-api"
	WebAppInvoke     InvokeFrom = "web-app"
	ExploreInvoke    InvokeFrom = "explore"
	DebuggerInvoke   InvokeFrom = "debugger"
)

type ToolParameterOption struct {
	Value string             `json:"value" yaml:"value"` // The value of the option
	Label *common.I18nObject `json:"label" yaml:"label"` // The label of the option
}

type ToolParameter struct {
	Name             string                 `json:"name" yaml:"name"`                           // The name of the parameter
	Label            *common.I18nObject     `json:"label" yaml:"label"`                         // The label presented to the user
	HumanDescription *common.I18nObject     `json:"human_description" yaml:"human_description"` // The description presented to the user
	Placeholder      *common.I18nObject     `json:"placeholder" yaml:"placeholder"`             // The placeholder presented to the user
	Type             ToolParameterType      `json:"type" yaml:"type"`                           // The type of the parameter
	Form             ToolParameterForm      `json:"form" yaml:"form"`                           // The form of the parameter, schema/form/llm
	LLMDescription   string                 `json:"llm_description" yaml:"llm_description"`     // Description set by LLM
	Required         bool                   `json:"required" yaml:"required"`                   // Whether the parameter is required
	Default          any                    `json:"default" yaml:"default"`                     // Default value for the parameter
	Min              float64                `json:"min" yaml:"min"`                             // Minimum value
	Max              float64                `json:"max" yaml:"max"`                             // Maximum value
	Options          []*ToolParameterOption `json:"options" yaml:"options"`                     // Options for select type
}

type ToolDescription struct {
	Human *common.I18nObject `json:"human" yaml:"human"`
	LLM   string             `json:"llm" yaml:"llm"`
}

type ToolRuntimeVariable struct {
	Type     ToolRuntimeVariableType `json:"type" yaml:"type"`           // The type of the variable
	Name     string                  `json:"name" yaml:"name"`           // The name of the variable
	Position int                     `json:"position" yaml:"position"`   // The position of the variable
	ToolName string                  `json:"tool_name" yaml:"tool_name"` // The name of the tool
}

type ToolRuntimeVariablePool struct {
	ConversationID string                 `json:"conversation_id" yaml:"conversation_id"` // The conversation id
	UserID         string                 `json:"user_id" yaml:"user_id"`                 // The user id
	TenantID       string                 `json:"tenant_id" yaml:"tenant_id"`             // The tenant id of assistant
	Pool           []*ToolRuntimeVariable `json:"pool" yaml:"pool"`                       // The pool of variables
}

type ToolIdentity struct {
	Author   string             `json:"author" yaml:"author"`                 // The author of the tool
	Name     string             `json:"name" yaml:"name"`                     // The name of the tool
	Label    *common.I18nObject `json:"label" yaml:"label"`                   // The label of the tool
	Provider string             `json:"provider" yaml:"provider"`             // The provider of the tool
	Icon     string             `json:"icon,omitempty" yaml:"icon,omitempty"` // The optional icon of the tool
}
