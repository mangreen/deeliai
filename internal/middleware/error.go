package middleware

import (
	"github.com/gin-gonic/gin"
)

// ErrorMiddleware 是一個處理服務層錯誤的 Gin Middleware
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // 先執行後面的 handler

		// 如果 context 中有錯誤
		err := c.Errors.Last()
		if err == nil {
			return
		}

		// 在此處可以根據 err.Type 和 err.Meta 決定 status code
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
	}
}
