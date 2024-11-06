package model

import (
	"context"
	"encoding/json"

	"github.com/lunarianss/Luna/internal/api-server/entities/model_provider"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/repo"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
)

type ModelDomain struct {
	ModelRepo repo.ModelRepo
}

func NewModelDomain(modelRepo repo.ModelRepo) *ModelDomain {
	return &ModelDomain{
		ModelRepo: modelRepo,
	}
}

func (mpd *ModelDomain) GetDefaultModelInstance(ctx context.Context, tenantID, modelType string) (*model_provider.ModelInstance, error) {

	return nil, nil

}

func (mpd *ModelDomain) AddOrUpdateCustomModelCredentials(ctx context.Context, providerConfiguration *model_provider.ProviderConfiguration, credentialParam map[string]interface{}, modelType, modelName string) error {

	modelRecord, credentials, err := mpd.validateProviderCredentials(ctx, providerConfiguration, credentialParam, modelType, modelName)

	if err != nil {
		return err
	}

	byteCredentials, err := json.Marshal(credentials)

	if err != nil {
		return errors.WithCode(code.ErrEncodingJSON, err.Error())
	}

	if modelRecord != nil {
		modelRecord.EncryptedConfig = string(byteCredentials)
		modelRecord.IsValid = 1

		if err := mpd.ModelRepo.UpdateModel(ctx, modelRecord); err != nil {
			return err
		}

	} else {
		model := &model.ProviderModel{
			ProviderName:    providerConfiguration.Provider.Provider,
			ModelName:       modelName,
			ModelType:       modelType,
			EncryptedConfig: string(byteCredentials),
			IsValid:         1,
			TenantID:        providerConfiguration.TenantId,
		}

		if err := mpd.ModelRepo.CreateModel(ctx, model); err != nil {
			return err
		}
	}
	return nil
}

func (mpd *ModelDomain) validateProviderCredentials(ctx context.Context, providerConfiguration *model_provider.ProviderConfiguration, credentials map[string]interface{}, modelType, modeName string) (*model.ProviderModel, map[string]interface{}, error) {

	model, err := mpd.ModelRepo.GetTenantModel(ctx, providerConfiguration.TenantId, providerConfiguration.Provider.Provider, modeName, modelType)

	if err != nil {
		return nil, nil, err
	}

	// credentials 对 apikey 进行 validate and encrypt
	return model, credentials, nil
}
