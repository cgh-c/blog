package main

import (
	"log"
	"os"

	"blog-backend/internal/config"
	"blog-backend/internal/database"
	"blog-backend/internal/router"
)

func main() {
	// Load .env file if present (development convenience)
	loadEnvFile(".env")

	cfg := config.Load()

	// Initialize database and run migrations
	db := database.Init(cfg.DBPath, cfg.IsProduction())

	// Seed initial admin user and default settings
	database.SeedAdmin(db, cfg.AdminUsername, cfg.AdminPassword)
	database.SeedSettings(db)

	// Create upload directory
	os.MkdirAll(cfg.UploadDir, 0755)

	// Setup router with all routes
	r := router.Setup(cfg, db)

	log.Printf("Blog backend starting on :%s (env: %s)", cfg.AppPort, cfg.AppEnv)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// loadEnvFile reads a .env file and sets environment variables.
// Silently ignored if file doesn't exist.
func loadEnvFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	for _, line := range splitLines(string(data)) {
		line = trimString(line)
		if line == "" || line[0] == '#' {
			continue
		}

		eqIdx := -1
		for i := 0; i < len(line); i++ {
			if line[i] == '=' {
				eqIdx = i
				break
			}
		}
		if eqIdx < 0 {
			continue
		}

		key := trimString(line[:eqIdx])
		value := trimString(line[eqIdx+1:])

		// Only set if not already set (env vars take precedence)
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func trimString(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}
