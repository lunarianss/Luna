package service

import (
	"context"

	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	biz_entity_console_app_statistic "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/console_app_statistic"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
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

func (ss *StatisticService) DailyMessages(ctx context.Context, appID, accountID, start, end string) (*biz_entity_console_app_statistic.StatisticDailyConversations, error) {

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

	return &biz_entity_console_app_statistic.StatisticDailyConversations{
		Data: statistics,
	}, nil
}

func (ss *StatisticService) DailyConversations(ctx context.Context, appID, accountID, start, end string) (*biz_entity_console_app_statistic.StatisticDailyConversations, error) {

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

	return &biz_entity_console_app_statistic.StatisticDailyConversations{
		Data: statistics,
	}, nil
}

func (ss *StatisticService) DailyUsers(ctx context.Context, appID, accountID, start, end string) (*biz_entity_console_app_statistic.StatisticDailyUser, error) {
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

	return &biz_entity_console_app_statistic.StatisticDailyUser{
		Data: statistics,
	}, nil
}

func (ss *StatisticService) AverageInteractions(ctx context.Context, appID, accountID, start, end string) (*biz_entity_console_app_statistic.StatisticAverageInteraction, error) {
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

	return &biz_entity_console_app_statistic.StatisticAverageInteraction{
		Data: statistics,
	}, nil
}

func (ss *StatisticService) SumTokenCosts(ctx context.Context, appID, accountID, start, end string) (*biz_entity_console_app_statistic.StatisticTokenCosts, error) {
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

	return &biz_entity_console_app_statistic.StatisticTokenCosts{
		Data: statistics,
	}, nil
}
