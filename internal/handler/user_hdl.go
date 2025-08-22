package handler

import (
	"errors"
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

// @Summary 註冊新使用者
// @Description 使用者註冊一個新帳號
// @Tags users
// @Accept json
// @Produce json
// @Param request body SignupRequest true "註冊請求"
// @Success 201 {object} StandardResponse{data=model.User}
// @Failure 400 {object} ErrorResponse "無效的請求"
// @Failure 500 {object} ErrorResponse "內部伺服器錯誤"
// @Router /signup [post]
func (h *UserHandler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	user, err := h.userService.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, err, "Failed to create user")
		return
	}

	RespondWithSuccess(c, http.StatusCreated, "SignUp success", user)
}

// LoginRequest 登入請求的資料結構
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// @Summary 使用者登入
// @Description 使用者憑 E-mail 和密碼登入
// @Tags users
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登入請求"
// @Success 200 {object} StandardResponse{data=object{token=string}}
// @Failure 400 {object} ErrorResponse "無效的請求"
// @Failure 401 {object} ErrorResponse "憑證無效"
// @Failure 500 {object} ErrorResponse "內部伺服器錯誤"
// @Router /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	user, err := h.userService.Authenticate(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		// 建議回傳通用的錯誤訊息，避免暴露使用者不存在等細節
		RespondWithError(c, http.StatusUnauthorized, err, "Invalid email or password")
		return
	}

	token, err := h.AuthService.GenerateToken(user.Email)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	RespondWithSuccess(c, http.StatusOK, "Login success", gin.H{"token": token})
}

// @Summary 獲取使用者個人資料
// @Description 透過 JWT 驗證獲取使用者個人資料
// @Tags users
// @Security BearerAuth
// @Param Authorization header string true "JWT token" default(Bearer <your_JWT_token>)
// @Produce json
// @Success 200 {object} StandardResponse{data=model.User} "成功獲取使用者資料"
// @Failure 401 {object} ErrorResponse "未授權，JWT 驗證失敗"
// @Failure 500 {object} ErrorResponse "內部伺服器錯誤"
// @Router /me [get]
func (h *UserHandler) Me(c *gin.Context) {
	// 從 context 中取出我們在 middleware 存入的使用者 ID
	emailAny, exists := c.Get("email")
	if !exists {
		RespondWithError(c, http.StatusUnauthorized, errors.New(""), "User not authenticated")
		return
	}

	// 根據 email 查詢使用者資訊
	user, err := h.userService.FindByEmail(c.Request.Context(), emailAny.(string))
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	RespondWithSuccess(c, http.StatusOK, "Login success", user)
}
