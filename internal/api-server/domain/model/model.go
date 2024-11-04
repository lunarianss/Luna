package model

import (
	"encoding/json"

	entities "github.com/lunarianss/Luna/internal/api-server/entities/provider"
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

func (mpd *ModelDomain) AddOrUpdateCustomModelCredentials(providerConfiguration *entities.ProviderConfiguration, credentialParam map[string]interface{}, modelType, modelName string) error {

	modelRecord, credentials, err := mpd.validateProviderCredentials(providerConfiguration, credentialParam, modelType, modelName)

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

		if err := mpd.ModelRepo.UpdateModel(modelRecord); err != nil {
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

		if err := mpd.ModelRepo.CreateModel(model); err != nil {
			return err
		}
	}
	return nil
}

func (mpd *ModelDomain) validateProviderCredentials(providerConfiguration *entities.ProviderConfiguration, credentials map[string]interface{}, modelType, modeName string) (*model.ProviderModel, map[string]interface{}, error) {

	model, err := mpd.ModelRepo.GetTenantModel(providerConfiguration.TenantId, providerConfiguration.Provider.Provider, modeName, modelType)

	if err != nil {
		return nil, nil, err
	}

	// credentials 对 apikey 进行 validate and encrypt
	return model, credentials, nil
}
