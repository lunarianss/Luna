package biz_entity

import (
	"github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/po_entity"
)

type NotionInfo struct {
	NotionWorkspaceID string
	NotionJobID       string
	NotionPageType    string
	Document          *Document
	TenantID          string
}

type WebsiteInfo struct {
	Provider        string
	JobID           string
	Url             string
	Mode            string
	TenantID        string
	OnlyMainContent bool
}

type ExtractSetting struct {
	DatasourceType string
	UploadFile     *po_entity.UploadFile
	NotionInfo     *NotionInfo
	WebsiteInfo    *WebsiteInfo
	DocumentModel  string
}
