package biz_entity

type VariableKey string

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

func (tc *ToolRuntimeConfiguration) GetAllRuntimeParameters() []*ToolParameter {
	return tc.Parameters
}

func (tc *ToolRuntimeConfiguration) CreateBlobMessage(blob []byte, meta map[string]any, saveAs string) *ToolInvokeMessage {
	return &ToolInvokeMessage{
		Type:    BLOB,
		Message: blob,
		Meta:    meta,
		SaveAs:  saveAs,
	}
}
