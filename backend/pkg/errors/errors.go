package errors

import (
	"fmt"
	"net/http"
)

// AppError 应用错误接口
// 所有自定义错误都实现这个接口
type AppError interface {
	error                 // 继承 Go 的标准 error 接口
	StatusCode() int      // HTTP 状态码
	ErrorCode() string    // 错误代码（用于前端识别错误类型）
	ErrorMessage() string // 错误消息
}

// ========== 基础错误结构 ==========

// baseError 基础错误结构体
type baseError struct {
	statusCode   int
	errorCode    string
	errorMessage string
}

// Error 实现 error 接口
func (e *baseError) Error() string {
	return e.errorMessage
}

// StatusCode 返回 HTTP 状态码
func (e *baseError) StatusCode() int {
	return e.statusCode
}

// ErrorCode 返回错误代码
func (e *baseError) ErrorCode() string {
	return e.errorCode
}

// ErrorMessage 返回错误消息
func (e *baseError) ErrorMessage() string {
	return e.errorMessage
}

// ========== 常见错误类型 ==========

// NewBadRequestError 创建 400 错误（请求参数错误）
func NewBadRequestError(message string) AppError {
	return &baseError{
		statusCode:   http.StatusBadRequest, // 400
		errorCode:    "BAD_REQUEST",
		errorMessage: message,
	}
}

// NewUnauthorizedError 创建 401 错误（未授权）
func NewUnauthorizedError(message string) AppError {
	return &baseError{
		statusCode:   http.StatusUnauthorized, // 401
		errorCode:    "UNAUTHORIZED",
		errorMessage: message,
	}
}

// NewForbiddenError 创建 403 错误（禁止访问）
func NewForbiddenError(message string) AppError {
	return &baseError{
		statusCode:   http.StatusForbidden, // 403
		errorCode:    "FORBIDDEN",
		errorMessage: message,
	}
}

// NewNotFoundError 创建 404 错误（资源不存在）
func NewNotFoundError(message string) AppError {
	return &baseError{
		statusCode:   http.StatusNotFound, // 404
		errorCode:    "NOT_FOUND",
		errorMessage: message,
	}
}

// NewConflictError 创建 409 错误（资源冲突，如用户名已存在）
func NewConflictError(message string) AppError {
	return &baseError{
		statusCode:   http.StatusConflict, // 409
		errorCode:    "CONFLICT",
		errorMessage: message,
	}
}

// NewInternalServerError 创建 500 错误（服务器内部错误）
func NewInternalServerError(message string) AppError {
	return &baseError{
		statusCode:   http.StatusInternalServerError, // 500
		errorCode:    "INTERNAL_SERVER_ERROR",
		errorMessage: message,
	}
}

// ========== 特定业务错误 ==========

// NewValidationError 创建验证错误
func NewValidationError(field string, message string) AppError {
	return &baseError{
		statusCode:   http.StatusBadRequest,
		errorCode:    "VALIDATION_ERROR",
		errorMessage: fmt.Sprintf("%s: %s", field, message),
	}
}

// NewDatabaseError 创建数据库错误
func NewDatabaseError(operation string, err error) AppError {
	return &baseError{
		statusCode:   http.StatusInternalServerError,
		errorCode:    "DATABASE_ERROR",
		errorMessage: fmt.Sprintf("数据库操作失败 [%s]: %v", operation, err),
	}
}

// ErrorResponse Swagger 错误响应结构
type ErrorResponse struct {
	Success bool      `json:"success" example:"false"`
	Error   ErrorInfo `json:"error"`
}

// ErrorInfo 错误信息
type ErrorInfo struct {
	Code    string `json:"code" example:"BAD_REQUEST"`
	Message string `json:"message" example:"请求参数错误"`
}
