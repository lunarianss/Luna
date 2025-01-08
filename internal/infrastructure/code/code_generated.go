// Copyright 2024 Benjamin <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Code generated by "codegen -type=int /Users/max/Documents/coding/Backend/Golang/Personal/Luna/internal/infrastructure/code"; DO NOT EDIT.

package code

import "github.com/lunarianss/Luna/infrastructure/errors" // init register error codes defines in this source code to `github.com/lunarianss/Luna/infrastructure/errors
func init() {
	errors.Enroll(ErrAppMapMode, 500, "Error occurred while attempt to index from appTemplate using mode")
	errors.Enroll(ErrAppNotFoundRelatedConfig, 500, "App config is not found")
	errors.Enroll(ErrAppStatusNotNormal, 500, "App is not active")
	errors.Enroll(ErrAppCodeNotFound, 500, "App code not found")
	errors.Enroll(ErrAppSiteDisabled, 400, "Site is disabled")
	errors.Enroll(ErrAppApiDisabled, 400, "Api is disabled")
	errors.Enroll(ErrAppTokenExceed, 400, "Count of app token is exceeded")
	errors.Enroll(ErrNotFoundJobID, 400, "Not found job ID")
	errors.Enroll(ErrVDBQueryError, 400, "Occurred error when vector similarity search")
	errors.Enroll(ErrVDBConstructError, 400, "Occurred error when construct vdb response")
	errors.Enroll(ErrEmailCode, 500, "Error occurred when email code is incorrect")
	errors.Enroll(ErrTokenEmail, 500, "Error occurred when email is incorrect")
	errors.Enroll(ErrTenantAlreadyExist, 500, "Error occurred when tenant is already exist")
	errors.Enroll(ErrAccountBanned, 500, "Error occurred when user is banned but still to operate")
	errors.Enroll(ErrTenantStatusArchive, 400, "Error occurred when tenant's status is archive")
	errors.Enroll(ErrSuccess, 200, "OK")
	errors.Enroll(ErrUnknown, 500, "Internal server error")
	errors.Enroll(ErrBind, 400, "Error occurred while request body is not incorrect")
	errors.Enroll(ErrValidation, 400, "Validation failed")
	errors.Enroll(ErrForbidden, 403, "You don't have the permission")
	errors.Enroll(ErrPageNotFound, 404, "Page not found")
	errors.Enroll(ErrResourceNotFound, 404, "The requested URL was not found on the server. If you entered the URL manually please check your spelling and try again")
	errors.Enroll(ErrRestfulId, 400, "Error occurred while parse restful id from url")
	errors.Enroll(ErrRunTimeCaller, 500, "Error occurred while call a system call")
	errors.Enroll(ErrRunTimeConfig, 500, "Error occurred while runtime config is nil")
	errors.Enroll(ErrMQSend, 500, "Error occurred while send sync message")
	errors.Enroll(ErrConcurrentLock, 500, "Please don't click repeatedly")
	errors.Enroll(ErrDatabase, 500, "Database error")
	errors.Enroll(ErrRecordNotFound, 500, "Database record not found")
	errors.Enroll(ErrScanToField, 400, "Database scan error to field")
	errors.Enroll(ErrVDB, 400, "Vector Database error")
	errors.Enroll(ErrRedis, 400, "Redis error")
	errors.Enroll(ErrEncrypt, 401, "Error occurred while encrypting the user password")
	errors.Enroll(ErrSignatureInvalid, 401, "Signature is invalid")
	errors.Enroll(ErrExpired, 401, "Token expired")
	errors.Enroll(ErrInvalidAuthHeader, 401, "Invalid authorization header")
	errors.Enroll(ErrMissingHeader, 401, "The `Authorization` header was empty")
	errors.Enroll(ErrPasswordIncorrect, 401, "Password was incorrect")
	errors.Enroll(ErrPermissionDenied, 403, "Permission denied")
	errors.Enroll(ErrEncodingFailed, 500, "Encoding failed due to an error with the data")
	errors.Enroll(ErrDecodingFailed, 500, "Decoding failed due to an error with the data")
	errors.Enroll(ErrInvalidJSON, 500, "Data is not valid JSON")
	errors.Enroll(ErrEncodingJSON, 500, "JSON data could not be encoded")
	errors.Enroll(ErrDecodingJSON, 500, "JSON data could not be decoded")
	errors.Enroll(ErrInvalidYaml, 500, "Data is not valid Yaml")
	errors.Enroll(ErrEncodingYaml, 500, "Yaml data could not be encoded")
	errors.Enroll(ErrDecodingYaml, 500, "Yaml data could not be decoded")
	errors.Enroll(ErrEncodingBase64, 500, "Base64 data could not be encoded")
	errors.Enroll(ErrDecodingBase64, 500, "Base64 data could not be decoded")
	errors.Enroll(ErrTokenGenerate, 500, "Error occurred when Token generate")
	errors.Enroll(ErrTokenExpired, 500, "Error occurred when Token expired")
	errors.Enroll(ErrTokenInvalid, 401, "Token invalid")
	errors.Enroll(ErrTokenMethodErr, 500, "Unexpected signing method")
	errors.Enroll(ErrTokenInsNotFound, 500, "Jwt instance is not found")
	errors.Enroll(ErrRefreshTokenNotFound, 500, "Refresh token is not found in redis")
	errors.Enroll(ErrTokenMissBearer, 401, "The token does not conform to the format")
	errors.Enroll(ErrRedisSetKey, 500, "Error occurred when set key, value to redis")
	errors.Enroll(ErrRedisSetExpire, 500, "Error occurred when set expire  to redis")
	errors.Enroll(ErrRedisRuntime, 500, "Error occurred when invoke redis api")
	errors.Enroll(ErrRedisDataExpire, 500, "Error occurred when data is expired")
	errors.Enroll(ErrRSAGenerate, 500, "Error occurred when generate pair of rsa key")
	errors.Enroll(ErrGinNotExistAccountInfo, 400, "Error occurred when get account info from gin context")
	errors.Enroll(ErrGinNotExistAppSiteInfo, 400, "Error occurred when get app site info from gin context")
	errors.Enroll(ErrGinNotExistServiceTokenInfo, 400, "Error occurred when get app service token info from gin context")
	errors.Enroll(ErrOnlyOverrideConfigInDebugger, 500, "Error occurred while attempt to override config in non-debug mode")
	errors.Enroll(ErrModelEmptyInConfig, 500, "Error occurred while attempt to index model from config")
	errors.Enroll(ErrRequiredCorrectProvider, 500, "Error occurred when provider is not found or provider isn't include in the provider list")
	errors.Enroll(ErrRequiredModelName, 500, "Error occurred when model name is not found in model config")
	errors.Enroll(ErrRequiredCorrectModel, 500, "Error occurred when model is not found or model isn't include in the model list")
	errors.Enroll(ErrRequiredOverrideConfig, 500, "Config_from is ARGS that override_config_dict is required")
	errors.Enroll(ErrNotFoundModelRegistry, 500, "Model registry is not found in the model registry list")
	errors.Enroll(ErrNotFoundToolRegistry, 500, "Tool registry is not found in the tool registry list")
	errors.Enroll(ErrRequiredPromptType, 500, "Prompt type is required when convert to prompt template")
	errors.Enroll(ErrQuotaExceed, 500, "Your quota for Luna Hosted Model Provider has been exhausted,Please go to Settings -> Model Provider to complete your own provider credentials")
	errors.Enroll(ErrAudioType, 500, "Audio type error: only support extensions like mp3, mp4, mpeg, mpga, m4a, wav, webm, amr")
	errors.Enroll(ErrAudioFileToLarge, 500, "Audio file is to large")
	errors.Enroll(ErrAudioFileEmpty, 500, "Audio file is empty")
	errors.Enroll(ErrAudioTextEmpty, 500, "Audio text is empty")
	errors.Enroll(ErrTTSWebSocket, 500, "Failed to connect tts websocket")
	errors.Enroll(ErrTTSWebSocketParse, 500, "Failed to parse tts websocket message")
	errors.Enroll(ErrContextTimeout, 500, "Context timeout")
	errors.Enroll(ErrTTSWebSocketWrite, 500, "Failed to write tts websocket message")
	errors.Enroll(ErrTencentARS, 500, "Tencent ARS service error")
	errors.Enroll(ErrNotStreamAgent, 500, "Agent Chat App does not support blocking mode")
	errors.Enroll(ErrInvokeTool, 500, "Failed to invoke agent tool")
	errors.Enroll(ErrToolParameter, 500, "Failed to parse tool parameter")
	errors.Enroll(ErrInvokeToolUnConvertAble, 500, "Failed to convert to tool message")
	errors.Enroll(ErrProviderMapModel, 500, "Error occurred while attempt to index from providerMpa using provider")
	errors.Enroll(ErrProviderNotHaveIcon, 500, "Error occurred while provider entity doesn't have icon property")
	errors.Enroll(ErrToOriginModelType, 500, "Error occurred while convert to origin model type")
	errors.Enroll(ErrDefaultModelNotFound, 500, "Error occurred while default model is not exist")
	errors.Enroll(ErrModelSchemaNotFound, 500, "Error occurred while attempt to index from predefined models using model name")
	errors.Enroll(ErrAllModelsEmpty, 500, "Error occurred when all models are empty")
	errors.Enroll(ErrModelNotHavePosition, 500, "Error occurred when models haven't position definition")
	errors.Enroll(ErrModelNotHavePrice, 500, "Error occurred when models haven't price definition")
	errors.Enroll(ErrModelNotHaveEndPoint, 500, "Error occurred when models haven't url endpoint")
	errors.Enroll(ErrModelUrlNotConvertUrl, 500, "Error occurred when models url interface{} convert ot string ")
	errors.Enroll(ErrTypeOfPromptMessage, 500, "When prompt type is user, the type of message is neither string or []*promptMessageContent")
	errors.Enroll(ErrCallLargeLanguageModel, 500, "Error occurred when call large language model post api")
	errors.Enroll(ErrConvertDelimiterString, 500, "Error occurred when convert delimiter to string")
	errors.Enroll(ErrNotSetManagerForProvider, 500, "Error occurred when not set manager for provider")
	errors.Enroll(ErrTTSModelNotVoice, 500, "Error occurred when tts model doesn't have voice")
}
