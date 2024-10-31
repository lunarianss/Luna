package entities

type ConfigurationMethod string

const (
	PREDEFINED_MODEL   ConfigurationMethod = "predefined-model"
	CUSTOMIZABLE_MODEL ConfigurationMethod = "customizable-model"
)

type FormType string

const (
	TEXT_INPUT   FormType = "text-input"
	SECRET_INPUT FormType = "secret-input"
	SELECT       FormType = "select"
	RADIO        FormType = "radio"
	SWITCH       FormType = "switch"
)

type I18nObject struct {
	Zh_Hans string `json:"zh_hans"`
	En_US   string `json:"en_US"`
}

type ProviderHelpEntity struct {
	Title I18nObject `json:"title"`
	Url   I18nObject `json:"url"`
}

type FormShowOnObject struct {
	Variable string `json:"variable"`
	Value    string `json:"value"`
}

type CredentialFormSchema struct {
	Variable     string             `json:"variable"`
	Label        I18nObject         `json:"label"`
	Type         FormType           `json:"type"`
	Required     bool               `json:"required"`
	DefaultValue string             `json:"default"`
	MaxLength    int                `json:"max_length"`
	ShowOn       []FormShowOnObject `json:"show_on"`
}

type FieldModelSchema struct {
	Label       I18nObject `json:"label"`
	PlaceHolder I18nObject `json:"place_holder"`
}

type ModelCredentialSchema struct {
	Model                 FieldModelSchema       `json:"model"`
	CredentialFormSchemas []CredentialFormSchema `json:"credential_form"`
}

type ProviderCredentialSchema struct {
	CredentialFormSchemas []CredentialFormSchema
}

type ProviderEntity struct {
	Provider    string     `json:"provider"`
	Label       I18nObject `json:"label"`
	Description I18nObject `json:"description"`
	Icon_small  I18nObject `json:"icon_small"`
	Icon_large  I18nObject `json:"icon_large"`
	Background  string     `json:"background"`
	// Help string For Front
	Help ProviderHelpEntity `json:"help"`
	// Checkbox model type For Front
	SupportedModelTypes []ModelType `json:"supported_model_types"`
	// Settings(predefined-model) or Add Model(customizable-model) for Front
	ConfigurationMethods []ConfigurationMethod `json:"configuration_methods"`
	Models               []ProviderModel       `json:"models"`
	// Add Provider by add api key directly, Form component when click setting(like input: api key)
	ProviderCredentialSchema ProviderCredentialSchema `json:"provider_credential_schema"`
	// Add Provider by add single model, form component when click add model(like input: model name, api key)
	ModelCredentialSchema ModelCredentialSchema `json:"model_credential_schema"`
}
