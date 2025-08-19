package repository

import (
	"context"
	"deeliai/internal/model"
)

// UserRepository 定義了使用者資料的存取方法
// 這是 Service 層唯一需要知道的 "契約"
type UserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}
