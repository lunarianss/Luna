package code

const (
	// ErrProviderMapModel - 500: Error occurred while attempt to index from providerMpa using provider.
	ErrProviderMapModel int = iota + 110001
	// ErrProviderNotHaveIcon - 500: Error occurred while provider entity doesn't have icon property.
	ErrProviderNotHaveIcon
	// ErrToOriginModelType - 500: Error occurred while convert to origin model type.
	ErrToOriginModelType
	// ErrDefaultModelNotFound - 500: Error occurred while trying to convert default model to unknown.
	ErrDefaultModelNotFound
	// ErrModelSchemaNotFound - 500: Error occurred while attempt to index from predefined models using model name.
	ErrModelSchemaNotFound
	// ErrAllModelsEmpty - 500: Error occurred when all models are empty.
	ErrAllModelsEmpty
)
