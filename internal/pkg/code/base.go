// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package code

// Common: basic errors.
// Code must start with 1xxxxx.
const (
	// ErrSuccess - 200: OK.
	ErrSuccess int = iota + 100001

	// ErrUnknown - 500: Internal server error.
	ErrUnknown

	// ErrBind - 400: Error occurred while binding the request body to the struct.
	ErrBind

	// ErrValidation - 400: Validation failed.
	ErrValidation

	// ErrPageNotFound - 404: Page not found.
	ErrPageNotFound

	// ErrRestfulId - 400: Error occurred while parse restful id from url.
	ErrRestfulId

	// ErrRunTimeCaller - 500: Error occurred while call go inner function.
	ErrRunTimeCaller

	// ErrRunTimeConfig - 500: Error occurred while runtime config is nil.
	ErrRunTimeConfig
)

// common: database errors.
const (
	// ErrDatabase - 500: Database error.
	ErrDatabase int = iota + 100101
	// ErrRecordNotFound - 500: Database record not found.
	ErrRecordNotFound
	// ErrScanToField - 400: Database scan error to field.
	ErrScanToField
)

// common: authorization and authentication errors.
const (
	// ErrEncrypt - 401: Error occurred while encrypting the user password.
	ErrEncrypt int = iota + 100201

	// ErrSignatureInvalid - 401: Signature is invalid.
	ErrSignatureInvalid

	// ErrExpired - 401: Token expired.
	ErrExpired

	// ErrInvalidAuthHeader - 401: Invalid authorization header.
	ErrInvalidAuthHeader

	// ErrMissingHeader - 401: The `Authorization` header was empty.
	ErrMissingHeader

	// ErrPasswordIncorrect - 401: Password was incorrect.
	ErrPasswordIncorrect

	// PermissionDenied - 403: Permission denied.
	ErrPermissionDenied
)

// common: encode/decode errors.
const (
	// ErrEncodingFailed - 500: Encoding failed due to an error with the data.
	ErrEncodingFailed int = iota + 100301

	// ErrDecodingFailed - 500: Decoding failed due to an error with the data.
	ErrDecodingFailed

	// ErrInvalidJSON - 500: Data is not valid JSON.
	ErrInvalidJSON

	// ErrEncodingJSON - 500: JSON data could not be encoded.
	ErrEncodingJSON

	// ErrDecodingJSON - 500: JSON data could not be decoded.
	ErrDecodingJSON

	// ErrInvalidYaml - 500: Data is not valid Yaml.
	ErrInvalidYaml

	// ErrEncodingYaml - 500: Yaml data could not be encoded.
	ErrEncodingYaml

	// ErrDecodingYaml - 500: Yaml data could not be decoded.
	ErrDecodingYaml
)

const (
	// ErrTokenInvalid - 500: Error occurred when Token generate.
	ErrTokenGenerate int = iota + 100401
	// ErrTokenExpired - 500: Error occurred when Token expired.
	ErrTokenExpired
	// ErrTokenInvalid - 401: Token invalid.
	ErrTokenInvalid
	// ErrTokenMethodErr - 500: Unexpected signing method.
	ErrTokenMethodErr
	// ErrTokenInsNotFound - 500: Jwt instance is not found.
	ErrTokenInsNotFound
	// ErrRefreshTokenNotFound - 500: Refresh token is not found in redis.
	ErrRefreshTokenNotFound
)

const (
	// ErrRedisSetKey - 500: Error occurred when set key, value to redis.
	ErrRedisSetKey int = iota + 100501

	// ErrRedisSetExpire - 500: Error occurred when set expire  to redis.
	ErrRedisSetExpire
	// ErrRedisRuntime - 500: Error occurred when invoke redis api.
	ErrRedisRuntime
	// ErrRedisDataExpire - 500: Error occurred when redis data is expired.
	ErrRedisDataExpire
)
