package manage

import (
	svc "github.com/Go5303/gin_demo/internal/service/manage"
	"github.com/Go5303/gin_demo/pkg/response"
	"github.com/gin-gonic/gin"
)

// Index handles manage index request with Redis demo
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
