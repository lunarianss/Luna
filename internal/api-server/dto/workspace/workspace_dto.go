// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

import "github.com/lunarianss/Luna/internal/infrastructure/field"

type CurrentTenantInfo struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Plan           string         `json:"plan"`
	Status         string         `json:"status"`
	CreateAt       int64          `json:"create_at"`
	InTrail        bool           `json:"in_trail"`
	TrialEndReason string         `json:"trial_end_reason"`
	Role           string         `json:"role"`
	CustomConfig   map[string]any `json:"custom_config"`
	Current        field.BitBool  `json:"current"`
}

type CurrentTenant struct {
	Workspaces []*CurrentTenantInfo `json:"workspaces"`
}
