package base

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/lunarianss/Hurricane/internal/api-server/model-runtime/entities"
	"github.com/lunarianss/Hurricane/internal/pkg/code"
	"github.com/lunarianss/Hurricane/pkg/errors"
)

type IModelProviderRepo interface {
	ValidateProviderCredentials() error
}

type ModelProvider struct {
	ProviderSchema entities.ProviderEntity
	ModelConfPath  string
}

func (mp *ModelProvider) GetProviderSchema() (*entities.ProviderEntity, error) {
	providerName := filepath.Base(mp.ModelConfPath)
	providerSchemaPath := fmt.Sprintf("%s/%s.yaml", mp.ModelConfPath, providerName)
	providerContent, err := os.ReadFile(providerSchemaPath)

	if err != nil {
		return nil, errors.WithCode(code.ErrRunTimeCaller, err.Error())
	}

	provider := &entities.ProviderEntity{}
	err = yaml.Unmarshal(providerContent, provider)

	if err != nil {
		return nil, errors.WithCode(code.ErrRunTimeCaller, err.Error())
	}

	return provider, nil
}
