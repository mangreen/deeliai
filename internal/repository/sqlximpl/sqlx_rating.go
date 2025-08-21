package sqlximpl

import (
	"context"
	"fmt"

	"deeliai/internal/model"
	"deeliai/internal/repository"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type sqlxRatingRepository struct {
	db *sqlx.DB
}

func NewRatingRepository(db *sqlx.DB) repository.RatingRepository {
	return &sqlxRatingRepository{db: db}
}

// CreateOrUpdate 創建或更新使用者的評分
func (r *sqlxRatingRepository) CreateOrUpdate(ctx context.Context, rating *model.Rating) (*model.Rating, error) {
	var createdRating model.Rating
	// 使用 ON CONFLICT DO UPDATE 來處理 upsert (新增或更新)
	query := `
		INSERT INTO ratings (user_email, article_id, rating)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_email, article_id) DO UPDATE
		SET rating = EXCLUDED.rating, updated_at = now()
		RETURNING *
	`
	err := r.db.QueryRowxContext(ctx, query, rating.UserEmail, rating.ArticleID, rating.Rating).StructScan(&createdRating)
	if err != nil {
		return nil, fmt.Errorf("failed to create or update rating: %w", err)
	}

	return &createdRating, nil
}

// FindRatingByUserEmailAndArticleID 取得使用者對單篇文章的評分
func (r *sqlxRatingRepository) FindRatingByUserEmailAndArticleID(ctx context.Context, userEmail string, articleID uuid.UUID) (*model.Rating, error) {
	var rating model.Rating
	query := `SELECT id, user_email, article_id, rating, created_at FROM ratings WHERE user_email = $1 AND article_id = $2 LIMIT 1`
	err := r.db.GetContext(ctx, &rating, query, userEmail, articleID)
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

// Delete 刪除使用者的評分
func (r *sqlxRatingRepository) Delete(ctx context.Context, userEmail string, articleID uuid.UUID) error {
	query := `DELETE FROM ratings WHERE user_email = $1 AND article_id = $2`
	result, err := r.db.ExecContext(ctx, query, userEmail, articleID)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("rating not found")
	}
	return nil
}
