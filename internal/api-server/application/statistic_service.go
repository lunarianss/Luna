package service

import (
	"context"

	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
)

type StatisticService struct {
	chatDomain    *chatDomain.ChatDomain
	accountDomain *accountDomain.AccountDomain
	appDomain     *appDomain.AppDomain
}

func NewStatisticService(chatDomain *chatDomain.ChatDomain, accountDomain *accountDomain.AccountDomain, appDomain *appDomain.AppDomain) *StatisticService {
	return &StatisticService{
		chatDomain:    chatDomain,
		accountDomain: accountDomain,
		appDomain:     appDomain,
	}
}

func (ss *StatisticService) DailyMessages(ctx context.Context, appID, accountID, start, end string) (*biz_entity.StatisticDailyConversations, error) {

	accountRecord, err := ss.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, err
	}

	tenant, _, err := ss.accountDomain.GetCurrentTenantOfAccount(ctx, accountRecord.ID)

	if err != nil {
		return nil, err
	}

	appModel, err := ss.appDomain.AppRepo.GetTenantApp(ctx, appID, tenant.ID)

	if err != nil {
		return nil, err
	}

	statistics, err := ss.chatDomain.MessageRepo.StatisticDailyMessages(ctx, appModel.ID, start, end, accountRecord.Timezone)

	if err != nil {
		return nil, err
	}

	return &biz_entity.StatisticDailyConversations{
		Data: statistics,
	}, nil
}

func (ss *StatisticService) DailyConversations(ctx context.Context, appID, accountID, start, end string) (*biz_entity.StatisticDailyConversations, error) {

	accountRecord, err := ss.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, err
	}

	tenant, _, err := ss.accountDomain.GetCurrentTenantOfAccount(ctx, accountRecord.ID)

	if err != nil {
		return nil, err
	}

	appModel, err := ss.appDomain.AppRepo.GetTenantApp(ctx, appID, tenant.ID)

	if err != nil {
		return nil, err
	}

	statistics, err := ss.chatDomain.MessageRepo.StatisticDailyConversations(ctx, appModel.ID, start, end, accountRecord.Timezone)

	if err != nil {
		return nil, err
	}

	return &biz_entity.StatisticDailyConversations{
		Data: statistics,
	}, nil
}

func (ss *StatisticService) DailyUsers(ctx context.Context, appID, accountID, start, end string) (*biz_entity.StatisticDailyUser, error) {
	accountRecord, err := ss.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, err
	}

	tenant, _, err := ss.accountDomain.GetCurrentTenantOfAccount(ctx, accountRecord.ID)

	if err != nil {
		return nil, err
	}

	appModel, err := ss.appDomain.AppRepo.GetTenantApp(ctx, appID, tenant.ID)

	if err != nil {
		return nil, err
	}

	statistics, err := ss.chatDomain.MessageRepo.StatisticDailyUsers(ctx, appModel.ID, start, end, accountRecord.Timezone)

	if err != nil {
		return nil, err
	}

	return &biz_entity.StatisticDailyUser{
		Data: statistics,
	}, nil
}

func (ss *StatisticService) AverageInteractions(ctx context.Context, appID, accountID, start, end string) (*biz_entity.StatisticAverageInteraction, error) {
	accountRecord, err := ss.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, err
	}

	tenant, _, err := ss.accountDomain.GetCurrentTenantOfAccount(ctx, accountRecord.ID)

	if err != nil {
		return nil, err
	}

	appModel, err := ss.appDomain.AppRepo.GetTenantApp(ctx, appID, tenant.ID)

	if err != nil {
		return nil, err
	}

	statistics, err := ss.chatDomain.MessageRepo.StatisticAverageSessionInteraction(ctx, appModel.ID, start, end, accountRecord.Timezone)

	if err != nil {
		return nil, err
	}

	return &biz_entity.StatisticAverageInteraction{
		Data: statistics,
	}, nil
}

func (ss *StatisticService) SumTokenCosts(ctx context.Context, appID, accountID, start, end string) (*biz_entity.StatisticTokenCosts, error) {
	accountRecord, err := ss.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, err
	}

	tenant, _, err := ss.accountDomain.GetCurrentTenantOfAccount(ctx, accountRecord.ID)

	if err != nil {
		return nil, err
	}

	appModel, err := ss.appDomain.AppRepo.GetTenantApp(ctx, appID, tenant.ID)

	if err != nil {
		return nil, err
	}

	statistics, err := ss.chatDomain.MessageRepo.StatisticTokenCosts(ctx, appModel.ID, start, end, accountRecord.Timezone)

	if err != nil {
		return nil, err
	}

	return &biz_entity.StatisticTokenCosts{
		Data: statistics,
	}, nil
}
