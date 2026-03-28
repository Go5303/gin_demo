package middleware

import (
	"encoding/json"
	"strings"

	"github.com/Go5303/gin_demo/pkg/cache"
	"github.com/Go5303/gin_demo/pkg/response"
	"github.com/gin-gonic/gin"
)

// MobileUserInfo stores mobile user info from token
type MobileUserInfo struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
}

const MobileUserKey = "mobile_user_info"

// MobileAuth middleware - token from Header, validate against Redis
func MobileAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip auth for login routes
		if strings.HasSuffix(c.FullPath(), "/login") {
			c.Next()
			return
		}

		// Get token from header
		token := c.GetHeader("Token")
		if token == "" {
			response.AuthError(c, "请先登录")
			c.Abort()
			return
		}

		// Look up token in Redis
		val, err := cache.Get("mobile:token:" + token)
		if err != nil || val == "" {
			response.AuthError(c, "登录已过期，请重新登录")
			c.Abort()
			return
		}

		// Parse user info
		var user MobileUserInfo
		if err := json.Unmarshal([]byte(val), &user); err != nil {
			response.AuthError(c, "无效的登录信息")
			c.Abort()
			return
		}

		c.Set(MobileUserKey, &user)
		c.Next()
	}
}

// GetMobileUser extracts mobile user info from context
func GetMobileUser(c *gin.Context) *MobileUserInfo {
	val, exists := c.Get(MobileUserKey)
	if !exists {
		return nil
	}
	return val.(*MobileUserInfo)
}
