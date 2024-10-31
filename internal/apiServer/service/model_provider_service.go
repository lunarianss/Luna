package service

import domain "github.com/lunarianss/Hurricane/internal/apiServer/domain/model_provider"

type ModelProviderService struct {
	ModelProviderDomain domain.ModelProviderDomain
}

func NewModelProviderService(modelProviderDomain domain.ModelProviderDomain) *ModelProviderService {
	return &ModelProviderService{ModelProviderDomain: modelProviderDomain}
}

// 1. 数据库获取租户的 provider map[string]Provider
// 2. 配置文件中获取所有 provider entities
// 3. 获取租户工作空间的 configurations （区分订阅/未订阅）
func (mpSrv *ModelProviderService) GetProviderList(tenant_id string, model_type string) (interface{}, error) {

    mpSrv.ModelProviderDomain.ModelProviderRepo.Get

}
