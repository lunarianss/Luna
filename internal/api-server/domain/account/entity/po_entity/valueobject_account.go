// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package po_entity

type AccountStatus string

const (
	PENDING       AccountStatus = "pending"
	UNINITIALIZED AccountStatus = "uninitialized"
	ACTIVE        AccountStatus = "active"
	BANNED        AccountStatus = "banned"
	CLOSED        AccountStatus = "closed"
)

type TenantStatus string

const (
	TENANT_NORMAL TenantStatus = "normal"
	ARCHIVE       TenantStatus = "archive"
)

type TenantAccountJoinRole string

const (
	OWNER            TenantAccountJoinRole = "owner"
	EDITOR           TenantAccountJoinRole = "editor"
	ADMIN            TenantAccountJoinRole = "admin"
	NORMAL           TenantAccountJoinRole = "normal"
	DATASET_OPERATOR TenantAccountJoinRole = "dataset_operator"
)
