package handler

import (
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

type RantingRequest struct {
	Rating int `json:"rating" binding:"required,gte=1,lte=5"`
}

// RateArticle 處理 POST /articles/:id/rate 請求
func (h *RatingHandler) RateArticle(c *gin.Context) {
	articleID := c.Param("id")
	articleUUID, err := uuid.Parse(articleID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid article id"})
		return
	}

	emailAny, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req RantingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating value, must be between 1 and 5"})
		return
	}

	rating, err := h.ratingService.RateArticle(c.Request.Context(), emailAny.(string), articleUUID, req.Rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rating)
}

// GetRating 處理 GET /articles/:id/rate 請求
func (h *RatingHandler) GetRating(c *gin.Context) {
	articleID := c.Param("id")
	articleUUID, err := uuid.Parse(articleID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid article id"})
		return
	}

	emailAny, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	rating, err := h.ratingService.GetRating(c.Request.Context(), emailAny.(string), articleUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rating not found"})
		return
	}

	c.JSON(http.StatusOK, rating)
}

// DeleteRating 處理 DELETE /articles/:id/rate 請求
func (h *RatingHandler) DeleteRating(c *gin.Context) {
	articleID := c.Param("id")
	articleUUID, err := uuid.Parse(articleID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid article id"})
		return
	}

	emailAny, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	err = h.ratingService.Delete(c.Request.Context(), emailAny.(string), articleUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
