package model

type ArticleScore struct {
	Article
	Score int `db:"score"`
}
