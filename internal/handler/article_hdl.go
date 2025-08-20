package handler

import (
	"net/http"
	"strconv"

	"deeliai/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ArticleHandler struct {
	articleService *service.ArticleService
}

func NewArticleHandler(s *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleService: s}
}

type PostArticleRequest struct {
	URL string `json:"url" binding:"required,url"`
}

// PostArticle 處理 POST /articles 請求
func (h *ArticleHandler) PostArticle(c *gin.Context) {
	var req PostArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emailAny, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	article, err := h.articleService.CreateArticle(c.Request.Context(), req.URL, emailAny.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create article"})
		return
	}

	c.JSON(http.StatusCreated, article)
}

// GetArticles 處理 GET /articles 請求 (支援分頁)
func (h *ArticleHandler) GetArticles(c *gin.Context) {
	emailAny, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	articles, err := h.articleService.GetArticles(c.Request.Context(), emailAny.(string), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get articles"})
		return
	}

	c.JSON(http.StatusOK, articles)
}

// DeleteArticle 處理 DELETE /articles/:id 請求
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	articleID := c.Param("id")
	emailAny, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	articleUUID, err := uuid.Parse(articleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid article id"})
		return
	}

	err = h.articleService.DeleteArticle(c.Request.Context(), articleUUID, emailAny.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
