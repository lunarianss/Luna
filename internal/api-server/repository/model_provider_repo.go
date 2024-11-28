// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repo_impl

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/repository"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/internal/pkg/mysql"
	"github.com/lunarianss/Luna/pkg/errors"
	"gorm.io/gorm"
)

type ModelProviderRepoImpl struct {
	db *gorm.DB
}

var _ repository.ModelRepo = (*ModelProviderRepoImpl)(nil)

func NewModelProviderRepoImpl(db *gorm.DB) repository.ModelRepo {
	return &ModelProviderRepoImpl{db}
}

func (md *ModelProviderRepoImpl) GetTenantModel(ctx context.Context, tenantId, providerName, modelName, modelType string) (*po_entity.ProviderModel, error) {
	var model *po_entity.ProviderModel

	if err := md.db.Scopes(mysql.IDDesc()).Where("tenant_id = ? and provider_name = ? and model_name = ? and model_type = ?", tenantId, providerName, modelName, modelType).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, errors.WithCode(code.ErrDatabase, err.Error())
		}
	}
	return model, nil
}

func (md *ModelProviderRepoImpl) UpdateModel(ctx context.Context, model *po_entity.ProviderModel) error {
	if err := md.db.Updates(model).Error; err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil
}

func (md *ModelProviderRepoImpl) CreateModel(ctx context.Context, model *po_entity.ProviderModel) error {
	if err := md.db.Create(model).Error; err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil
}

func (md *ModelProviderRepoImpl) GetTenantDefaultModel(ctx context.Context, tenantID, modelType string) (*po_entity.TenantDefaultModel, error) {
	var defaultModel po_entity.TenantDefaultModel
	if err := md.db.Scopes(mysql.IDDesc()).Where("tenant_id = ? and model_type = ?", tenantID, modelType).First(&defaultModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, errors.WithCode(code.ErrDatabase, err.Error())
		}
	}
	return &defaultModel, nil
}

func (md *ModelProviderRepoImpl) CreateTenantDefaultModel(ctx context.Context, tenantDefaultModel *po_entity.TenantDefaultModel) (*po_entity.TenantDefaultModel, error) {
	if err := md.db.Create(tenantDefaultModel).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return tenantDefaultModel, nil
}
