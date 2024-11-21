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
| ErrAppNotFoundRelatedConfig | 110102 | 500 | Error occurred while attempt to find app related config |
| ErrEmailCode | 110301 | 500 | Error occurred when email code is incorrect |
| ErrTokenEmail | 110302 | 500 | Error occurred when email is incorrect |
| ErrTenantAlreadyExist | 110303 | 500 | Error occurred when tenant is already exist |
| ErrAccountBanned | 110304 | 500 | Error occurred when user is banned but still to operate |
| ErrSuccess | 100001 | 200 | OK |
| ErrUnknown | 100002 | 500 | Internal server error |
| ErrBind | 100003 | 400 | Error occurred while binding the request body to the struct |
| ErrValidation | 100004 | 400 | Validation failed |
| ErrPageNotFound | 100005 | 404 | Page not found |
| ErrRestfulId | 100006 | 400 | Error occurred while parse restful id from url |
| ErrRunTimeCaller | 100007 | 500 | Error occurred while call go inner function |
| ErrRunTimeConfig | 100008 | 500 | Error occurred while runtime config is nil |
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
| ErrTokenMissBearer | 100407 | 401 | Token miss a header of Bearer  |
| ErrRedisSetKey | 100501 | 500 | Error occurred when set key, value to redis |
| ErrRedisSetExpire | 100502 | 500 | Error occurred when set expire  to redis |
| ErrRedisRuntime | 100503 | 500 | Error occurred when invoke redis api |
| ErrRedisDataExpire | 100504 | 500 | Error occurred when redis data is expired |
| ErrRSAGenerate | 100601 | 500 | Error occurred when generate pair of rsa key |
| ErrGinNotExistAccountInfo | 100701 | 400 | Error occurred when get account info from gin context |
| ErrOnlyOverrideConfigInDebugger | 110201 | 500 | Error occurred while attempt to override config in non-debug mode |
| ErrModelEmptyInConfig | 110202 | 500 | Error occurred while attempt to index model from config |
| ErrRequiredCorrectProvider | 110203 | 500 | Error occurred when provider is not found or provider isn't include in the provider list |
| ErrRequiredModelName | 110204 | 500 | Error occurred when model name is not found in model config |
| ErrRequiredCorrectModel | 110205 | 500 | Error occurred when model is not found or model isn't include in the model list |
| ErrRequiredOverrideConfig | 110206 | 500 | Config_from is ARGS that override_config_dict is required |
| ErrNotFoundModelRegistry | 110207 | 500 | Model registry is not found in the model registry list |
| ErrRequiredPromptType | 110208 | 500 | Prompt type is required when convert to prompt template |
| ErrProviderMapModel | 110001 | 500 | Error occurred while attempt to index from providerMpa using provider |
| ErrProviderNotHaveIcon | 110002 | 500 | Error occurred while provider entity doesn't have icon property |
| ErrToOriginModelType | 110003 | 500 | Error occurred while convert to origin model type |
| ErrDefaultModelNotFound | 110004 | 500 | Error occurred while trying to convert default model to unknown |
| ErrModelSchemaNotFound | 110005 | 500 | Error occurred while attempt to index from predefined models using model name |
| ErrAllModelsEmpty | 110006 | 500 | Error occurred when all models are empty |
| ErrModelNotHavePosition | 110007 | 500 | Error occurred when models haven't position definition |
| ErrModelNotHaveEndPoint | 110008 | 500 | Error occurred when models haven't url endpoint |
| ErrModelUrlNotConvertUrl | 110009 | 500 | Error occurred when models url interface{} convert ot string  |
| ErrTypeOfPromptMessage | 110010 | 500 | When prompt type is user, the type of message is neither string or []*promptMessageContent |
| ErrCallLargeLanguageModel | 110011 | 500 | Error occurred when call large language model post api |
| ErrConvertDelimiterString | 110012 | 500 | Error occurred when convert delimiter to string |

