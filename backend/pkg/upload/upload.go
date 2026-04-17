package upload

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

// Config holds upload configuration.
type Config struct {
	Dir            string   // directory to save uploaded files
	MaxSize        int64    // max file size in bytes
	AllowedTypes   []string // allowed MIME types
	AllowedExtensions []string // allowed file extensions (e.g. ".jpg", ".png")
}

// DefaultConfig returns a sensible default upload config.
func DefaultConfig() Config {
	return Config{
		Dir:            "./uploads",
		MaxSize:        5 * 1024 * 1024, // 5MB
		AllowedTypes:   []string{"image/jpeg", "image/png", "image/webp", "image/gif"},
		AllowedExtensions: []string{".jpg", ".jpeg", ".png", ".webp", ".gif"},
	}
}

// Uploader handles file upload with validation.
type Uploader struct {
	config Config
}

// NewUploader creates a new Uploader with the given config.
func NewUploader(config Config) *Uploader {
	return &Uploader{config: config}
}

// Upload validates and saves an uploaded file.
// Returns the relative URL path to the saved file.
func (u *Uploader) Upload(c *gin.Context, fieldName string) (string, error) {
	file, header, err := c.Request.FormFile(fieldName)
	if err != nil {
		return "", fmt.Errorf("failed to read uploaded file: %w", err)
	}
	defer file.Close()

	if err := u.validate(header); err != nil {
		return "", err
	}

	// Generate random filename to prevent path traversal and collisions
	filename, err := u.generateFilename(header.Filename)
	if err != nil {
		return "", fmt.Errorf("failed to generate filename: %w", err)
	}

	// Organize by year/month subdirectories
	subDir := time.Now().Format("2006/01")
	saveDir := filepath.Join(u.config.Dir, subDir)
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	savePath := filepath.Join(saveDir, filename)
	if err := c.SaveUploadedFile(header, savePath); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Return the URL path (relative to upload dir)
	urlPath := fmt.Sprintf("/uploads/%s/%s", subDir, filename)
	return urlPath, nil
}

func (u *Uploader) validate(header *multipart.FileHeader) error {
	// Check file size
	if header.Size > u.config.MaxSize {
		return fmt.Errorf("file too large: max %d bytes allowed", u.config.MaxSize)
	}

	// Check extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	extAllowed := false
	for _, allowed := range u.config.AllowedExtensions {
		if ext == allowed {
			extAllowed = true
			break
		}
	}
	if !extAllowed {
		return fmt.Errorf("file type not allowed: %s", ext)
	}

	// Check MIME type
	contentType := header.Header.Get("Content-Type")
	typeAllowed := false
	for _, allowed := range u.config.AllowedTypes {
		if contentType == allowed {
			typeAllowed = true
			break
		}
	}
	if !typeAllowed {
		return fmt.Errorf("MIME type not allowed: %s", contentType)
	}

	return nil
}

func (u *Uploader) generateFilename(original string) (string, error) {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	ext := strings.ToLower(filepath.Ext(original))
	return hex.EncodeToString(randomBytes) + ext, nil
}
