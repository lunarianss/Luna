package model_provider

import "github.com/lunarianss/Hurricane/internal/api-server/repo"

type ModelProviderDomain struct {
	ModelProviderRepo repo.ModelProviderRepo
}

func NewBlogDomain(modelProviderRepo repo.ModelProviderRepo) *ModelProviderDomain {
	return &ModelProviderDomain{
		ModelProviderRepo: modelProviderRepo,
	}
}
