// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package po_entity

type ProviderType string

type ProviderQuotaType string

const (
	CUSTOM ProviderType = "custom"
	SYSTEM ProviderType = "system"
)

const (
	PAID ProviderQuotaType = "paid"

	FREE ProviderQuotaType = "free"

	TRIAL ProviderQuotaType = "trial"
)
