package app

import "github.com/lunarianss/Luna/internal/api-server/repo"

type AppDomain struct {
	AppRepo        repo.AppRepo
	AppRunningRepo repo.AppRunningRepo
	MessageRepo    repo.MessageRepo
}

func NewAppDomain(appRepo repo.AppRepo, appRunningRepo repo.AppRunningRepo, messageRepo repo.MessageRepo) *AppDomain {
	return &AppDomain{
		AppRepo:        appRepo,
		AppRunningRepo: appRunningRepo,
		MessageRepo:    messageRepo,
	}
}
