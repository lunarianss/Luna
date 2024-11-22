package app

import "github.com/lunarianss/Luna/internal/api-server/repo"

type AppDomain struct {
	AppRepo        repo.AppRepo
	AppRunningRepo repo.AppRunningRepo
}

func NewAppDomain(appRepo repo.AppRepo, appRunningRepo repo.AppRunningRepo) *AppDomain {
	return &AppDomain{
		AppRepo:        appRepo,
		AppRunningRepo: appRunningRepo,
	}
}
