package biz_entity

import "github.com/lunarianss/Luna/internal/api-server/domain/agent/entity/po_entity"

type MessageType string

var (
	TEXT       MessageType = "text"
	IMAGE      MessageType = "image"
	LINK       MessageType = "link"
	BLOB       MessageType = "blob"
	JSON       MessageType = "json"
	IMAGE_LINK MessageType = "image_link"
	FILE       MessageType = "file"
)

type ToolInvokeMessage struct {
	Type    MessageType    `json:"type"`
	Message any            `json:"message"`
	Meta    map[string]any `json:"meta"`
	SaveAs  string         `json:"save_as"`
}

type ToolEngineInvokeMeta struct {
	TimeCost   float64        `json:"time_cost"`
	Error      string         `json:"error"`
	ToolConfig map[string]any `json:"tool_config"`
}

func ConvertToPoMeta(meta *ToolEngineInvokeMeta) *po_entity.ToolEngineInvokeMeta {
	return &po_entity.ToolEngineInvokeMeta{
		TimeCost:   meta.TimeCost,
		Error:      meta.Error,
		ToolConfig: meta.ToolConfig,
	}
}

func ErrorInvokeMetaIns(err string) *ToolEngineInvokeMeta {
	return &ToolEngineInvokeMeta{
		Error:      err,
		TimeCost:   0,
		ToolConfig: make(map[string]any),
	}
}

type ToolEngineInvokeMessageFiles struct {
	MessageFile string `json:"message_file"`
	SaveAs      string `json:"save_as"`
}

type ToolEngineInvokeMessage struct {
	InvokeToolPrompt string                          `json:"invoke_tool_prompt"`
	MessageFiles     []*ToolEngineInvokeMessageFiles `json:"message_files"`
	ToolInvokeMeta   *ToolEngineInvokeMeta           `json:"tool_invoke_meta"`
}

type ToolArtifact struct {
	ToolCallID   string                `json:"tool_call_id"`
	ToolCallName string                `json:"tool_call_name"`
	ToolResponse string                `json:"tool_response"`
	Meta         *ToolEngineInvokeMeta `json:"meta"`
}
