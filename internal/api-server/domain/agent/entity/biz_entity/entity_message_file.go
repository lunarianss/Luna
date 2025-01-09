package biz_entity

type ToolFileMapping struct {
	ID             string
	Type           string
	TransferMethod string
	ToolFileID     string
}

type BuildFile struct {
	ID             string      `json:"id,omitempty"`
	TenantID       string      `json:"tenant_id,omitempty"`
	Type           string      `json:"type,omitempty"`
	TransferMethod string      `json:"transfer_method,omitempty"`
	RemoteUrl      string      `json:"remote_url,omitempty"`
	RelatedID      string      `json:"related_id,omitempty"`
	Filename       string      `json:"filename,omitempty"`
	Extension      string      `json:"extension,omitempty"`
	MimeType       string      `json:"mime_type,omitempty"`
	Size           int64       `json:"size"`
	ExtraConfig    interface{} `json:"extra_config,omitempty"`
	BelongsTo      string      `json:"belongs_to"`
	Url            string      `json:"url"`
}
