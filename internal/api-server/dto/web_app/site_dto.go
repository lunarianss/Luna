// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

import dto "github.com/lunarianss/Luna/internal/api-server/dto/app"

type GetWebSiteResponse struct {
	AppID      string          `json:"app_id"`
	EndUserID  string          `json:"end_user_id"`
	EnableSite int             `json:"enable_site"`
	Site       *dto.SiteDetail `json:"site"`
	Plan       string          `json:"plan"`
}
