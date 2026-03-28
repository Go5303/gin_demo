package response

import (
	"errors"
	"net/http"

	"github.com/Go5303/gin_demo/pkg/errcode"
	"github.com/gin-gonic/gin"
)

const (
	CodeSuccess  = 200
	CodeArgError = 80000 // ERROR_CODE_ARG
	CodeAuthFail = 80001 // ERROR_CODE_AUTH
	CodeError    = 80002 // General error
)

// R is the standard API response structure
type R struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success
func Success(c *gin.Context, data interface{}) {
	if data == nil {
		data = struct{}{}
	}
	c.JSON(http.StatusOK, R{
		Code:    CodeSuccess,
		Message: "ok",
		Data:    data,
	})
}

// SuccessMsg returns a success response with custom message
func SuccessMsg(c *gin.Context, data interface{}, message string) {
	if data == nil {
		data = struct{}{}
	}
	c.JSON(http.StatusOK, R{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// Error
func Error(c *gin.Context, message string) {
	c.JSON(http.StatusOK, R{
		Code:    CodeError,
		Message: message,
		Data:    struct{}{},
	})
}

// ErrorCode returns an error response with custom code
func ErrorCode(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, R{
		Code:    code,
		Message: message,
		Data:    struct{}{},
	})
}

// ErrorData returns an error response with extend data
func ErrorData(c *gin.Context, code int, message string, data interface{}) {
	if data == nil {
		data = struct{}{}
	}
	c.JSON(http.StatusOK, R{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// ArgError returns parameter error (ERROR_CODE_ARG)
func ArgError(c *gin.Context, message string) {
	ErrorCode(c, CodeArgError, message)
}

// AuthError returns auth error (ERROR_CODE_AUTH)
func AuthError(c *gin.Context, message string) {
	ErrorCode(c, CodeAuthFail, message)
}

// BizErr handles error from service layer.
// If err is *errcode.BizError, uses its code; otherwise falls back to CodeError.
func BizErr(c *gin.Context, err error) {
	var bizErr *errcode.BizError
	if errors.As(err, &bizErr) {
		ErrorCode(c, bizErr.Code, bizErr.Message)
		return
	}
	Error(c, err.Error())
}
