package middleware

import (
	"context"
	"time"

	"github.com/Go5303/gin_demo/pkg/response"
	"github.com/gin-gonic/gin"
)

// Timeout returns a middleware that aborts the request if it exceeds the given duration.
func Timeout(seconds int) gin.HandlerFunc {
	timeout := time.Duration(seconds) * time.Second

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{}, 1)
		go func() {
			c.Next()
			done <- struct{}{}
		}()

		select {
		case <-done:
			return
		case <-ctx.Done():
			c.Abort()
			response.Error(c, "request timeout")
		}
	}
}
