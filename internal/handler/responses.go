package handler

import (
	"github.com/gin-gonic/gin"
)

// StandardResponse 定義所有成功回應的標準格式
type StandardResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse 定義所有錯誤回應的標準格式
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// RespondWithSuccess 封裝成功回應
func RespondWithSuccess(c *gin.Context, httpStatus int, message string, data interface{}) {
	c.JSON(httpStatus, StandardResponse{
		Message: message,
		Data:    data,
	})
}

// RespondWithError 封裝錯誤回應
func RespondWithError(c *gin.Context, httpStatus int, err error, message string) {
	c.JSON(httpStatus, ErrorResponse{
		Error:   err.Error(),
		Message: message,
	})
}
