// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

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
	// ErrRequiredOverrideConfig - 500: Config_from is ARGS that override_config_dict is required.
	ErrRequiredOverrideConfig
	// ErrNotFoundModelRegistry - 500: Model registry is not found in the model registry list.
	ErrNotFoundModelRegistry
	// ErrRequiredPromptMessage - 500: Prompt type is required when convert to prompt template.
	ErrRequiredPromptType
	// ErrQuotaExceed - 500: Your quota for Luna Hosted Model Provider has been exhausted,Please go to Settings -> Model Provider to complete your own provider credentials.
	ErrQuotaExceed
)
