package handler

import (
	"log/slog"
	"net/http"

	"deeliai/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userService *service.UserService
	AuthService *service.AuthService
	validate    *validator.Validate
}

func NewUserHandler(userSvc *service.UserService, authSvc *service.AuthService) *UserHandler {
	return &UserHandler{
		userService: userSvc,
		AuthService: authSvc,
		validate:    validator.New(),
	}
}

// CreateUserRequest 定義了建立使用者時的請求體結構
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// 驗證請求參數
	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
		return
	}

	user, err := h.userService.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		slog.Error("Failed to create user", "error", err)
		c.Error(err) // 將錯誤傳給 ErrorMiddleware
		return
	}

	c.JSON(http.StatusCreated, user)
}

// LoginRequest 登入請求的資料結構
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginResponse 登入成功後的回應
type LoginResponse struct {
	Token string `json:"token"`
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := h.userService.Authenticate(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		// 建議回傳通用的錯誤訊息，避免暴露使用者不存在等細節
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := h.AuthService.GenerateToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: token})
}

func (h *UserHandler) Me(c *gin.Context) {
	// 從 context 中取出我們在 middleware 存入的使用者 ID
	emailAny, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 根據 email 查詢使用者資訊
	user, err := h.userService.FindByEmail(c.Request.Context(), emailAny.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
