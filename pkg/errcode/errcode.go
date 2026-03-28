package errcode

import "fmt"

// BizError represents a business error with code and message
type BizError struct {
	Code    int
	Message string
}

func (e *BizError) Error() string {
	return e.Message
}

// New creates a new BizError
func New(code int, message string) *BizError {
	return &BizError{Code: code, Message: message}
}

// Newf creates a new BizError with formatted message
func Newf(code int, format string, args ...any) *BizError {
	return &BizError{Code: code, Message: fmt.Sprintf(format, args...)}
}

// Common error codes
const (
	CodeSuccess  = 200
	CodeArgError = 80000 // 参数错误
	CodeAuthFail = 80001 // 鉴权失败
	CodeBizError = 80002 // 业务错误
	CodeSysError = 80003 // 系统错误
)

// Pre-defined errors
var (
	ErrParam         = New(CodeArgError, "参数错误")
	ErrAuth          = New(CodeAuthFail, "请先登录")
	ErrAuthExpired   = New(CodeAuthFail, "登录已过期，请重新登录")
	ErrUserNotFound  = New(CodeBizError, "用户不存在")
	ErrLoginFailed   = New(CodeBizError, "用户名或密码错误")
	ErrSystemError   = New(CodeSysError, "系统繁忙，请稍后重试")
)
