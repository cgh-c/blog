package database

import (
	"log"
	"os"
	"path/filepath"

	"blog-backend/internal/model"
	"blog-backend/pkg/auth"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Init opens the SQLite database, runs auto-migration, and returns the DB instance.
func Init(dbPath string, isProduction bool) *gorm.DB {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	logLevel := logger.Info
	if isProduction {
		logLevel = logger.Warn
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Enable WAL mode for better concurrent read performance
	db.Exec("PRAGMA journal_mode=WAL")
	db.Exec("PRAGMA foreign_keys=ON")

	// Auto-migrate all models
	if err := db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Tag{},
		&model.Article{},
		&model.Setting{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

// SeedAdmin creates the initial admin user if no users exist.
func SeedAdmin(db *gorm.DB, username, password string) {
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count > 0 {
		return
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		log.Fatalf("Failed to hash admin password: %v", err)
	}

	admin := model.User{
		Username:     username,
		PasswordHash: hash,
	}
	if err := db.Create(&admin).Error; err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	log.Printf("Admin user '%s' created successfully", username)
}

// SeedSettings creates default site settings if none exist.
func SeedSettings(db *gorm.DB) {
	var count int64
	db.Model(&model.Setting{}).Count(&count)
	if count > 0 {
		return
	}

	defaults := []model.Setting{
		{Key: "site_name", Value: "My Blog", Type: "string", IsPublic: true, Description: "Blog name displayed in navbar"},
		{Key: "site_description", Value: "A personal blog", Type: "string", IsPublic: true, Description: "Blog description for SEO and footer"},
		{Key: "about_content", Value: "# About Me\n\nWelcome to my blog!", Type: "text", IsPublic: true, Description: "About page content in Markdown"},
		{Key: "avatar_url", Value: "", Type: "string", IsPublic: true, Description: "Avatar URL for about page"},
		{Key: "social_links", Value: "[]", Type: "json", IsPublic: true, Description: "Social media links as JSON array"},
		{Key: "icp_number", Value: "", Type: "string", IsPublic: true, Description: "ICP filing number for footer"},
		{Key: "admin_email", Value: "", Type: "string", IsPublic: false, Description: "Admin email (internal use only)"},
	}

	db.Create(&defaults)
	log.Println("Default settings seeded")
}
