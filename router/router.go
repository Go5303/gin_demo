package router

import (
	"git.woda.ink/Woda_OA/config"
	manageHandler "git.woda.ink/Woda_OA/internal/handler/manage"
	mobileHandler "git.woda.ink/Woda_OA/internal/handler/mobile"
	"git.woda.ink/Woda_OA/internal/middleware"
	"github.com/gin-gonic/gin"
)

// Setup initializes all routes
func Setup(r *gin.Engine, cfg *config.AppConfig) {
	// Global middleware
	r.Use(middleware.CORS())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.Recovery())

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// ========================================
	// Backend Management API (/manage)
	// B端 - 管理后台接口
	// ========================================
	manageGroup := r.Group("/manage")
	manageGroup.Use(middleware.ManageAuth())
	{
		// User module - 用户相关
		user := manageGroup.Group("/user")
		{
			user.POST("/login", manageHandler.Login) // 登录（免鉴权）
		}

		// Index module - 首页相关
		index := manageGroup.Group("/index")
		{
			index.POST("/index", manageHandler.Index) // 首页（需登录）
		}

		// Excel demo (keep existing)
		manageGroup.GET("/excel/demo", manageHandler.ExcelDemo)
	}

	// ========================================
	// Mobile API (/mobile)
	// C端 - 移动端接口
	// ========================================
	mobileGroup := r.Group("/mobile")
	mobileGroup.Use(middleware.MobileAuth())
	{
		// User module - 用户相关
		user := mobileGroup.Group("/user")
		{
			user.POST("/login", mobileHandler.Login) // 登录（免鉴权）
		}

		// Index module - 首页相关
		index := mobileGroup.Group("/index")
		{
			index.POST("/index", mobileHandler.Index) // 首页（需登录）
		}
	}
}
