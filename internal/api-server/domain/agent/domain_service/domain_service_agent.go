package domain_service

import (
	"context"
	"strings"
	"sync"

	"github.com/lunarianss/Luna/internal/api-server/core/tools"
	"github.com/lunarianss/Luna/internal/api-server/domain/agent/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/agent/repository"
	repo_app "github.com/lunarianss/Luna/internal/api-server/domain/app/repository"
	po_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
)

type AgentDomain struct {
	*ToolTransformService
	*tools.ToolManager
	*sync.RWMutex

	repository.AgentRepo
	repo_app.AppRepo
}

func NewAgentDomain(ts *ToolTransformService, tm *tools.ToolManager, agentRepo repository.AgentRepo, appRepo repo_app.AppRepo) *AgentDomain {
	return &AgentDomain{
		ToolTransformService: ts,
		ToolManager:          tm,
		RWMutex:              &sync.RWMutex{},
		AgentRepo:            agentRepo,
		AppRepo:              appRepo,
	}
}

func (ad *AgentDomain) BuildMessageFile(ctx context.Context, msg *po_chat.Message, tenantID string, baseUrl, secretKey string) ([]*biz_entity.BuildFile, error) {

	var result []*biz_entity.BuildFile

	messageFiles, err := ad.AgentRepo.GetMessageFileByMessage(ctx, msg.ID)

	if err != nil {
		return nil, err
	}

	for _, messageFile := range messageFiles {
		if messageFile.TransferMethod == "tool_file" {

			urls := strings.Split(messageFile.URL, "/")
			toolFileName := strings.Split(urls[len(urls)-1], ".")

			toolFileID := toolFileName[0]
			extension := toolFileName[1]
			toolFile, err := ad.AgentRepo.GetToolFileByID(ctx, toolFileID)

			if err != nil {
				return nil, err
			}

			toolFileManager := ToolFileManager{}
			url, err := toolFileManager.SignFile(toolFileID, extension, secretKey, baseUrl)

			if err != nil {
				return nil, err
			}

			result = append(result, &biz_entity.BuildFile{
				ID:             messageFile.ID,
				Type:           messageFile.Type,
				TenantID:       tenantID,
				Filename:       toolFile.Name,
				TransferMethod: messageFile.TransferMethod,
				RemoteUrl:      toolFile.OriginalURL,
				RelatedID:      toolFile.ID,
				Extension:      extension,
				MimeType:       toolFile.MimeType,
				Size:           int64(toolFile.Size),
				ExtraConfig:    nil,
				BelongsTo:      messageFile.BelongsTo,
				Url:            url,
			})
		}
	}

	return result, nil

}

func (ad *AgentDomain) ListBuiltInTools(ctx context.Context, tenantID string) ([]*biz_entity.UserToolProvider, error) {

	var result []*biz_entity.UserToolProvider

	systemToolProviders, err := ad.ListBuiltInProviders()

	if err != nil {
		return nil, err
	}

	for _, systemToolProvider := range systemToolProviders {

		userProvider := ad.BuiltInProviderToUserProvider(systemToolProvider, nil, true)

		for _, tool := range systemToolProvider.Tools {
			userProvider.Tools = append(userProvider.Tools, ad.BuiltInToolToUserTool(tool, nil, tenantID, userProvider.Labels))
		}

		util.PatchI18nObject(userProvider)
		result = append(result, userProvider)
	}

	return result, nil
}

func (ad *AgentDomain) ListBuiltInLabels(ctx context.Context) ([]*biz_entity.ToolLabel, error) {
	ad.RLock()
	defer ad.RUnlock()

	return biz_entity.GetDefaultTools(), nil
}
