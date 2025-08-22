package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"deeliai/internal/interfaces"
	"deeliai/internal/model"
)

// UserService 包含業務邏輯
type UserService struct {
	userRepo interfaces.UserRepository // 依賴介面，而非實作
}

func NewUserService(repo interfaces.UserRepository) *UserService {
	return &UserService{
		userRepo: repo,
	}
}

// Register 是一個業務邏輯方法
func (s *UserService) Register(ctx context.Context, email, password string) (*model.User, error) {
	if email == "" || password == "" {
		return nil, errors.New("email & password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	newUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *UserService) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

// Authenticate 驗證使用者帳號與密碼，成功則回傳使用者資訊
func (s *UserService) Authenticate(ctx context.Context, email, password string) (*model.User, error) {
	// 1. 根據 email 從資料庫查詢使用者
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		// 建議統一回傳 "Invalid email or password" 以避免暴露使用者是否存在
		return nil, errors.New("invalid email or password")
	}

	// 2. 使用 bcrypt 比對使用者輸入的密碼與資料庫中的雜湊密碼
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// 比對失敗，回傳認證失敗
		return nil, errors.New("invalid email or password")
	}

	// 3. 認證成功，回傳使用者資訊
	return user, nil
}
