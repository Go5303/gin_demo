package middleware

import (
	"github.com/Go5303/uuid"
	"github.com/gin-gonic/gin"
)

// RequestID generates a unique request ID for each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.Uuid()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}
