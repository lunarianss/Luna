package biz_entity

type ToolProviderRuntime struct {
	*ToolProviderStatic
	Tools            []*ToolRuntimeConfiguration `json:"tools" yaml:"tools"`
	ConfPath         string                      `json:"conf_path"`
	ToolProviderName string                      `json:"tool_provider_name"`
}

func (tr *ToolProviderRuntime) GetToolLabels() []string {
	var categoryNames []string

	for _, tag := range tr.Identity.Tags {
		categoryNames = append(categoryNames, string(tag))
	}

	return categoryNames
}

func (tr *ToolProviderRuntime) NeedCredentials() bool {
	return len(tr.CredentialsSchema) != 0
}

func (tr *ToolProviderRuntime) GetTools() []*ToolRuntimeConfiguration {
	return tr.Tools
}

func (tr *ToolProviderRuntime) GetTool(toolName string) *ToolRuntimeConfiguration {
	for _, tool := range tr.GetTools() {
		if tool.Identity.Name == toolName {
			return tool
		}
	}

	return nil
}
