package service

import (
	"context"

	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
)

type StatisticService struct {
	chatDomain    *chatDomain.ChatDomain
	accountDomain *accountDomain.AccountDomain
}

func NewStatisticService(chatDomain *chatDomain.ChatDomain, accountDomain *accountDomain.AccountDomain) *StatisticService {
	return &StatisticService{
		chatDomain:    chatDomain,
		accountDomain: accountDomain,
	}
}

func (ss *StatisticService) DailyConversations(ctx context.Context, appID, accountID, start, end string) (*biz_entity.StatisticDailyConversations, error) {

	account, err := ss.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, err
	}

	statistics, err := ss.chatDomain.MessageRepo.StatisticDailyConversations(ctx, appID, start, end, account.Timezone)

	if err != nil {
		return nil, err
	}

	return &biz_entity.StatisticDailyConversations{
		Data: statistics,
	}, nil
}

func (ss *StatisticService) DailyUsers(ctx context.Context, appID, accountID, start, end string) (*biz_entity.StatisticDailyUser, error) {

	account, err := ss.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, err
	}

	statistics, err := ss.chatDomain.MessageRepo.StatisticDailyUsers(ctx, appID, start, end, account.Timezone)

	if err != nil {
		return nil, err
	}

	return &biz_entity.StatisticDailyUser{
		Data: statistics,
	}, nil
}
