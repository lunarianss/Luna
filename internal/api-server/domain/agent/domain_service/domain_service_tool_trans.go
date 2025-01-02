package domain_service

import (
	"fmt"

	"github.com/lunarianss/Luna/internal/api-server/config"
	"github.com/lunarianss/Luna/internal/api-server/domain/agent/biz_entity"
	common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"
)

type ToolTransformService struct {
	config *config.Config
}

func NewToolTransformService(config *config.Config) *ToolTransformService {
	return &ToolTransformService{
		config: config,
	}
}

func (tf *ToolTransformService) BuiltInProviderToUserProvider(providerRuntime *biz_entity.ToolProviderRuntime, dbProvider any, decryptCredentials bool) *biz_entity.UserToolProvider {
	result := &biz_entity.UserToolProvider{
		ID:     providerRuntime.Identity.Name,
		Author: providerRuntime.Identity.Author,
		Name:   providerRuntime.Identity.Name,
		Description: &common.I18nObject{
			En_US:   providerRuntime.Identity.Description.En_US,
			Zh_Hans: providerRuntime.Identity.Description.Zh_Hans,
		},
		Icon: fmt.Sprintf("%s/%s/%s/icon", tf.config.SystemOptions.IconBaseUrl, "v1/console/api/workspaces/current/tool-provider/builtin", providerRuntime.ToolProviderName),
		Label: &common.I18nObject{
			En_US:   providerRuntime.Identity.Label.En_US,
			Zh_Hans: providerRuntime.Identity.Label.Zh_Hans,
		},
		Type:                biz_entity.ToolProviderTypeBuiltIn,
		MaskedCredentials:   make(map[string]interface{}, 0),
		IsTeamAuthorization: false,
		Tools:               make([]*biz_entity.UserTool, 0),
		Labels:              providerRuntime.GetToolLabels(),
		// AllowDelete:         true,
	}

	for key := range providerRuntime.CredentialsSchema {
		result.MaskedCredentials[key] = ""
	}

	if !providerRuntime.NeedCredentials() {
		result.IsTeamAuthorization = true
		result.AllowDelete = false
	}

	return result
}

func (tf *ToolTransformService) BuiltInToolToUserTool(tool *biz_entity.ToolRuntimeConfiguration, credentials map[string]any, tenantID string, labels []string) *biz_entity.UserTool {
	return &biz_entity.UserTool{
		Author:      tool.Identity.Author,
		Name:        tool.Identity.Name,
		Label:       tool.Identity.Label,
		Description: tool.Description.Human,
		Parameters:  tool.Parameters,
		Labels:      labels,
	}
}
