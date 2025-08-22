package sqlximpl

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"deeliai/internal/model"
	"deeliai/internal/repository"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
		INSERT INTO ratings (user_email, article_id, scores, tags)
		SELECT 
			$1::text,      -- user_email
			$2::uuid,      -- article_id
			$3::int,       -- scores
			$4::text[]     -- tags
		WHERE EXISTS (
			SELECT 1 FROM articles 
			WHERE id = $2::uuid AND user_email = $1::text AND scrape_status = 'success'
		)
		ON CONFLICT (user_email, article_id) DO UPDATE
		SET scores = EXCLUDED.scores, updated_at = now()
		RETURNING id, user_email, article_id, scores, tags, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, rating.UserEmail, rating.ArticleID, rating.Scores, pq.Array(rating.Tags)).Scan(
		&createdRating.ID,
		&createdRating.UserEmail,
		&createdRating.ArticleID,
		&createdRating.Scores,
		pq.Array(&createdRating.Tags), // 🔑 這裡把 text[] 掃到 []string
		&createdRating.CreatedAt,
		&createdRating.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error("article not ready or not owned by user", "error", err)
			return nil, err
		}

		slog.Error("failed to create or update rating", "error", err)
		return nil, err
	}

	return &createdRating, nil
}

// FindRatingByUserEmailAndArticleID 取得使用者對單篇文章的評分
func (r *sqlxRatingRepository) FindRatingByUserEmailAndArticleID(ctx context.Context, userEmail string, articleID uuid.UUID) (*model.Rating, error) {
	var rating model.Rating
	query := `SELECT id, user_email, article_id, scores, tags, created_at, updated_at FROM ratings WHERE user_email = $1 AND article_id = $2 LIMIT 1`
	err := r.db.QueryRowxContext(ctx, query, userEmail, articleID).Scan(
		&rating.ID,
		&rating.UserEmail,
		&rating.ArticleID,
		&rating.Scores,
		pq.Array(&rating.Tags), // 🔑 這裡把 text[] 掃進 Go 的 []string
		&rating.CreatedAt,
		&rating.UpdatedAt,
	)
	if err != nil {
		slog.Error("failed to find rating", "error", err)
		return nil, err
	}

	return &rating, nil
}

// Delete 刪除使用者的評分
func (r *sqlxRatingRepository) Delete(ctx context.Context, userEmail string, articleID uuid.UUID) error {
	query := `DELETE FROM ratings WHERE user_email = $1 AND article_id = $2`
	result, err := r.db.ExecContext(ctx, query, userEmail, articleID)
	if err != nil {
		slog.Error("failed to delete rating", "error", err)
		return err
	}

	if rowsAffected, err := result.RowsAffected(); rowsAffected == 0 {
		slog.Error("rating not found", "error", err)
		return errors.New("rating not found")
	}

	return nil
}
