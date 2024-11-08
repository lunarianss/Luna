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
| ErrSuccess | 100001 | 200 | OK |
| ErrUnknown | 100002 | 500 | Internal server error |
| ErrBind | 100003 | 400 | Error occurred while binding the request body to the struct |
| ErrValidation | 100004 | 400 | Validation failed |
| ErrTokenInvalid | 100005 | 401 | Token invalid |
| ErrPageNotFound | 100006 | 404 | Page not found |
| ErrRestfulId | 100007 | 400 | Error occurred while parse restful id from url |
| ErrRunTimeCaller | 100008 | 500 | Error occurred while call go inner function |
| ErrRunTimeConfig | 100009 | 500 | Error occurred while runtime config is nil |
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
| ErrProviderMapModel | 110001 | 500 | Error occurred while attempt to index from providerMpa using provider |
| ErrProviderNotHaveIcon | 110002 | 500 | Error occurred while provider entity doesn't have icon property |
| ErrToOriginModelType | 110003 | 500 | Error occurred while convert to origin model type |
| ErrDefaultModelNotFound | 110004 | 500 | Error occurred while trying to convert default model to unknown |
| ErrModelSchemaNotFound | 110005 | 500 | Error occurred while attempt to index from predefined models using model name |
| ErrAllModelsEmpty | 110006 | 500 | Error occurred when all models are empty |

