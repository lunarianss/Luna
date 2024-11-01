// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dao

import (
	"github.com/lunarianss/Hurricane/internal/apiServer/model/v1"
	"github.com/lunarianss/Hurricane/internal/apiServer/model_runtime/entities"
	"github.com/lunarianss/Hurricane/internal/apiServer/model_runtime/model_providers"
	"github.com/lunarianss/Hurricane/internal/apiServer/repo"
	"github.com/lunarianss/Hurricane/internal/pkg/code"
	"github.com/lunarianss/Hurricane/pkg/errors"
	"gorm.io/gorm"
)

type ModelProviderDao struct {
	db *gorm.DB
}

var _ repo.ModelProviderRepo = (*ModelProviderDao)(nil)

func NewModelProvider(db *gorm.DB) *ModelProviderDao {
	return &ModelProviderDao{db}
}

// Get tenant's model providers
func (mpd *ModelProviderDao) GetTenantModelProviders(tenantId int64) ([]*model.Provider, error) {

	tenantProviders := []*model.Provider{}

	if err := mpd.db.Where("tenant_id = ?", tenantId).Find(&tenantProviders).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return tenantProviders, nil
}

// Get tenant's model providers mapped by provider name
func (mpd *ModelProviderDao) GetMapTenantModelProviders(tenantId int64) (map[string]*model.Provider, error) {
	providersMap := make(map[string]*model.Provider)
	tenantProviders, err := mpd.GetTenantModelProviders(tenantId)

	if err != nil {
		return nil, err
	}

	for _, tenantProvider := range tenantProviders {
		providersMap[tenantProvider.ProviderName] = tenantProvider
	}
	return providersMap, nil
}

// Get all inner Providers
func (mpd *ModelProviderDao) GetSystemProviders() ([]*entities.ProviderEntity, error) {
	return model_providers.Factory.GetProvidersFromDir()
}

// Get all inner Providers mapped by provider name
func (mpd *ModelProviderDao) GetMapSystemProviders() (map[string]*entities.ProviderEntity, error) {
	mapSystemProviders := make(map[string]*entities.ProviderEntity, model_providers.PROVIDER_COUNT)

	systemProviders, err := mpd.GetSystemProviders()

	if err != nil {
		return nil, err
	}

	for _, provider := range systemProviders {
		mapSystemProviders[provider.Provider] = provider
	}
	return mapSystemProviders, nil
}
