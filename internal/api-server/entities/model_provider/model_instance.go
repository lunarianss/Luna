package model_provider

import (
	"github.com/lunarianss/Luna/internal/api-server/entities/model_runtime"
)

type ModelStatus string

const (
	ACTIVE         ModelStatus = "active"
	NO_CONFIGURE   ModelStatus = "no-configure"
	QUOTA_EXCEEDED ModelStatus = "quota-exceeded"
	NO_PERMISSION  ModelStatus = "no-permission"
	DISABLED       ModelStatus = "disabled"
)

type ProviderModelBundle struct {
	Configuration     *ProviderConfiguration
	ProviderInstance  ModelProvider
	ModelTypeInstance *model_runtime.AIModel
}

type ModelInstance struct {
	ProviderModelBundle *ProviderModelBundle   `json:"provider_model_bundle"`
	Model               string                 `json:"model"`
	Provider            string                 `json:"provider"`
	Credentials         map[string]any         `json:"credentials"`
	ModelTypeInstance   *model_runtime.AIModel `json:"model_type_instance"`
}
