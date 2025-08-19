package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"deeliai/config"
	"deeliai/internal/handler"
	"deeliai/internal/repository/sqlximpl"
	"deeliai/internal/service"

	"github.com/MatusOllah/slogcolor"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	gin.ForceConsoleColor()

	// 1. 初始化結構化日誌
	slog.SetDefault(slog.New(slogcolor.NewHandler(os.Stderr, slogcolor.DefaultOptions)))

	// 2. 載入設定
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}
	slog.Info("Configuration loaded successfully")

	// 3. 初始化資料庫連線
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode)

	db, err := sqlx.Connect(cfg.Database.Driver, dsn)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// 4. 依賴注入：組裝 Repository, Service, Handler
	userRepo := sqlximpl.NewUserRepository(db)
	// articleRepo := repository.NewArticleRepository(db)

	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(cfg.App.JWTSecret)
	// articleService := service.NewArticleService(articleRepo)

	userHandler := handler.NewUserHandler(userService, authService)
	// articleHandler := handler.NewArticleHandler(articleService)

	// 5. 設定路由
	router := handler.SetupRouter(userHandler)
	slog.Info("Router setup complete")

	// 6. 建立 HTTP Server
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.App.Port),
		Handler: router,
	}

	// 7. 實現 Graceful Shutdown
	// 在一個新的 goroutine 中啟動 server，避免阻塞
	go func() {
		slog.Info(fmt.Sprintf("Server starting on port %d", cfg.App.Port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// 等待中斷訊號 (SIGINT or SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 阻塞直到接收到訊號
	slog.Info("Shutting down server...")

	// 建立一個有超時的 context，例如 5 秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 呼叫 server.Shutdown() 進行優雅關閉
	// 這會等待正在處理的請求結束，但不再接受新請求
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown:", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exiting gracefully.")
}
