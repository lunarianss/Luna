package dto

type ListIconUri struct {
	Provider string `json:"provider" uri:"provider"`
}

type MessageAgentThought struct {
	ID             string            `json:"id"`
	MessageID      string            `json:"message_id"`
	MessageChainID string            `json:"chain_id"`
	Position       int               `json:"position"`
	Thought        string            `json:"thought"`
	Tool           string            `json:"tool"`
	ToolLabelsStr  string            `json:"tool_labels"`
	ToolInput      map[string]string `json:"tool_input"`
	Observation    map[string]string `json:"observation"`
	Message        string            `json:"message"`
	MessageFiles   []string          `json:"files"`
}

type MessageFile struct {
	ID             string `json:"id"`
	Type           string `json:"type"`
	TransferMethod string `json:"transfer_method"`
	URL            string `json:"url"`
	BelongsTo      string `json:"belongs_to"`
	MimeType       string `json:"mime_type"`
	Size           int64  `json:"size"`
	FileName       string `json:"filename"`
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
