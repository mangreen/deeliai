package sqlximpl

import (
	"context"
	"database/sql"
	"deeliai/internal/model"
	"deeliai/internal/repository"

	"github.com/jmoiron/sqlx"
)

type sqlxUserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repository.UserRepository {
	return &sqlxUserRepository{db: db}
}

func (r *sqlxUserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	newUser := &model.User{}
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING email, password, created_at, updated_at`
	// 對於支援 RETURNING 的資料庫 (如 PostgreSQL)，可以這樣取回 ID
	// 對於 MySQL，需要用 LastInsertId()
	err := r.db.QueryRowxContext(ctx, query, user.Email, user.Password).StructScan(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (r *sqlxUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT email, password, created_at, updated_at FROM users WHERE email=$1`
	err := r.db.GetContext(ctx, user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 或者回傳一個自訂的 not found error
		}
		return nil, err
	}
	return user, nil
}
