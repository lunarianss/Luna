package biz_entity

type ModelWithProvider struct {
	*ProviderModelWithStatus
	Provider *SimpleModelProvider `json:"provider"`
}
