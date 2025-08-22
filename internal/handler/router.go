package handler

import (
	"deeliai/docs"
	"deeliai/internal/middleware"
	"log/slog"
	"time"

	_ "deeliai/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title DeeliAI API
// @version 1.0
// @description 一個用於文章收藏與推薦的 AI 應用程式。
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func SetupRouter(userHandler *UserHandler, articleHandler *ArticleHandler, ratingHandler *RatingHandler, recommendHandler *RecommendHandler) *gin.Engine {
	// gin.ReleaseMode or gin.DebugMode
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

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

	// Swagger 文件路由
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 路由分組
	r.POST("/signup", userHandler.Signup)
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
