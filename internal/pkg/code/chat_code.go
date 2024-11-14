package code

const (
	// ErrAppMapMode - 500: Error occurred while attempt to override config in non-debug mode.
	ErrOnlyOverrideConfigInDebugger int = iota + 110201
	// ErrModelEmptyInConfig - 500: Error occurred while attempt to index model from config.
	ErrModelEmptyInConfig
	// ErrRequiredCorrectProvider - 500: Error occurred when provider is not found or provider isn't include in the provider list.
	ErrRequiredCorrectProvider
	// ErrRequiredCorrectProvider - 500: Error occurred when model name is not found in model config.
	ErrRequiredModelName
	// ErrRequiredCorrectModel - 500: Error occurred when model is not found or model isn't include in the model list.
	ErrRequiredCorrectModel
	// ErrRequiredOverrideConfig - 500: config_from is ARGS that override_config_dict is required
	ErrRequiredOverrideConfig
	// ErrNotFoundModelRegistry - 500: model registry is not found in the model registry  list
	ErrNotFoundModelRegistry
)
