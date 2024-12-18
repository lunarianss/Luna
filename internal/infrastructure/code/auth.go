// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package code

const (
	// ErrEmailCode - 500: Error occurred when email code is incorrect.
	ErrEmailCode int = iota + 110301
	// ErrTokenEmail - 500: Error occurred when email is incorrect.
	ErrTokenEmail
	// ErrTenantAlreadyExist - 500: Error occurred when tenant is already exist.
	ErrTenantAlreadyExist
	// ErrAccountBanned - 500: Error occurred when user is banned but still to operate.
	ErrAccountBanned
	// ErrTenantStatusArchive - 400: Error occurred when tenant's status is archive.
	ErrTenantStatusArchive
)
