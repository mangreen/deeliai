package sqlximpl

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"deeliai/internal/interfaces"
	"deeliai/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type sqlxArticleRepository struct {
	db *sqlx.DB
}

func NewArticleRepository(db *sqlx.DB) interfaces.ArticleRepository {
	return &sqlxArticleRepository{db: db}
}

// Create 將新文章記錄存入資料庫
func (r *sqlxArticleRepository) Create(ctx context.Context, article *model.Article) (*model.Article, error) {
	newArticle := &model.Article{}
	query := `INSERT INTO articles (user_email, url) VALUES ($1, $2) RETURNING *`
	// 對於支援 RETURNING 的資料庫 (如 PostgreSQL)，可以這樣取回 ID
	// 對於 MySQL，需要用 LastInsertId()
	err := r.db.QueryRowxContext(ctx, query, article.UserEmail, article.URL).StructScan(newArticle)
	if err != nil {
		slog.Error("Failed to create article", "error", err)
		return nil, err
	}
	return newArticle, nil
}

// UpdateMetadata 更新文章的 Metadata
func (r *sqlxArticleRepository) UpdateMetadata(ctx context.Context, articleID uuid.UUID, title, description, imageURL string) error {
	query := `UPDATE articles SET title=$1, description=$2, image_url=$3, scrape_status='success', updated_at=$4 WHERE id=$5`
	_, err := r.db.ExecContext(ctx, query, title, description, imageURL, time.Now(), articleID)
	if err != nil {
		slog.Error("Failed to update article metadata", "error", err)
		return err
	}

	return nil
}

// MarkScrapeFailed 標記爬取失敗並增加重試次數
func (r *sqlxArticleRepository) MarkScrapeFailed(ctx context.Context, articleID uuid.UUID) error {
	query := `UPDATE articles SET scrape_status='failed', retry_count=retry_count+1, updated_at=$1 WHERE id=$2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), articleID)
	if err != nil {
		slog.Error("Failed to marke scrape failed", "error", err)
		return err
	}

	return nil
}

// ListByUserEmail 根據使用者 ID 取得文章列表
func (r *sqlxArticleRepository) ListByUserEmail(ctx context.Context, userEmail string, limit, offset int) ([]model.Article, error) {
	var articles []model.Article
	query := `SELECT id, user_email, url, title, description, image_url, scrape_status, created_at FROM articles WHERE user_email = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	err := r.db.SelectContext(ctx, &articles, query, userEmail, limit, offset)
	if err != nil {
		slog.Error("Failed to list articles by email", "error", err)
		return nil, err
	}

	return articles, nil
}

// FindByID 根據文章 ID 取得單篇文章
func (r *sqlxArticleRepository) FindByID(ctx context.Context, articleID uuid.UUID) (*model.Article, error) {
	article := &model.Article{}
	query := `SELECT id, user_email, url, title, description, image_url, scrape_status, created_at FROM articles WHERE id = $1 LIMIT 1`
	err := r.db.GetContext(ctx, article, query, articleID)
	if err != nil {
		slog.Error("Failed to get article by id", "error", err)
		return nil, err
	}

	return article, nil
}

// FindByIDAndUserEmail 根據文章 ID 和使用者 ID 取得單篇文章
func (r *sqlxArticleRepository) FindByIDAndUserEmail(ctx context.Context, articleID uuid.UUID, userEmail string) (*model.Article, error) {
	article := &model.Article{}
	query := `SELECT id, user_email, url, title, description, image_url, scrape_status, created_at FROM articles WHERE id = $1 AND user_email = $2 LIMIT 1`
	err := r.db.GetContext(ctx, article, query, articleID, userEmail)
	if err != nil {
		slog.Error("Failed to get article by id & email", "error", err)
		return nil, err
	}

	return article, nil
}

// Delete 刪除文章
func (r *sqlxArticleRepository) Delete(ctx context.Context, articleID uuid.UUID, userEmail string) error {
	query := `DELETE FROM articles WHERE id = $1 AND user_email = $2`
	res, err := r.db.ExecContext(ctx, query, articleID, userEmail)
	if err != nil {
		slog.Error("Failed to delete article", "error", err)
		return err
	}

	if rowsAffected, err := res.RowsAffected(); rowsAffected == 0 {
		slog.Error("article not found or user not authorized", "error", err)
		return errors.New("article not found or user not authorized")
	}

	return nil
}

// FindFailedScrapes 尋找失敗且重試次數未達上限的文章
func (r *sqlxArticleRepository) FindFailedScrapes(ctx context.Context) ([]model.Article, error) {
	var articles []model.Article
	query := `SELECT id, url, retry_count FROM articles WHERE scrape_status = 'failed' AND retry_count < 3`
	err := r.db.SelectContext(ctx, &articles, query)
	if err != nil {
		slog.Error("Failed to get failed scrapes", "error", err)
		return nil, err
	}

	return articles, nil
}

func (r *sqlxArticleRepository) ListRecommendArticles(ctx context.Context, userEmail string) ([]model.Article, error) {
	query := `
        WITH user_tag_weights AS (
            SELECT unnest(tags) AS tag, SUM(scores) AS weight
            FROM ratings
            WHERE user_email = $1
            GROUP BY tag
        )
        SELECT a.*, COALESCE(SUM(t.weight), 0) AS score
        FROM articles a
        JOIN ratings r ON a.id = r.article_id
        JOIN LATERAL unnest(r.tags) AS rt(tag) ON TRUE
        LEFT JOIN user_tag_weights t ON rt.tag = t.tag
        WHERE NOT EXISTS (
            SELECT 1 FROM ratings r2
            WHERE r2.article_id = a.id
              AND r2.user_email = $1
        )
        GROUP BY a.id
        ORDER BY score DESC
        LIMIT 10
    `

	var articles []interfaces.ArticleScore
	err := r.db.SelectContext(ctx, &articles, query, userEmail)
	if err != nil {
		slog.Error("Failed to find recommend articles", "error", err)
		return nil, err
	}

	result := make([]model.Article, len(articles))
	for i, a := range articles {
		result[i] = a.Article
	}

	return result, nil
}

// FindLatestArticles 找出最新的文章
func (r *sqlxArticleRepository) FindLatestArticles(ctx context.Context, userEmail string, limit int) ([]model.Article, error) {
	var articles []model.Article
	query := `SELECT id, url, title, description, image_url FROM articles WHERE user_email != $1 ORDER BY created_at DESC LIMIT $2`
	err := r.db.SelectContext(ctx, &articles, query, userEmail, limit)
	if err != nil {
		slog.Error("Failed to find latest articles", "error", err)
		return nil, err
	}

	return articles, nil
}
