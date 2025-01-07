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
	// ErrNotFoundToolRegistry - 500: tool registry is not found in the tool registry list.
	ErrNotFoundToolRegistry
	// ErrRequiredPromptMessage - 500: Prompt type is required when convert to prompt template.
	ErrRequiredPromptType
	// ErrQuotaExceed - 500: Your quota for Luna Hosted Model Provider has been exhausted,Please go to Settings -> Model Provider to complete your own provider credentials.
	ErrQuotaExceed
	// ErrAudioType - 500: Audio type error: only support extensions like mp3, mp4, mpeg, mpga, m4a, wav, webm, amr.
	ErrAudioType
	// ErrAudioFileToLarge - 500: Audio file is to large.
	ErrAudioFileToLarge
	// ErrAudioFileEmpty - 500: Audio file is empty.
	ErrAudioFileEmpty
	// ErrAudioTextEmpty - 500: Audio text is empty.
	ErrAudioTextEmpty
	// ErrTTSWebSocket - 500: Failed to connect tts websocket.
	ErrTTSWebSocket
	// ErrTTSWebSocketParse - 500: Failed to parse tts websocket message.
	ErrTTSWebSocketParse
	// ErrContextTimeout - 500: Context timeout.
	ErrContextTimeout
	// ErrTTSWebSocketWrite - 500: Failed to write tts websocket message.
	ErrTTSWebSocketWrite
	// ErrTencentARS - 500: Tencent ARS service error.
	ErrTencentARS
	// ErrNotStreamAgent - 500: Agent Chat App does not support blocking mode.
	ErrNotStreamAgent
	// ErrInvokeTool - 500: Failed to invoke agent tool.
	ErrInvokeTool
	// ErrInvokeToolUnConvertAble - 500: Failed to convert to tool message.
	ErrInvokeToolUnConvertAble
)
