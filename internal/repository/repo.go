package repository

import (
	"context"
	"deeliai/internal/model"

	"github.com/google/uuid"
)

// UserRepository 定義了使用者資料的存取方法
// 這是 Service 層唯一需要知道的 "契約"
type UserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
}

type ArticleRepository interface {
	Create(ctx context.Context, article *model.Article) (*model.Article, error)
	UpdateMetadata(ctx context.Context, articleID uuid.UUID, title, description, imageURL string) error
	MarkScrapeFailed(ctx context.Context, articleID uuid.UUID) error
	ListByUserEmail(ctx context.Context, userEmail string, limit, offset int) ([]model.Article, error)
	FindByID(ctx context.Context, articleID uuid.UUID) (*model.Article, error)
	FindByIDAndUserEmail(ctx context.Context, articleID uuid.UUID, userEmail string) (*model.Article, error)
	Delete(ctx context.Context, articleID uuid.UUID, userEmail string) error
	FindFailedScrapes(ctx context.Context) ([]model.Article, error)
}
