package middleware

import (
	"time"

	"github.com/Go5303/gin_demo/pkg/logger"
	"github.com/gin-gonic/gin"
)

// RequestLogger logs request info
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.String()

		logger.WithCtx(c).Infof("[HTTP] %d | %13v | %15s | %-7s %s",
			status, latency, clientIP, method, path)
	}
}
