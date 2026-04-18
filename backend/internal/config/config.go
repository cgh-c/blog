package config

import (
	"os"
	"time"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	AppEnv  string
	AppPort string

	DBPath string

	JWTSecret string
	JWTExpiry time.Duration

	AdminUsername string
	AdminPassword string

	UploadDir     string
	UploadMaxSize int64

	CORSOrigins    []string
	SecureCookie   bool
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	cfg := &Config{
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8080"),

		DBPath: getEnv("DB_PATH", "./data/blog.db"),

		JWTSecret: getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		JWTExpiry: parseDuration(getEnv("JWT_EXPIRY", "24h")),

		AdminUsername: getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "admin123"),

		UploadDir:     getEnv("UPLOAD_DIR", "./uploads"),
		UploadMaxSize: 5 * 1024 * 1024, // 5MB

		CORSOrigins:  splitComma(getEnv("CORS_ORIGINS", "http://localhost:5173,http://localhost:5174")),
		SecureCookie: getEnv("SECURE_COOKIE", "false") == "true",
	}

	return cfg
}

// IsProduction returns true if running in production environment.
func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 24 * time.Hour
	}
	return d
}

func splitComma(s string) []string {
	if s == "" {
		return nil
	}
	var result []string
	for _, part := range split(s, ',') {
		trimmed := trim(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func split(s string, sep byte) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}

func trim(s string) string {
	start, end := 0, len(s)
	for start < end && s[start] == ' ' {
		start++
	}
	for end > start && s[end-1] == ' ' {
		end--
	}
	return s[start:end]
}
