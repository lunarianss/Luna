package app_generate

import "github.com/lunarianss/Luna/internal/api-server/repo"

type AppGenerateDomain struct {
	AppRepo repo.AppRepo
}

func NewAppDomain(appRepo repo.AppRepo) *AppGenerateDomain {
	return &AppGenerateDomain{
		AppRepo: appRepo,
	}
}
