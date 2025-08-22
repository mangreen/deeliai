// internal/handler/recommendation.go
package handler

import (
	"deeliai/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RecommendHandler struct {
	recService *service.RecommendService
}

func NewRecommendHandler(s *service.RecommendService) *RecommendHandler {
	return &RecommendHandler{recService: s}
}

// GetRecommendations 處理 GET /recommendations 請求
func (h *RecommendHandler) GetRecommendations(c *gin.Context) {
	emailAny, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	recommendations, err := h.recService.GetSimpleRecommendations(c.Request.Context(), emailAny.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get recommendations"})
		return
	}

	c.JSON(http.StatusOK, recommendations)
}
