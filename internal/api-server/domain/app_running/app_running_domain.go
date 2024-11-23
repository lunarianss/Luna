package domain

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/repo"
)

type AppRunningDomain struct {
	AppRunningRepo repo.AppRunningRepo
}

func NewAppRunningDomain(appRunningRepo repo.AppRunningRepo) *AppRunningDomain {
	return &AppRunningDomain{
		AppRunningRepo: appRunningRepo,
	}
}

func (ard *AppRunningDomain) CreateEndUser(ctx context.Context, appRecord *model.App) (*model.EndUser, error) {
	sessionID, err := ard.AppRunningRepo.GenerateSessionForEndUser(ctx)

	if err != nil {
		return nil, err
	}

	endUser := &model.EndUser{
		TenantID:  appRecord.TenantID,
		AppID:     appRecord.ID,
		Type:      "browser",
		SessionID: sessionID,
	}

	endUserRecord, err := ard.AppRunningRepo.CreateEndUser(ctx, endUser, nil)

	if err != nil {
		return nil, err
	}

	return endUserRecord, nil
}
