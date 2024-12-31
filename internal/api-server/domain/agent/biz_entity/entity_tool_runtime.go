package biz_entity

type ToolRuntimeConfiguration struct {
	*ToolStaticConfiguration
	TenantID          string         `json:"tenant_id"`
	ToolID            string         `json:"tool_id"`
	InvokeFrom        InvokeFrom     `json:"invoke_from"`
	ToolInvokeFrom    ToolInvokeFrom `json:"tool_invoke_from" `
	Credentials       map[string]any `json:"credentials" `
	RuntimeParameters map[string]any `json:"runtime_parameters"`
	ConfPath          string         `json:"conf_path"`
}
