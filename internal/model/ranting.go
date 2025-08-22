package model

import (
	"time"

	"github.com/google/uuid"
)

type Rating struct {
	ID        uuid.UUID `db:"id" json:"id"`
	UserEmail string    `db:"user_email" json:"user_email"`
	ArticleID uuid.UUID `db:"article_id" json:"article_id"`
	Scores    int       `db:"scores" json:"scores"`
	Tags      []string  `db:"tags" json:"tags"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
