// internal/service/article.go
package service

import (
	"context"
	"fmt"

	"deeliai/internal/interfaces"
	"deeliai/internal/model"

	"github.com/google/uuid"
)

type ArticleService struct {
	articleRepo interfaces.ArticleRepository
	producer    interfaces.QueueProducer // 依賴介面
}

func NewArticleService(repo interfaces.ArticleRepository, producer interfaces.QueueProducer) *ArticleService {
	return &ArticleService{
		articleRepo: repo,
		producer:    producer,
	}
}

// CreateArticle 處理文章儲存和爬取任務分派
func (s *ArticleService) CreateArticle(ctx context.Context, url, userEmail string) (*model.Article, error) {
	article := &model.Article{
		UserEmail: userEmail,
		URL:       url,
	}

	// 1. 儲存文章到資料庫，狀態為 pending
	createdArticle, err := s.articleRepo.Create(ctx, article)
	if err != nil {
		return nil, err
	}

	// 2. 將文章 ID 推入爬取佇列，讓 worker 處理
	// 這裡直接呼叫 producer 的 Produce 方法，不關心底層是誰
	if err := s.producer.Produce(createdArticle.ID.String()); err != nil {
		return nil, fmt.Errorf("failed to produce message to queue: %w", err)
	}

	return createdArticle, nil
}

// GetArticles 取得使用者儲存的文章列表
func (s *ArticleService) GetArticles(ctx context.Context, userEmail string, page, limit int) ([]model.Article, error) {
	offset := (page - 1) * limit
	return s.articleRepo.ListByUserEmail(ctx, userEmail, limit, offset)
}

// DeleteArticle 刪除使用者收藏的文章
func (s *ArticleService) DeleteArticle(ctx context.Context, articleUUID uuid.UUID, userEmail string) error {
	return s.articleRepo.Delete(ctx, articleUUID, userEmail)
}
