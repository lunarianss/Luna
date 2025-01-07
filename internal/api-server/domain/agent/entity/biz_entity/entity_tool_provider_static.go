package biz_entity

type ToolProviderStatic struct {
	Identity          *ToolProviderIdentity               `json:"identity" yaml:"identity"`

	CredentialsSchema map[string]*ToolProviderCredentials `json:"credentials_for_provider" yaml:"credentials_for_provider"`
}
