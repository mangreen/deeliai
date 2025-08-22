package service

import (
	"context"
	"fmt"

	"deeliai/internal/interfaces"
	"deeliai/internal/model"

	"github.com/google/uuid"
)

type RatingService struct {
	ratingRepo interfaces.RatingRepository
}

func NewRatingService(repo interfaces.RatingRepository) *RatingService {
	return &RatingService{ratingRepo: repo}
}

// RateArticle 為文章評分
func (s *RatingService) RateArticle(ctx context.Context, userEmail string, articleUUID uuid.UUID, scores int, tags []string) (*model.Rating, error) {
	if scores < 1 || scores > 5 {
		return nil, fmt.Errorf("rating must be between 1 and 5")
	}

	rating := &model.Rating{
		UserEmail: userEmail,
		ArticleID: articleUUID,
		Scores:    scores,
		Tags:      tags,
	}

	return s.ratingRepo.CreateOrUpdate(ctx, rating)
}

// GetRating 取得使用者對文章的評分
func (s *RatingService) GetRating(ctx context.Context, userEmail string, articleUUID uuid.UUID) (*model.Rating, error) {
	return s.ratingRepo.FindRatingByUserEmailAndArticleID(ctx, userEmail, articleUUID)
}

// Delete 刪除使用者的評分
func (s *RatingService) Delete(ctx context.Context, userEmail string, articleUUID uuid.UUID) error {
	return s.ratingRepo.Delete(ctx, userEmail, articleUUID)
}
