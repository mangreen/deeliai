package handler

import (
	"errors"
	"net/http"

	"deeliai/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RatingHandler struct {
	ratingService *service.RatingService
}

func NewRatingHandler(s *service.RatingService) *RatingHandler {
	return &RatingHandler{ratingService: s}
}

// @Summary 評分並標記文章
// @Description 為指定文章評分並新增標籤
// @Tags ratings
// @Security BearerAuth
// @Param Authorization header string true "JWT token" default(Bearer <your_JWT_token>)
// @Accept json
// @Produce json
// @Param id path string true "文章 ID"
// @Param request body RateArticleRequest true "評分與標籤"
// @Success 200 {object} StandardResponse{data=model.Rating} "評分成功"
// @Failure 400 {object} ErrorResponse "無效的請求或評分值"
// @Failure 401 {object} ErrorResponse "未授權"
// @Failure 500 {object} ErrorResponse "內部伺服器錯誤"
// @Router /articles/{id}/rate [post]
func (h *RatingHandler) RateArticle(c *gin.Context) {
	articleID := c.Param("id")
	articleUUID, err := uuid.Parse(articleID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, err, "Invalid article id")
		return
	}

	emailAny, exists := c.Get("email")
	if !exists {
		RespondWithError(c, http.StatusUnauthorized, errors.New(""), "User not authenticated")
		return
	}

	var req RateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	if len(req.Tags) < 1 {
		RespondWithError(c, http.StatusBadRequest, err, "At least 1 tag")
		return
	}

	rating, err := h.ratingService.RateArticle(c.Request.Context(), emailAny.(string), articleUUID, req.Scores, req.Tags)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	RespondWithSuccess(c, http.StatusCreated, "Post success", rating)
}

// @Summary 獲取使用者對文章的評分
// @Description 獲取使用者對指定文章的評分與標籤
// @Tags ratings
// @Security BearerAuth
// @Param Authorization header string true "JWT token" default(Bearer <your_JWT_token>)
// @Produce json
// @Param id path string true "文章 ID"
// @Success 200 {object} StandardResponse{data=model.Rating} "成功獲取評分"
// @Failure 400 {object} ErrorResponse "無效的文章 ID"
// @Failure 401 {object} ErrorResponse "未授權"
// @Failure 404 {object} ErrorResponse "找不到評分"
// @Failure 500 {object} ErrorResponse "內部伺服器錯誤"
// @Router /articles/{id}/rate [get]
func (h *RatingHandler) GetRating(c *gin.Context) {
	articleID := c.Param("id")
	articleUUID, err := uuid.Parse(articleID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, err, "Invalid article id")
		return
	}

	emailAny, exists := c.Get("email")
	if !exists {
		RespondWithError(c, http.StatusUnauthorized, errors.New(""), "User not authenticated")
		return
	}

	rating, err := h.ratingService.GetRating(c.Request.Context(), emailAny.(string), articleUUID)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	RespondWithSuccess(c, http.StatusOK, "Get success", rating)
}

// @Summary 刪除評分
// @Description 刪除使用者對指定文章的評分與標籤
// @Tags ratings
// @Security BearerAuth
// @Param Authorization header string true "JWT token" default(Bearer <your_JWT_token>)
// @Param id path string true "文章 ID"
// @Success 204 "評分成功刪除"
// @Failure 400 {object} ErrorResponse "無效的文章 ID"
// @Failure 401 {object} ErrorResponse "未授權"
// @Failure 404 {object} ErrorResponse "找不到評分"
// @Failure 500 {object} ErrorResponse "內部伺服器錯誤"
// @Router /articles/{id}/rate [delete]
func (h *RatingHandler) DeleteRating(c *gin.Context) {
	articleID := c.Param("id")
	articleUUID, err := uuid.Parse(articleID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, err, "Invalid article id")
		return
	}

	emailAny, exists := c.Get("email")
	if !exists {
		RespondWithError(c, http.StatusUnauthorized, errors.New(""), "User not authenticated")
		return
	}

	err = h.ratingService.Delete(c.Request.Context(), emailAny.(string), articleUUID)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	RespondWithSuccess(c, http.StatusOK, "Delete success", nil)
}
