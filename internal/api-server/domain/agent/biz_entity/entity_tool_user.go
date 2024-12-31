package biz_entity

import common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"

type UserToolProvider struct {
	ID                  string                 `json:"id"`                              // The provider ID
	Author              string                 `json:"author"`                          // The provider author
	Name                string                 `json:"name" `                           // Identifier
	Description         *common.I18nObject     `json:"description" `                    // Description
	Icon                string                 `json:"icon" `                           // Icon path or URL
	Label               *common.I18nObject     `json:"label" `                          // Label
	Type                ToolProviderType       `json:"type" `                           // Provider type
	MaskedCredentials   map[string]interface{} `json:"team_credentials" `               // Masked credentials
	OriginalCredentials map[string]interface{} `json:"original_credentials,omitempty" ` // Original credentials
	IsTeamAuthorization bool                   `json:"is_team_authorization" `          // Is team authorization required
	AllowDelete         bool                   `json:"allow_delete" `                   // Allow deletion
	Tools               []*UserTool            `json:"tools"`                           // List of tools
	Labels              []string               `json:"labels"`                          // List of labels
}
