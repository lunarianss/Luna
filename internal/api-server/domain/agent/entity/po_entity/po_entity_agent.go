package po_entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ToolEngineInvokeMeta struct {
	TimeCost   float64        `json:"time_cost"`
	Error      string         `json:"error"`
	ToolConfig map[string]any `json:"tool_config"`
}

type MessageAgentThought struct {
	ID               string                           `json:"id" gorm:"column:id"`
	MessageID        string                           `json:"message_id" gorm:"column:message_id"`
	MessageChainID   string                           `json:"message_chain_id" gorm:"column:message_chain_id"`
	Position         int                              `json:"position" gorm:"column:position"`
	Thought          string                           `json:"thought" gorm:"column:thought"`
	Tool             string                           `json:"tool" gorm:"column:tool"`
	ToolLabelsStr    string                           `json:"tool_labels_str" gorm:"column:tool_labels_str;default:{}"`
	ToolMetaStr      map[string]*ToolEngineInvokeMeta `json:"tool_meta_str" gorm:"column:tool_meta_str;default:{};serializer:json"`
	ToolInput        map[string]map[string]any        `json:"tool_input" gorm:"column:tool_input;serializer:json"`
	Observation      map[string]string                `json:"observation" gorm:"column:observation;serializer:json"`
	ToolProcessData  string                           `json:"tool_process_data" gorm:"column:tool_process_data"`
	Message          string                           `json:"message" gorm:"column:message"`
	MessageToken     int                              `json:"message_token" gorm:"column:message_token"`
	MessageUnitPrice float64                          `json:"message_unit_price" gorm:"column:message_unit_price"`
	MessagePriceUnit float64                          `json:"message_price_unit" gorm:"column:message_price_unit;default:0.001"`
	MessageFiles     []string                         `json:"message_files" gorm:"column:message_files;serializer:json"`
	Answer           string                           `json:"answer" gorm:"column:answer"`
	AnswerToken      int                              `json:"answer_token" gorm:"column:answer_token"`
	AnswerUnitPrice  float64                          `json:"answer_unit_price" gorm:"column:answer_unit_price"`
	AnswerPriceUnit  float64                          `json:"answer_price_unit" gorm:"column:answer_price_unit;default:0.001"`
	Tokens           int                              `json:"tokens" gorm:"column:tokens"`
	TotalPrice       float64                          `json:"total_price" gorm:"column:total_price"`
	Currency         string                           `json:"currency" gorm:"column:currency"`
	Latency          float64                          `json:"latency" gorm:"column:latency"`
	CreatedByRole    string                           `json:"created_by_role" gorm:"column:created_by_role"`
	CreatedBy        string                           `json:"created_by" gorm:"column:created_by"`
	CreatedAt        int64                            `json:"created_at" gorm:"column:created_at"`
}

func (*MessageAgentThought) TableName() string {
	return "message_agent_thoughts"
}

func (m *MessageAgentThought) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.NewString()
	return
}

func (m *MessageAgentThought) BeforeSave(tx *gorm.DB) (err error) {
	// 初始化 map 字段，确保它们有默认值
	if m.ToolMetaStr == nil {
		m.ToolMetaStr = make(map[string]*ToolEngineInvokeMeta)
	}
	if m.ToolInput == nil {
		m.ToolInput = make(map[string]map[string]any)
	}
	if m.Observation == nil {
		m.Observation = make(map[string]string)
	}
	return nil
}
