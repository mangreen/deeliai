package handler

// SignupRequest 定義了建立使用者時的請求體結構
type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type PostArticleRequest struct {
	URL string `json:"url" binding:"required,url"`
}

type RateArticleRequest struct {
	Scores int      `json:"scores" binding:"required,gte=1,lte=5"`
	Tags   []string `json:"tags" binding:"required"`
}
