package po_entity

import "github.com/lunarianss/Luna/internal/infrastructure/field"

type DocumentSegment struct {
	ID            string `json:"id" gorm:"column:id"`
	TenantID      string `json:"tenant_id" gorm:"column:tenant_id"`
	DatasetID     string `json:"dataset_id" gorm:"column:dataset_id"`
	DocumentID    string `json:"document_id" gorm:"column:document_id"`
	Position      int    `json:"position" gorm:"column:position"`
	Content       string `json:"content" gorm:"column:content"`
	Answer        string `json:"answer" gorm:"column:answer"`
	WordCount     int    `json:"word_count" gorm:"column:word_count"`
	Tokens        int    `json:"tokens" gorm:"column:tokens"`
	Keywords      string `json:"keywords" gorm:"column:keywords"`
	IndexNodeID   string `json:"index_node_id" gorm:"column:index_node_id"`
	IndexNodeHash string `json:"index_node_hash" gorm:"column:index_node_hash"`
	HitCount      int    `json:"hit_count" gorm:"column:hit_count"`
	Enabled       bool   `json:"enabled" gorm:"column:enabled"`
	DisabledAt    int64  `json:"disabled_at" gorm:"column:disabled_at"`
	DisabledBy    string `json:"disabled_by" gorm:"column:disabled_by"`
	Status        string `json:"status" gorm:"column:status;default:waiting"`
	CreatedBy     string `json:"created_by" gorm:"column:created_by"`
	CreatedAt     int64  `json:"created_at" gorm:"column:created_at"`
	UpdatedBy     string `json:"updated_by" gorm:"column:updated_by"`
	UpdatedAt     int64  `json:"updated_at" gorm:"column:updated_at"`
	IndexingAt    int64  `json:"indexing_at" gorm:"column:indexing_at"`
	CompletedAt   int64  `json:"completed_at" gorm:"column:completed_at"`
	Error         string `json:"error" gorm:"column:error"`
	StoppedAt     int64  `json:"stopped_at" gorm:"column:stopped_at"`
}

func (DocumentSegment) TableName() string {
	return "document_segments"
}

type Document struct {
	ID                   string         `json:"id" gorm:"column:id"`
	TenantID             string         `json:"tenant_id" gorm:"column:tenant_id"`
	DatasetID            string         `json:"dataset_id" gorm:"column:dataset_id"`
	Position             int            `json:"position" gorm:"column:position"`
	DataSourceType       string         `json:"data_source_type" gorm:"column:data_source_type"`
	DataSourceInfo       string         `json:"data_source_info" gorm:"column:data_source_info"`
	DatasetProcessRuleID string         `json:"dataset_process_rule_id" gorm:"column:dataset_process_rule_id"`
	Batch                string         `json:"batch" gorm:"column:batch"`
	Name                 string         `json:"name" gorm:"column:name"`
	CreatedFrom          string         `json:"created_from" gorm:"column:created_from"`
	CreatedBy            string         `json:"created_by" gorm:"column:created_by"`
	CreatedAPIRequestID  string         `json:"created_api_request_id" gorm:"column:created_api_request_id"`
	CreatedAt            string         `json:"created_at" gorm:"column:created_at"`
	ProcessingStartedAt  string         `json:"processing_started_at" gorm:"column:processing_started_at"`
	FileID               string         `json:"file_id" gorm:"column:file_id"`
	WordCount            int            `json:"word_count" gorm:"column:word_count"`
	ParsingCompletedAt   string         `json:"parsing_completed_at" gorm:"column:parsing_completed_at"`
	CleaningCompletedAt  string         `json:"cleaning_completed_at" gorm:"column:cleaning_completed_at"`
	SplittingCompletedAt string         `json:"splitting_completed_at" gorm:"column:splitting_completed_at"`
	Tokens               int            `json:"tokens" gorm:"column:tokens"`
	IndexingLatency      float64        `json:"indexing_latency" gorm:"column:indexing_latency"`
	CompletedAt          string         `json:"completed_at" gorm:"column:completed_at"`
	IsPaused             field.BitBool  `json:"is_paused" gorm:"column:is_paused;default:0"`
	PausedBy             string         `json:"paused_by" gorm:"column:paused_by"`
	PausedAt             string         `json:"paused_at" gorm:"column:paused_at"`
	Error                string         `json:"error" gorm:"column:error"`
	StoppedAt            string         `json:"stopped_at" gorm:"column:stopped_at"`
	IndexingStatus       string         `json:"indexing_status" gorm:"column:indexing_status;default:'waiting'"`
	Enabled              field.BitBool  `json:"enabled" gorm:"column:enabled;default:1"`
	DisabledAt           string         `json:"disabled_at" gorm:"column:disabled_at"`
	DisabledBy           string         `json:"disabled_by" gorm:"column:disabled_by"`
	Archived             field.BitBool  `json:"archived" gorm:"column:archived;default:0"`
	ArchivedReason       string         `json:"archived_reason" gorm:"column:archived_reason"`
	ArchivedBy           string         `json:"archived_by" gorm:"column:archived_by"`
	ArchivedAt           string         `json:"archived_at" gorm:"column:archived_at"`
	UpdatedAt            string         `json:"updated_at" gorm:"column:updated_at"`
	DocType              string         `json:"doc_type" gorm:"column:doc_type"`
	DocMetadata          map[string]any `json:"doc_metadata" gorm:"column:doc_metadata;serializer:json"`
	DocForm              string         `json:"doc_form" gorm:"column:doc_form;default:'text_model'"`
	DocLanguage          string         `json:"doc_language" gorm:"column:doc_language"`
}

func (Document) TableName() string {
	return "documents"
}

type DatasetPermission struct {
	ID            string        `json:"id" gorm:"column:id"`
	DatasetID     string        `json:"dataset_id" gorm:"column:dataset_id"`
	AccountID     string        `json:"account_id" gorm:"column:account_id"`
	TenantID      string        `json:"tenant_id" gorm:"column:tenant_id"`
	HasPermission field.BitBool `json:"has_permission" gorm:"column:has_permission;default:1"`
	CreatedAt     int64         `json:"created_at" gorm:"column:created_at"`
}

func (DatasetPermission) TableName() string {
	return "dataset_permissions"
}



