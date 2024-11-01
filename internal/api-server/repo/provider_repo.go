// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repo

import (
	"github.com/lunarianss/Luna/internal/api-server/model-runtime/entities"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
)

type ModelProviderRepo interface {
	// Get tenant's model providers
	GetTenantModelProviders(tenantId int64) ([]*model.Provider, error)
	// Get tenant's model providers mapped by provider name
	GetMapTenantModelProviders(tenantId int64) (map[string][]*model.Provider, error)
	// Get all inner Providers
	GetSystemProviders() ([]*entities.ProviderEntity, error)
	// Get all inner Providers mapped by provider name
	GetMapSystemProviders() (map[string]*entities.ProviderEntity, error)
}
