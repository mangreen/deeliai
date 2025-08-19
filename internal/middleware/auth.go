package middleware

import (
	"deeliai/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 驗證 JWT 並將使用者 ID 存入 Gin context
func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims, err := authService.ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 將使用者 ID 存入 Gin context，以便後續的 handler 使用
		c.Set("email", claims.Email)
		c.Next()
	}
}
