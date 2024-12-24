# 错误码

！！IAM 系统错误码列表，由 `codegen -type=int -doc` 命令生成，不要对此文件做任何更改。

## 功能说明

如果返回结果中存在 `code` 字段，则表示调用 API 接口失败。例如：

```json
{
  "code": 100101,
  "message": "Database error"
}
```

上述返回中 `code` 表示错误码，`message` 表示该错误的具体信息。每个错误同时也对应一个 HTTP 状态码，比如上述错误码对应了 HTTP 状态码 500(Internal Server Error)。

## 错误码列表

IAM 系统支持的错误码列表如下：

| Identifier | Code | HTTP Code | Description |
| ---------- | ---- | --------- | ----------- |
| ErrAppMapMode | 110101 | 500 | Error occurred while attempt to index from appTemplate using mode |
| ErrAppNotFoundRelatedConfig | 110102 | 500 | App config is not found |
| ErrAppStatusNotNormal | 110103 | 500 | App is not active |
| ErrAppCodeNotFound | 110104 | 500 | App code not found |
| ErrAppSiteDisabled | 110105 | 400 | Site is disabled |
| ErrAppApiDisabled | 110106 | 400 | Api is disabled |
| ErrAppTokenExceed | 110107 | 400 | Count of app token is exceeded |
| ErrEmailCode | 110301 | 500 | Error occurred when email code is incorrect |
| ErrTokenEmail | 110302 | 500 | Error occurred when email is incorrect |
| ErrTenantAlreadyExist | 110303 | 500 | Error occurred when tenant is already exist |
| ErrAccountBanned | 110304 | 500 | Error occurred when user is banned but still to operate |
| ErrTenantStatusArchive | 110305 | 400 | Error occurred when tenant's status is archive |
| ErrSuccess | 100001 | 200 | OK |
| ErrUnknown | 100002 | 500 | Internal server error |
| ErrBind | 100003 | 400 | Error occurred while request body is not incorrect |
| ErrValidation | 100004 | 400 | Validation failed |
| ErrForbidden | 100005 | 403 | You don't have the permission |
| ErrPageNotFound | 100006 | 404 | Page not found |
| ErrResourceNotFound | 100007 | 404 | The requested URL was not found on the server. If you entered the URL manually please check your spelling and try again |
| ErrRestfulId | 100008 | 400 | Error occurred while parse restful id from url |
| ErrRunTimeCaller | 100009 | 500 | Error occurred while call a system call |
| ErrRunTimeConfig | 100010 | 500 | Error occurred while runtime config is nil |
| ErrMQSend | 100011 | 500 | Error occurred while send message to mq |
| ErrConcurrentLock | 100012 | 500 | Please don't click repeatedly |
| ErrDatabase | 100101 | 500 | Database error |
| ErrRecordNotFound | 100102 | 500 | Database record not found |
| ErrScanToField | 100103 | 400 | Database scan error to field |
| ErrEncrypt | 100201 | 401 | Error occurred while encrypting the user password |
| ErrSignatureInvalid | 100202 | 401 | Signature is invalid |
| ErrExpired | 100203 | 401 | Token expired |
| ErrInvalidAuthHeader | 100204 | 401 | Invalid authorization header |
| ErrMissingHeader | 100205 | 401 | The `Authorization` header was empty |
| ErrPasswordIncorrect | 100206 | 401 | Password was incorrect |
| ErrPermissionDenied | 100207 | 403 | Permission denied |
| ErrEncodingFailed | 100301 | 500 | Encoding failed due to an error with the data |
| ErrDecodingFailed | 100302 | 500 | Decoding failed due to an error with the data |
| ErrInvalidJSON | 100303 | 500 | Data is not valid JSON |
| ErrEncodingJSON | 100304 | 500 | JSON data could not be encoded |
| ErrDecodingJSON | 100305 | 500 | JSON data could not be decoded |
| ErrInvalidYaml | 100306 | 500 | Data is not valid Yaml |
| ErrEncodingYaml | 100307 | 500 | Yaml data could not be encoded |
| ErrDecodingYaml | 100308 | 500 | Yaml data could not be decoded |
| ErrTokenGenerate | 100401 | 500 | Error occurred when Token generate |
| ErrTokenExpired | 100402 | 500 | Error occurred when Token expired |
| ErrTokenInvalid | 100403 | 401 | Token invalid |
| ErrTokenMethodErr | 100404 | 500 | Unexpected signing method |
| ErrTokenInsNotFound | 100405 | 500 | Jwt instance is not found |
| ErrRefreshTokenNotFound | 100406 | 500 | Refresh token is not found in redis |
| ErrTokenMissBearer | 100407 | 401 | The token does not conform to the format |
| ErrRedisSetKey | 100501 | 500 | Error occurred when set key, value to redis |
| ErrRedisSetExpire | 100502 | 500 | Error occurred when set expire  to redis |
| ErrRedisRuntime | 100503 | 500 | Error occurred when invoke redis api |
| ErrRedisDataExpire | 100504 | 500 | Error occurred when data is expired |
| ErrRSAGenerate | 100601 | 500 | Error occurred when generate pair of rsa key |
| ErrGinNotExistAccountInfo | 100701 | 400 | Error occurred when get account info from gin context |
| ErrGinNotExistAppSiteInfo | 100702 | 400 | Error occurred when get app site info from gin context |
| ErrGinNotExistServiceTokenInfo | 100703 | 400 | Error occurred when get app service token info from gin context |
| ErrOnlyOverrideConfigInDebugger | 110201 | 500 | Error occurred while attempt to override config in non-debug mode |
| ErrModelEmptyInConfig | 110202 | 500 | Error occurred while attempt to index model from config |
| ErrRequiredCorrectProvider | 110203 | 500 | Error occurred when provider is not found or provider isn't include in the provider list |
| ErrRequiredModelName | 110204 | 500 | Error occurred when model name is not found in model config |
| ErrRequiredCorrectModel | 110205 | 500 | Error occurred when model is not found or model isn't include in the model list |
| ErrRequiredOverrideConfig | 110206 | 500 | Config_from is ARGS that override_config_dict is required |
| ErrNotFoundModelRegistry | 110207 | 500 | Model registry is not found in the model registry list |
| ErrRequiredPromptType | 110208 | 500 | Prompt type is required when convert to prompt template |
| ErrQuotaExceed | 110209 | 500 | Your quota for Luna Hosted Model Provider has been exhausted,Please go to Settings -> Model Provider to complete your own provider credentials |
| ErrAudioType | 110210 | 500 | Audio type error: only support extensions like mp3, mp4, mpeg, mpga, m4a, wav, webm, amr |
| ErrAudioFileToLarge | 110211 | 500 | Audio file is to large |
| ErrAudioFileEmpty | 110212 | 500 | Audio file is empty |
| ErrAudioTextEmpty | 110213 | 500 | Audio text is empty |
| ErrTTSWebSocket | 110214 | 500 | Failed to connect tts websocket |
| ErrTTSWebSocketParse | 110215 | 500 | Failed to parse tts websocket message |
| ErrContextTimeout | 110216 | 500 | Context timeout |
| ErrTTSWebSocketWrite | 110217 | 500 | Failed to write tts websocket message |
| ErrTencentARS | 110218 | 500 | Tencent ARS service error |
| ErrProviderMapModel | 110001 | 500 | Error occurred while attempt to index from providerMpa using provider |
| ErrProviderNotHaveIcon | 110002 | 500 | Error occurred while provider entity doesn't have icon property |
| ErrToOriginModelType | 110003 | 500 | Error occurred while convert to origin model type |
| ErrDefaultModelNotFound | 110004 | 500 | Error occurred while default model is not exist |
| ErrModelSchemaNotFound | 110005 | 500 | Error occurred while attempt to index from predefined models using model name |
| ErrAllModelsEmpty | 110006 | 500 | Error occurred when all models are empty |
| ErrModelNotHavePosition | 110007 | 500 | Error occurred when models haven't position definition |
| ErrModelNotHavePrice | 110008 | 500 | Error occurred when models haven't price definition |
| ErrModelNotHaveEndPoint | 110009 | 500 | Error occurred when models haven't url endpoint |
| ErrModelUrlNotConvertUrl | 110010 | 500 | Error occurred when models url interface{} convert ot string  |
| ErrTypeOfPromptMessage | 110011 | 500 | When prompt type is user, the type of message is neither string or []*promptMessageContent |
| ErrCallLargeLanguageModel | 110012 | 500 | Error occurred when call large language model post api |
| ErrConvertDelimiterString | 110013 | 500 | Error occurred when convert delimiter to string |
| ErrNotSetManagerForProvider | 110014 | 500 | Error occurred when not set manager for provider |
| ErrTTSModelNotVoice | 110015 | 500 | Error occurred when tts model doesn't have voice |

