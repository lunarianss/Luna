// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package code

const (
	// ErrAppMapMode - 500: Error occurred while attempt to index from appTemplate using mode.
	ErrAppMapMode int = iota + 110101
	// ErrAppNotFoundRelatedConfig - 500: App config is not found.
	ErrAppNotFoundRelatedConfig
	// ErrAppStatusNotNormal - 500: App is not active.
	ErrAppStatusNotNormal
	// ErrAppCodeNotFound - 500: App code not found.
	ErrAppCodeNotFound
	// ErrAppSiteDisabled - 400: Site is disabled.
	ErrAppSiteDisabled
	// ErrAppApiDisabled - 400: Api is disabled.
	ErrAppApiDisabled
	// ErrAppTokenExceed - 400: Count of app token is exceeded.
	ErrAppTokenExceed
)
