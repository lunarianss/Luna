// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repo

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/entities/model_provider"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
)

type ModelProviderRepo interface {
	// UpdateProvider updates Provider by gorm updates
	UpdateProvider(ctx context.Context, provider *model.Provider) error
	// UpdateProvider updates Provider by gorm updates
	CreateProvider(ctx context.Context, provider *model.Provider) error
	// Get tenant's model providers
	GetTenantModelProviders(ctx context.Context, tenantId string) ([]*model.Provider, error)
	// Get tenant's model providers mapped by provider name
	GetMapTenantModelProviders(ctx context.Context, tenantId string) (map[string][]*model.Provider, error)
	// Get all inner Providers
	GetSystemProviders(ctx context.Context) ([]*model_provider.ProviderEntity, error)
	// Get all inner Providers mapped by provider name
	GetMapSystemProviders(ctx context.Context) (map[string]*model_provider.ProviderEntity, error)
	// Get provider path
	GetProviderPath(ctx context.Context, provider string) (string, error)
	// GerProviderEntity get the provider entity by provider name
	GetProviderEntity(ctx context.Context, provider string) (*model_provider.ProviderEntity, error)
	// GetProviderInstance get the provider entity by provider name
	GetProviderInstance(ctx context.Context, provider string) (*model_provider.ModelProvider, error)
	// GetProviders get all provider by searchProvider
	GetTenantProvider(ctx context.Context, tenant string, providerName string, providerType string) (*model.Provider, error)
}
