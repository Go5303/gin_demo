package mobile

import (
	svc "github.com/Go5303/gin_demo/internal/service/mobile"
	"github.com/Go5303/gin_demo/pkg/logger"
	"github.com/Go5303/gin_demo/pkg/response"
	"github.com/gin-gonic/gin"
)

// Index handles mobile index request with Redis demo
func Index(c *gin.Context) {
	var req svc.IndexReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.L.Warn("mobile index param error", "err", err)
		response.ArgError(c, "参数错误")
		return
	}

	data, err := svc.Index(&req)
	if err != nil {
		logger.L.Error("mobile index failed", "err", err)
		response.BizErr(c, err)
		return
	}

	logger.L.Debug("mobile index success")
	response.Success(c, data)
}
