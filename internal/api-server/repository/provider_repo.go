// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repo_impl

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/lunarianss/Luna/infrastructure/errors"
	model_providers "github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_providers"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/repository"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/lunarianss/Luna/internal/infrastructure/mysql"
)

type ProviderRepoImpl struct {
	db *gorm.DB
}

var _ repository.ProviderRepo = (*ProviderRepoImpl)(nil)

func NewProviderRepoImpl(db *gorm.DB) *ProviderRepoImpl {
	return &ProviderRepoImpl{db}
}

func (mpd *ProviderRepoImpl) GetTenantProvider(ctx context.Context, tenantId string, providerName string, providerType string) (*po_entity.Provider, error) {
	var provider *po_entity.Provider

	if err := mpd.db.Scopes(mysql.IDDesc()).Where("tenant_id = ? and provider_name = ? and provider_type = ?", tenantId, providerName, providerType).First(&provider).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, errors.WithSCode(code.ErrDatabase, err.Error())
		}
	}
	return provider, nil
}

// Get tenant's model providers
func (mpd *ProviderRepoImpl) GetTenantModelProviders(ctx context.Context, tenantId string) ([]*po_entity.Provider, error) {

	var tenantProviders []*po_entity.Provider

	if err := mpd.db.Where("tenant_id = ? and is_valid = ?", tenantId, 1).Find(&tenantProviders).Error; err != nil {
		return nil, errors.WithSCode(code.ErrDatabase, err.Error())
	}
	return tenantProviders, nil
}

// Get tenant's model providers mapped by provider name
func (mpd *ProviderRepoImpl) GetMapTenantModelProviders(ctx context.Context, tenantId string) (map[string][]*po_entity.Provider, error) {
	providersMap := make(map[string][]*po_entity.Provider)
	tenantProviders, err := mpd.GetTenantModelProviders(ctx, tenantId)

	if err != nil {
		return nil, err
	}

	for _, tenantProvider := range tenantProviders {
		providersMap[tenantProvider.ProviderName] = append(providersMap[tenantProvider.ProviderName], tenantProvider)
	}
	return providersMap, nil
}

// Get all inner Providers
func (mpd *ProviderRepoImpl) GetSystemProviders(ctx context.Context) ([]*biz_entity.ProviderStaticConfiguration, []string, error) {
	return model_providers.Factory.GetProvidersFromDir()
}

// Get all inner Providers mapped by provider name
func (mpd *ProviderRepoImpl) GetMapSystemProviders(ctx context.Context) (map[string]*biz_entity.ProviderStaticConfiguration, []string, error) {
	mapSystemProviders := make(map[string]*biz_entity.ProviderStaticConfiguration, model_providers.PROVIDER_COUNT)

	systemProviders, orderedProviders, err := mpd.GetSystemProviders(ctx)

	if err != nil {
		return nil, nil, err
	}

	for _, provider := range systemProviders {
		mapSystemProviders[provider.Provider] = provider
	}

	return mapSystemProviders, orderedProviders, nil
}

func (mpd *ProviderRepoImpl) GetProviderPath(ctx context.Context, provider string) (string, error) {
	providerPath, err := model_providers.Factory.ResolveProviderDirPath()

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", providerPath, provider), nil
}

func (mpd *ProviderRepoImpl) GetProviderEntity(ctx context.Context, provider string) (*biz_entity.ProviderStaticConfiguration, error) {
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

func (mpd *ProviderRepoImpl) UpdateProvider(ctx context.Context, provider *po_entity.Provider) error {
	if err := mpd.db.Updates(provider).Error; err != nil {
		return errors.WithSCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (mpd *ProviderRepoImpl) CreateProvider(ctx context.Context, provider *po_entity.Provider) error {
	if err := mpd.db.Create(provider).Error; err != nil {
		return errors.WithSCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (mpd *ProviderRepoImpl) GetProviderInstance(ctx context.Context, provider string) (*biz_entity.ProviderRuntime, error) {
	return model_providers.Factory.GetProviderInstance(provider)
}
