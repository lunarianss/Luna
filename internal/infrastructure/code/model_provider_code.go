// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package code

const (
	// ErrProviderMapModel - 500: Error occurred while attempt to index from providerMpa using provider.
	ErrProviderMapModel int = iota + 110001
	// ErrProviderNotHaveIcon - 500: Error occurred while provider entity doesn't have icon property.
	ErrProviderNotHaveIcon
	// ErrToOriginModelType - 500: Error occurred while convert to origin model type.
	ErrToOriginModelType
	// ErrDefaultModelNotFound - 500: Error occurred while default model is not exist.
	ErrDefaultModelNotFound
	// ErrModelSchemaNotFound - 500: Error occurred while attempt to index from predefined models using model name.
	ErrModelSchemaNotFound
	// ErrAllModelsEmpty - 500: Error occurred when all models are empty.
	ErrAllModelsEmpty
	// ErrAllModelsEmpty - 500: Error occurred when models haven't position definition.
	ErrModelNotHavePosition
	// ErrModelNotHavePrice - 500: Error occurred when models haven't price definition.
	ErrModelNotHavePrice
	// ErrModelNotHaveEndPoint - 500: Error occurred when models haven't url endpoint.
	ErrModelNotHaveEndPoint
	// ErrModelUrlNotConvertUrl - 500: Error occurred when models url interface{} convert ot string .
	ErrModelUrlNotConvertUrl
	// ErrTypeOfPromptMessage - 500: When prompt type is user, the type of message is neither string or []*promptMessageContent.
	ErrTypeOfPromptMessage
	// ErrCallLargeLanguageModel - 500: Error occurred when call large language model post api.
	ErrCallLargeLanguageModel
	// ErrConvertDelimiterString - 500: Error occurred when convert delimiter to string.
	ErrConvertDelimiterString
	// ErrNotSetManagerForProvider - 500: Error occurred when not set manager for provider.
	ErrNotSetManagerForProvider
	// ErrTTSModelNotVoice - 500: Error occurred when tts model doesn't have voice.
	ErrTTSModelNotVoice
)
