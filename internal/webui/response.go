package webui

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构体
type Response[T any] struct {
	Code int    `json:"code"`
	Data T      `json:"data,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

// 响应状态码
const (
	CodeSuccess = 0
	CodeError   = 1
)

// Success 返回成功响应
func Success[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, Response[T]{
		Code: CodeSuccess,
		Data: data,
	})
}

// SuccessWithStatus 返回成功响应，可指定HTTP状态码
func SuccessWithStatus[T any](c *gin.Context, statusCode int, data T) {
	c.JSON(statusCode, Response[T]{
		Code: CodeSuccess,
		Data: data,
	})
}

// Error 返回错误响应
func Error(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response[any]{
		Code: CodeError,
		Msg:  msg,
	})
}

// ErrorWithStatus 返回错误响应，可指定HTTP状态码
func ErrorWithStatus(c *gin.Context, statusCode int, msg string) {
	c.JSON(statusCode, Response[any]{
		Code: CodeError,
		Msg:  msg,
	})
}

// BadRequest 返回400错误
func BadRequest(c *gin.Context, msg string) {
	ErrorWithStatus(c, http.StatusBadRequest, msg)
}

// NotFound 返回404错误
func NotFound(c *gin.Context, msg string) {
	ErrorWithStatus(c, http.StatusNotFound, msg)
}

// InternalError 返回500错误
func InternalError(c *gin.Context, msg string) {
	ErrorWithStatus(c, http.StatusInternalServerError, msg)
}
