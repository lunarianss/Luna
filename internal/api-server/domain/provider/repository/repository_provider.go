// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"context"

	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/po_entity"
)

type ProviderRepo interface {
	// UpdateProvider updates Provider by gorm updates
	UpdateProvider(ctx context.Context, provider *po_entity.Provider) error
	// UpdateProvider updates Provider by gorm updates
	CreateProvider(ctx context.Context, provider *po_entity.Provider) error
	// Get tenant's model providers
	GetTenantModelProviders(ctx context.Context, tenantId string) ([]*po_entity.Provider, error)
	// Get tenant's model providers mapped by provider name
	GetMapTenantModelProviders(ctx context.Context, tenantId string) (map[string][]*po_entity.Provider, error)
	// Get all inner Providers
	GetSystemProviders(ctx context.Context) ([]*biz_entity.ProviderStaticConfiguration, error)
	// Get all inner Providers mapped by provider name
	GetMapSystemProviders(ctx context.Context) (map[string]*biz_entity.ProviderStaticConfiguration, error)
	// Get provider path
	GetProviderPath(ctx context.Context, provider string) (string, error)
	// GerProviderEntity get the provider entity by provider name
	GetProviderEntity(ctx context.Context, provider string) (*biz_entity.ProviderStaticConfiguration, error)
	// GetProviderInstance get the provider entity by provider name
	GetProviderInstance(ctx context.Context, provider string) (*biz_entity.ProviderRuntime, error)
	// GetProviders get all provider by searchProvider
	GetTenantProvider(ctx context.Context, tenant string, providerName string, providerType string) (*po_entity.Provider, error)
}
