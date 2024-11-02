// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dao

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/lunarianss/Luna/internal/api-server/model-runtime/entities"
	model_providers "github.com/lunarianss/Luna/internal/api-server/model-runtime/model-providers"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/repo"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
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

	if err := mpd.db.Where("tenant_id = ? and is_valid = ?", tenantId, 1).Find(&tenantProviders).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return tenantProviders, nil
}

// Get tenant's model providers mapped by provider name
func (mpd *ModelProviderDao) GetMapTenantModelProviders(tenantId int64) (map[string][]*model.Provider, error) {
	providersMap := make(map[string][]*model.Provider)
	tenantProviders, err := mpd.GetTenantModelProviders(tenantId)

	if err != nil {
		return nil, err
	}

	for _, tenantProvider := range tenantProviders {
		providersMap[tenantProvider.ProviderName] = append(providersMap[tenantProvider.ProviderName], tenantProvider)
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

func (mpd *ModelProviderDao) GetProviderPath(provider string) (string, error) {
	providerPath, err := model_providers.Factory.ResolveProviderDirPath()

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", providerPath, provider), nil
}

func (mpd *ModelProviderDao) GetProviderEntity(provider string) (*entities.ProviderEntity, error) {
	modelProvider, err := model_providers.Factory.GetProviderInstance(provider)

	if err != nil {
		return nil, err
	}

	providerEntity, err := modelProvider.GetProviderSchema()

	if err != nil {
		return nil, err
	}
	return providerEntity, nil
}
