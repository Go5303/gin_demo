package mobile

import (
	svc "git.woda.ink/Woda_OA/internal/service/mobile"
	"git.woda.ink/Woda_OA/pkg/response"
	"github.com/gin-gonic/gin"
)

// Index handles mobile index request with Redis demo
func Index(c *gin.Context) {
	var req svc.IndexReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ArgError(c, "参数错误")
		return
	}

	data, err := svc.Index(&req)
	if err != nil {
		response.BizErr(c, err)
		return
	}

	response.Success(c, data)
}
