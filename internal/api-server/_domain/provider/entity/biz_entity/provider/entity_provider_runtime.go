package biz_entity

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	common "github.com/lunarianss/Luna/internal/api-server/_domain/provider/entity/biz_entity/common_relation"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/_domain/provider/entity/biz_entity/provider/model_provider"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
	"gopkg.in/yaml.v3"
)

type ProviderRuntime struct {
	ProviderSchema   *ProviderStaticConfiguration
	ModelConfPath    string
	ModelInstanceMap map[string]*biz_entity.AIModelRuntime
}

func (mp *ProviderRuntime) Models(modelType common.ModelType) ([]*biz_entity.AIModelStaticConfiguration, error) {

	providerEntity, err := mp.GetProviderSchema()

	if !slices.Contains(providerEntity.SupportedModelTypes, modelType) {
		return nil, nil
	}

	AIModelInstance := mp.GetModelInstance(modelType)

	if err != nil {
		return nil, err
	}

	return AIModelInstance.PredefinedModels()
}

func (mp *ProviderRuntime) GetProviderSchema() (*ProviderStaticConfiguration, error) {
	providerName := filepath.Base(mp.ModelConfPath)
	providerSchemaPath := fmt.Sprintf("%s/%s.yaml", mp.ModelConfPath, providerName)
	providerContent, err := os.ReadFile(providerSchemaPath)

	if err != nil {
		return nil, errors.WithCode(code.ErrRunTimeCaller, err.Error())
	}

	provider := &ProviderStaticConfiguration{}
	err = yaml.Unmarshal(providerContent, provider)

	if err != nil {
		return nil, errors.WithCode(code.ErrRunTimeCaller, err.Error())
	}

	if provider.ModelCredentialSchema != nil {
		for _, c := range provider.ModelCredentialSchema.CredentialFormSchemas {
			if c.ShowOn == nil {
				c.ShowOn = []*FormShowOnObject{}
			}
		}
	}

	if provider.ProviderCredentialSchema != nil {
		for _, c := range provider.ProviderCredentialSchema.CredentialFormSchemas {
			if c.ShowOn == nil {
				c.ShowOn = []*FormShowOnObject{}
			}
		}
	}

	provider.PatchI18nObject()
	return provider, nil
}

func (mp *ProviderRuntime) GetModelInstance(modelType common.ModelType) *biz_entity.AIModelRuntime {
	providerName := filepath.Base(mp.ModelConfPath)
	modelSchemaPath := fmt.Sprintf("%s/%s", mp.ModelConfPath, modelType)
	mp.ModelInstanceMap = make(map[string]*biz_entity.AIModelRuntime)

	if _, ok := mp.ModelInstanceMap[fmt.Sprintf("%s.%s", providerName, modelType)]; ok {
		return mp.ModelInstanceMap[fmt.Sprintf("%s.%s", providerName, modelType)]
	}

	AIModel := &biz_entity.AIModelRuntime{
		ModelType:     modelType,
		ModelConfPath: modelSchemaPath,
	}

	mp.ModelInstanceMap[fmt.Sprintf("%s.%s", providerName, modelType)] = AIModel

	return AIModel
}
