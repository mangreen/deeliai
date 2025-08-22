package model

import (
	"time"

	"github.com/google/uuid"
)

type Article struct {
	ID           uuid.UUID `db:"id" json:"id"`
	UserEmail    string    `db:"user_email" json:"user_email"`
	URL          string    `db:"url" json:"url"`
	Title        *string   `db:"title" json:"title,omitempty"`
	Description  *string   `db:"description" json:"description,omitempty"`
	ImageURL     *string   `db:"image_url" json:"image_url,omitempty"`
	ScrapeStatus string    `db:"scrape_status" json:"scrape_status"`
	RetryCount   int       `db:"retry_count" json:"-"` // 不顯示給使用者
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
