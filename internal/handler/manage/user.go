package manage

import (
	svc "git.woda.ink/Woda_OA/internal/service/manage"
	"git.woda.ink/Woda_OA/pkg/response"
	"github.com/gin-gonic/gin"
)

// Login handles admin login request
func Login(c *gin.Context) {
	var req svc.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ArgError(c, "用户名和密码不能为空")
		return
	}

	data, err := svc.Login(&req)
	if err != nil {
		response.BizErr(c, err)
		return
	}

	response.SuccessMsg(c, data, "登录成功")
}
