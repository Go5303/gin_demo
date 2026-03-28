package middleware

import (
	"encoding/json"
	"strings"

	"git.woda.ink/Woda_OA/pkg/cache"
	"git.woda.ink/Woda_OA/pkg/response"
	"github.com/gin-gonic/gin"
)

// AdminInfo stores admin user information
type AdminInfo struct {
	ID   int    `json:"id"`
	GID  int    `json:"gid"`
	Name string `json:"name"`
	SP   int    `json:"sp"`
}

const AdminInfoKey = "admin_info"

// ManageAuth middleware - token from Header, validate against Redis
func ManageAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip auth for login routes
		if strings.HasSuffix(c.FullPath(), "/login") {
			c.Next()
			return
		}
		if strings.HasSuffix(c.FullPath(), "/manage/excel/demo") {
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
		val, err := cache.Get("manage:token:" + token)
		if err != nil || val == "" {
			response.AuthError(c, "登录已过期，请重新登录")
			c.Abort()
			return
		}

		// Parse admin info
		var admin AdminInfo
		if err := json.Unmarshal([]byte(val), &admin); err != nil {
			response.AuthError(c, "无效的登录信息")
			c.Abort()
			return
		}

		c.Set(AdminInfoKey, &admin)
		c.Next()
	}
}

// GetAdminInfo extracts admin info from context
func GetAdminInfo(c *gin.Context) *AdminInfo {
	val, exists := c.Get(AdminInfoKey)
	if !exists {
		return nil
	}
	return val.(*AdminInfo)
}
