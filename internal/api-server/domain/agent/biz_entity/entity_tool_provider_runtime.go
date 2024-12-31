package biz_entity

type ToolProviderRuntime struct {
	*ToolProviderStatic
	ConfPath         string `json:"conf_path"`
	ToolProviderName string `json:"tool_provider_name"`
}
