package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 定義 JWT 中包含的資料
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type AuthService struct {
	JWTSecret string
}

func NewAuthService(JWTSecret string) *AuthService {
	return &AuthService{JWTSecret: JWTSecret}
}

// GenerateToken 根據使用者 ID 產生 JWT
func (s *AuthService) GenerateToken(email string) (string, error) {
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Token 有效期限 24 小時
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.JWTSecret))
}

// ParseToken 解析並驗證 JWT，成功則回傳 Claims
func (s *AuthService) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}
