package app

import "github.com/lunarianss/Luna/internal/api-server/repo"

type AppDomain struct {
	AppRepo repo.AppRepo
}

func NewAppDomain(appRepo repo.AppRepo) *AppDomain {
	return &AppDomain{
		AppRepo: appRepo,
	}
}
