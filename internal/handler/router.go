package handler

import (
	"deeliai/internal/middleware"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *UserHandler, articleHandler *ArticleHandler, ratingHandler *RatingHandler, recommendHandler *RecommendHandler) *gin.Engine {
	// gin.ReleaseMode or gin.DebugMode
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// 自訂的錯誤處理 Middleware
	r.Use(middleware.ErrorMiddleware())

	// 使用自訂的結構化日誌 Middleware
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		slog.Debug("request",
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.ClientIP(),
			"latency", latency,
			"user_agent", c.Request.UserAgent(),
		)
	})

	// 路由分組
	r.POST("/signup", userHandler.CreateUser)
	r.POST("/login", userHandler.Login)
	r.GET("/me", middleware.AuthMiddleware(userHandler.AuthService), userHandler.Me)

	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.AuthMiddleware(userHandler.AuthService))
	{
		// 文章收藏 API
		apiV1.POST("/articles", articleHandler.PostArticle)
		apiV1.GET("/articles", articleHandler.GetArticles)
		apiV1.DELETE("/articles/:id", articleHandler.DeleteArticle)

		apiV1.POST("/articles/:id/rate", ratingHandler.RateArticle)
		apiV1.GET("/articles/:id/rate", ratingHandler.GetRating)
		apiV1.DELETE("/articles/:id/rate", ratingHandler.DeleteRating)

		apiV1.GET("/recommendations", recommendHandler.GetRecommendations)
	}

	return r
}
