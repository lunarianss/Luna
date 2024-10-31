package base

type IModelProviderRepo interface {
	ValidateProviderCredentials() error
	GetProviderSchema()
}
