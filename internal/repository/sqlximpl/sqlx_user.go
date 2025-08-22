package sqlximpl

import (
	"context"
	"deeliai/internal/interfaces"
	"deeliai/internal/model"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type sqlxUserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) interfaces.UserRepository {
	return &sqlxUserRepository{db: db}
}

func (r *sqlxUserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	newUser := &model.User{}
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *`
	// 對於支援 RETURNING 的資料庫 (如 PostgreSQL)，可以這樣取回 ID
	// 對於 MySQL，需要用 LastInsertId()
	err := r.db.QueryRowxContext(ctx, query, user.Email, user.Password).StructScan(newUser)
	if err != nil {
		slog.Error("Failed to create user", "error", err)
		return nil, err
	}

	return newUser, nil
}

func (r *sqlxUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT email, password, created_at, updated_at FROM users WHERE email=$1`
	err := r.db.GetContext(ctx, user, query, email)
	if err != nil {
		slog.Error("Failed to get user by email", "error", err)
		return nil, err
	}

	return user, nil
}
