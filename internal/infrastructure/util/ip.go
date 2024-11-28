// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package util

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func ExtractRemoteIP(c *gin.Context) string {
	if cfConnectingIP := c.GetHeader("CF-Connecting-IP"); cfConnectingIP != "" {
		return cfConnectingIP
	}

	if xForwardedFor := c.GetHeader("X-Forwarded-For"); xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	return c.ClientIP()
}
