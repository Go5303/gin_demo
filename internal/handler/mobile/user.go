package mobile

import (
	svc "git.woda.ink/Woda_OA/internal/service/mobile"
	"git.woda.ink/Woda_OA/pkg/logger"
	"git.woda.ink/Woda_OA/pkg/response"
	"github.com/gin-gonic/gin"
)

// Login handles mobile user login request
func Login(c *gin.Context) {
	var req svc.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.L.Warn("mobile login param error", "err", err)
		response.ArgError(c, "用户名和密码不能为空")
		return
	}

	logger.L.Debug("mobile login attempt", "username", req.Username)

	data, err := svc.Login(&req)
	if err != nil {
		logger.L.Error("mobile login failed", "username", req.Username, "err", err)
		response.BizErr(c, err)
		return
	}

	logger.L.Info("mobile login success", "username", req.Username)
	response.SuccessMsg(c, data, "登录成功")
}
