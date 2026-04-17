package service

import (
	"errors"
	"time"

	"blog-backend/internal/model"
	"blog-backend/pkg/auth"

	"gorm.io/gorm"
)

// AuthService handles authentication business logic.
type AuthService struct {
	db *gorm.DB
}

// NewAuthService creates a new AuthService.
func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

// ErrInvalidCredentials is returned when login credentials are wrong.
// The message is deliberately generic to avoid leaking whether the username exists.
var ErrInvalidCredentials = errors.New("invalid username or password")

// Login validates credentials and returns the user on success.
func (s *AuthService) Login(username, password string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		// Don't distinguish between "user not found" and "wrong password"
		return nil, ErrInvalidCredentials
	}

	if err := auth.CheckPassword(password, user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Update last login time
	now := time.Now()
	s.db.Model(&user).Update("last_login_at", now)

	return &user, nil
}
