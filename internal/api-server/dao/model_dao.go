package dao

import (
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/repo"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/internal/pkg/mysql"
	"github.com/lunarianss/Luna/pkg/errors"
	"gorm.io/gorm"
)

type ModelDao struct {
	db *gorm.DB
}

var _ repo.ModelRepo = (*ModelDao)(nil)

func NewModelDao(db *gorm.DB) *ModelDao {
	return &ModelDao{db}
}

func (md *ModelDao) GetTenantModel(tenantId, providerName, modelName, modelType string) (*model.ProviderModel, error) {
	var model *model.ProviderModel

	if err := md.db.Scopes(mysql.IDDesc()).Where("tenant_id = ? and provider_name = ? and model_name = ? and model_type = ?", tenantId, providerName, modelName, modelType).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, errors.WithCode(code.ErrDatabase, err.Error())
		}
	}
	return model, nil
}

func (md *ModelDao) UpdateModel(model *model.ProviderModel) error {
	if err := md.db.Updates(model).Error; err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil
}

func (md *ModelDao) CreateModel(model *model.ProviderModel) error {
	if err := md.db.Create(model).Error; err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil
}
