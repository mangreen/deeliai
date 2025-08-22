package handler

import (
	"errors"
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

// @Summary 提交新文章
// @Description 提交一個文章 URL，後台會自動爬取 metadata
// @Tags articles
// @Security BearerAuth
// @Param Authorization header string true "JWT token" default(Bearer <your_JWT_token>)
// @Param request body PostArticleRequest true "文章 URL"
// @Accept json
// @Produce json
// @Success 202 {object} StandardResponse{data=model.Article} "文章正在處理中"
// @Failure 400 {object} ErrorResponse "無效的請求或 URL"
// @Failure 401 {object} ErrorResponse "未授權"
// @Failure 500 {object} ErrorResponse "內部伺服器錯誤"
// @Router /articles [post]
func (h *ArticleHandler) PostArticle(c *gin.Context) {
	var req PostArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	emailAny, exists := c.Get("email")
	if !exists {
		RespondWithError(c, http.StatusUnauthorized, errors.New(""), "User not authenticated")
		return
	}

	article, err := h.articleService.CreateArticle(c.Request.Context(), req.URL, emailAny.(string))
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	RespondWithSuccess(c, http.StatusCreated, "Post success", article)
}

// @Summary 獲取文章列表
// @Description 獲取使用者收藏的文章列表
// @Tags articles
// @Security BearerAuth
// @Param Authorization header string true "JWT token" default(Bearer <your_JWT_token>)
// @Produce json
// @Param limit query int false "限制返回數量" default(10)
// @Param offset query int false "跳過數量" default(0)
// @Success 200 {object} StandardResponse{data=[]model.Article} "成功獲取文章列表"
// @Failure 401 {object} ErrorResponse "未授權"
// @Failure 500 {object} ErrorResponse "內部伺服器錯誤"
// @Router /articles [get]
func (h *ArticleHandler) GetArticles(c *gin.Context) {
	emailAny, exists := c.Get("email")
	if !exists {
		RespondWithError(c, http.StatusUnauthorized, errors.New(""), "User not authenticated")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	articles, err := h.articleService.GetArticles(c.Request.Context(), emailAny.(string), page, limit)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	RespondWithSuccess(c, http.StatusOK, "Get success", articles)
}

// @Summary 刪除文章
// @Description 刪除使用者收藏的指定文章
// @Tags articles
// @Security BearerAuth
// @Param Authorization header string true "JWT token" default(Bearer <your_JWT_token>)
// @Param id path string true "文章 ID"
// @Success 204 "文章成功刪除"
// @Failure 400 {object} ErrorResponse "無效的文章 ID"
// @Failure 401 {object} ErrorResponse "未授權"
// @Failure 404 {object} ErrorResponse "文章不存在"
// @Failure 500 {object} ErrorResponse "內部伺服器錯誤"
// @Router /articles/{id} [delete]
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	articleID := c.Param("id")
	emailAny, exists := c.Get("email")
	if !exists {
		RespondWithError(c, http.StatusUnauthorized, errors.New(""), "User not authenticated")
		return
	}

	articleUUID, err := uuid.Parse(articleID)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, err, "Invalid article id")
		return
	}

	err = h.articleService.DeleteArticle(c.Request.Context(), articleUUID, emailAny.(string))
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	RespondWithSuccess(c, http.StatusOK, "Delete success", nil)
}
