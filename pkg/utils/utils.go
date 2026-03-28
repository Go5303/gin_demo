package utils

import (
	"fmt"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetClientIP extracts the real client IP from request headers
func GetClientIP(c *gin.Context) string {
	// Check HTTP_CLIENT_IP
	if ip := c.GetHeader("Client-Ip"); ip != "" && ip != "unknown" {
		return sanitizeIP(ip)
	}
	// Check HTTP_X_FORWARDED_FOR
	if ip := c.GetHeader("X-Forwarded-For"); ip != "" && ip != "unknown" {
		// Take the first IP if multiple
		parts := strings.Split(ip, ",")
		return sanitizeIP(strings.TrimSpace(parts[0]))
	}
	// Fallback to RemoteAddr
	ip, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
	return sanitizeIP(ip)
}

func sanitizeIP(ip string) string {
	if net.ParseIP(ip) != nil {
		return ip
	}
	return ""
}

func GetEmptyValue(value, defaultVal string) string {
	if value == "" {
		return defaultVal
	}
	return value
}

func GetParam(c *gin.Context, key string) string {
	val := c.Query(key)
	if val == "" {
		val = c.PostForm(key)
	}
	return val
}

// GetParamInt gets an integer parameter
func GetParamInt(c *gin.Context, key string, defaultVal int) int {
	val := GetParam(c, key)
	if val == "" {
		return defaultVal
	}
	var result int
	_, _ = fmt.Sscanf(val, "%d", &result)
	return result
}
