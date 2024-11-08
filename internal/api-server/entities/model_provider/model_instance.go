package model_provider

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
	ProviderInstance  *ModelProvider
	ModelTypeInstance *AIModel
}

type ModelInstance struct {
	ProviderModelBundle *ProviderModelBundle `json:"provider_model_bundle"`
	Model               string               `json:"model"`
	Provider            string               `json:"provider"`
	Credentials         interface{}          `json:"credentials"`
	ModelTypeInstance   *AIModel             `json:"model_type_instance"`
}
