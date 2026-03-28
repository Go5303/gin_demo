package middleware

import (
	"log"
	"net/http"

	"github.com/Go5303/gin_demo/pkg/response"
	"github.com/gin-gonic/gin"
)

// Recovery catches panics and returns a JSON error response
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[Recovery] panic: %v", err)
				c.AbortWithStatusJSON(http.StatusOK, response.R{
					Code:    response.CodeError,
					Message: "服务器内部错误",
					Data:    struct{}{},
				})
			}
		}()
		c.Next()
	}
}
