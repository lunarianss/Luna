package service

import (
	"context"

	accountDomain "github.com/lunarianss/Luna/internal/api-server/_domain/account/domain_service"
	appDomain "github.com/lunarianss/Luna/internal/api-server/_domain/app/domain_service"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/_domain/chat/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/_domain/provider/domain_service"
	webAppDomain "github.com/lunarianss/Luna/internal/api-server/_domain/web_app/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/config"
	"github.com/lunarianss/Luna/internal/api-server/core/app/apps/chat"
	"github.com/lunarianss/Luna/internal/api-server/core/app/apps/entities"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
)

type WebChatService struct {
	webAppDomain   *webAppDomain.WebAppDomain
	accountDomain  *accountDomain.AccountDomain
	appDomain      *appDomain.AppDomain
	chatDomain     *chatDomain.ChatDomain
	providerDomain *domain_service.ProviderDomain
	config         *config.Config
}

func NewWebChatService(webAppDomain *webAppDomain.WebAppDomain, accountDomain *accountDomain.AccountDomain, appDomain *appDomain.AppDomain, config *config.Config, providerDomain *domain_service.ProviderDomain, chatDomain *chatDomain.ChatDomain) *WebChatService {
	return &WebChatService{
		webAppDomain:   webAppDomain,
		accountDomain:  accountDomain,
		appDomain:      appDomain,
		config:         config,
		providerDomain: providerDomain,
		chatDomain:     chatDomain,
	}
}

func (s *WebChatService) Chat(ctx context.Context, appID, endUserID string, args *dto.CreateChatMessageBody, invokeFrom entities.InvokeForm, streaming bool) error {

	appModel, err := s.appDomain.AppRepo.GetAppByID(ctx, appID)

	if err != nil {
		return err
	}

	endUserRecord, err := s.webAppDomain.WebAppRepo.GetEndUserByID(ctx, endUserID)

	if err != nil {
		return err
	}

	chatAppGenerator := chat.NewChatAppGenerator(s.appDomain, s.providerDomain, s.chatDomain)

	if err := chatAppGenerator.Generate(ctx, appModel, endUserRecord, args, invokeFrom, true); err != nil {
		return err
	}
	return nil
}
