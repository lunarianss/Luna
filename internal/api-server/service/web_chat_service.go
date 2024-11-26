package service

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/config"
	"github.com/lunarianss/Luna/internal/api-server/core/app/apps/chat"
	"github.com/lunarianss/Luna/internal/api-server/core/app/apps/entities"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/app_running"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat"
	providerDomain "github.com/lunarianss/Luna/internal/api-server/domain/provider"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
)

type WebChatService struct {
	appRunningDomain *domain.AppRunningDomain
	accountDomain    *accountDomain.AccountDomain
	appDomain        *appDomain.AppDomain
	chatDomain       *chatDomain.ChatDomain
	providerDomain   *providerDomain.ModelProviderDomain
	config           *config.Config
}

func NewWebChatService(appRunningDomain *domain.AppRunningDomain, accountDomain *accountDomain.AccountDomain, appDomain *appDomain.AppDomain, config *config.Config, providerDomain *providerDomain.ModelProviderDomain, chatDomain *chatDomain.ChatDomain) *WebChatService {
	return &WebChatService{
		appRunningDomain: appRunningDomain,
		accountDomain:    accountDomain,
		appDomain:        appDomain,
		config:           config,
		providerDomain:   providerDomain,
		chatDomain:       chatDomain,
	}
}

func (s *WebChatService) Chat(ctx context.Context, appID, endUserID string, args *dto.CreateChatMessageBody, invokeFrom entities.InvokeForm, streaming bool) error {

	appModel, err := s.appDomain.AppRepo.GetAppByID(ctx, appID)

	if err != nil {
		return err
	}

	endUserRecord, err := s.appRunningDomain.AppRunningRepo.GetEndUserByID(ctx, endUserID)

	if err != nil {
		return err
	}

	chatAppGenerator := chat.NewChatAppGenerator(s.appDomain, s.providerDomain, s.chatDomain)

	if err := chatAppGenerator.Generate(ctx, appModel, endUserRecord, args, invokeFrom, true); err != nil {
		return err
	}
	return nil
}