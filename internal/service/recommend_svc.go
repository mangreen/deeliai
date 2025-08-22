package service

import (
	"context"
	"deeliai/internal/interfaces"
	"deeliai/internal/model"
)

type RecommendService struct {
	articleRepo interfaces.ArticleRepository
	ratingRepo  interfaces.RatingRepository
}

func NewRecommendService(articleRepo interfaces.ArticleRepository, ratingRepo interfaces.RatingRepository) *RecommendService {
	return &RecommendService{
		articleRepo: articleRepo,
		ratingRepo:  ratingRepo,
	}
}

// GetSimpleRecommendations 實現簡單推薦演算法（結合加權標籤）
func (s *RecommendService) GetSimpleRecommendations(ctx context.Context, userEmail string) ([]model.Article, error) {
	// 從評分加權中獲取使用者偏好標籤
	articleScores, err := s.articleRepo.ListRecommendArticles(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	// 如果沒有高評分文章，則推薦最新的熱門文章
	if len(articleScores) == 0 {
		return s.articleRepo.FindLatestArticles(ctx, userEmail, 10)
	}

	// 呼叫新的 FindRelatedArticles 函式
	return articleScores, nil
}
