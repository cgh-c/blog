package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT payload.
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT token operations.
type JWTManager struct {
	secret        []byte
	expiry        time.Duration
	refreshWindow time.Duration // if remaining < refreshWindow, auto-refresh
}

// NewJWTManager creates a new JWT manager.
// secret: signing key; expiry: token lifetime; refreshWindow: auto-refresh threshold.
func NewJWTManager(secret string, expiry time.Duration, refreshWindow time.Duration) *JWTManager {
	return &JWTManager{
		secret:        []byte(secret),
		expiry:        expiry,
		refreshWindow: refreshWindow,
	}
}

// GenerateToken creates a new signed JWT for the given user.
func (m *JWTManager) GenerateToken(userID uint, username string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ParseToken validates and parses a JWT string, returning the claims.
func (m *JWTManager) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ShouldRefresh returns true if the token is valid but close to expiry.
func (m *JWTManager) ShouldRefresh(claims *Claims) bool {
	if claims.ExpiresAt == nil {
		return false
	}
	remaining := time.Until(claims.ExpiresAt.Time)
	return remaining > 0 && remaining < m.refreshWindow
}
