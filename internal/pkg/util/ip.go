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
