// internal/handler/recommendation.go
package handler

import (
	"deeliai/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RecommendHandler struct {
	recService *service.RecommendService
}

func NewRecommendHandler(s *service.RecommendService) *RecommendHandler {
	return &RecommendHandler{recService: s}
}

// @Summary 獲取文章推薦列表
// @Description 根據使用者的歷史評分與標籤，推薦相關的新文章
// @Tags recommendations
// @Security BearerAuth
// @Param Authorization header string true "JWT token" default(Bearer <your_JWT_token>)
// @Produce json
// @Success 200 {object} StandardResponse{data=[]model.Article} "成功獲取推薦文章列表"
// @Failure 401 {object} ErrorResponse "未授權，JWT 驗證失敗"
// @Failure 500 {object} ErrorResponse "內部伺服器錯誤"
// @Router /recommendations [get]
func (h *RecommendHandler) GetRecommendations(c *gin.Context) {
	emailAny, exists := c.Get("email")
	if !exists {
		RespondWithError(c, http.StatusUnauthorized, errors.New(""), "User not authenticated")
		return
	}

	recommendations, err := h.recService.GetSimpleRecommendations(c.Request.Context(), emailAny.(string))
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	RespondWithSuccess(c, http.StatusOK, "Get success", recommendations)
}
