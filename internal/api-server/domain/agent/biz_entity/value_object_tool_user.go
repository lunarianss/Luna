package biz_entity

import common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"

type UserToolProviderTypeLiteral string

const (
	Builtin  UserToolProviderTypeLiteral = "builtin"
	API      UserToolProviderTypeLiteral = "api"
	Workflow UserToolProviderTypeLiteral = "workflow"
)

// ToolProviderType defines the types of tool providers.
type ToolProviderType string

const (
	// Built-in tool provider type
	ToolProviderTypeBuiltIn ToolProviderType = "builtin"
	// Workflow tool provider type
	ToolProviderTypeWorkflow ToolProviderType = "workflow"
	// API tool provider type
	ToolProviderTypeAPI ToolProviderType = "api"
	// App tool provider type
	ToolProviderTypeApp ToolProviderType = "app"
	// Dataset retrieval tool provider type
	ToolProviderTypeDatasetRetrieval ToolProviderType = "dataset-retrieval"
)

type UserTool struct {
	Author      string             `json:"author"`      // The author of the tool
	Name        string             `json:"name"`        // Identifier
	Label       *common.I18nObject `json:"label"`       // Label
	Description *common.I18nObject `json:"description"` // Description
	Parameters  []ToolParameter    `json:"parameters"`  // List of tool parameters
	Labels      []string           `json:"labels"`      // List of labels
}
