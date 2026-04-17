package service

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"blog-backend/internal/model"
	"blog-backend/internal/repository"
)

// ArticleService handles article business logic.
type ArticleService struct {
	articleRepo *repository.ArticleRepo
	tagRepo     *repository.TagRepo
}

// NewArticleService creates a new ArticleService.
func NewArticleService(articleRepo *repository.ArticleRepo, tagRepo *repository.TagRepo) *ArticleService {
	return &ArticleService{
		articleRepo: articleRepo,
		tagRepo:     tagRepo,
	}
}

// CreateArticleInput holds the data needed to create an article.
type CreateArticleInput struct {
	Title      string           `json:"title" binding:"required,min=1,max=200"`
	Content    string           `json:"content" binding:"required"`
	Summary    string           `json:"summary" binding:"max=500"`
	CoverImage string           `json:"cover_image"`
	CategoryID *uint            `json:"category_id"`
	TagIDs     []uint           `json:"tag_ids"`
	Visibility model.Visibility `json:"visibility" binding:"required,oneof=public private draft"`
}

// UpdateArticleInput holds the data needed to update an article.
type UpdateArticleInput struct {
	Title      *string           `json:"title" binding:"omitempty,min=1,max=200"`
	Content    *string           `json:"content"`
	Summary    *string           `json:"summary" binding:"omitempty,max=500"`
	CoverImage *string           `json:"cover_image"`
	CategoryID *uint             `json:"category_id"`
	TagIDs     *[]uint           `json:"tag_ids"`
	Visibility *model.Visibility `json:"visibility" binding:"omitempty,oneof=public private draft"`
}

// ListPublic returns only public articles with pagination.
func (s *ArticleService) ListPublic(offset, limit int, categoryID, tagID *uint) ([]model.Article, int64, error) {
	vis := model.VisibilityPublic
	return s.articleRepo.List(repository.ArticleQuery{
		Visibility: &vis,
		CategoryID: categoryID,
		TagID:      tagID,
		Offset:     offset,
		Limit:      limit,
	})
}

// ListAll returns all articles (for admin) with pagination.
func (s *ArticleService) ListAll(offset, limit int, categoryID, tagID *uint, visibility *model.Visibility) ([]model.Article, int64, error) {
	return s.articleRepo.List(repository.ArticleQuery{
		Visibility: visibility,
		CategoryID: categoryID,
		TagID:      tagID,
		Offset:     offset,
		Limit:      limit,
	})
}

// GetPublicBySlug returns a public article by slug and increments view count.
func (s *ArticleService) GetPublicBySlug(slug string) (*model.Article, error) {
	article, err := s.articleRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	if article.Visibility != model.VisibilityPublic {
		return nil, fmt.Errorf("article not found")
	}
	s.articleRepo.IncrementViewCount(article.ID)
	return article, nil
}

// GetByID returns an article by ID (admin, all visibilities).
func (s *ArticleService) GetByID(id uint) (*model.Article, error) {
	return s.articleRepo.GetByID(id)
}

// Create creates a new article.
func (s *ArticleService) Create(input CreateArticleInput) (*model.Article, error) {
	article := model.Article{
		Slug:       generateSlug(input.Title),
		Title:      input.Title,
		Content:    input.Content,
		Summary:    input.Summary,
		CoverImage: input.CoverImage,
		CategoryID: input.CategoryID,
		Visibility: input.Visibility,
	}

	if input.Visibility == model.VisibilityPublic {
		now := time.Now()
		article.PublishedAt = &now
	}

	if err := s.articleRepo.Create(&article); err != nil {
		return nil, err
	}

	// Set tags
	if len(input.TagIDs) > 0 {
		tags, err := s.tagRepo.GetByIDs(input.TagIDs)
		if err != nil {
			return nil, err
		}
		if err := s.articleRepo.UpdateTags(&article, tags); err != nil {
			return nil, err
		}
	}

	return s.articleRepo.GetByID(article.ID)
}

// Update updates an existing article.
func (s *ArticleService) Update(id uint, input UpdateArticleInput) (*model.Article, error) {
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		article.Title = *input.Title
		article.Slug = generateSlug(*input.Title)
	}
	if input.Content != nil {
		article.Content = *input.Content
	}
	if input.Summary != nil {
		article.Summary = *input.Summary
	}
	if input.CoverImage != nil {
		article.CoverImage = *input.CoverImage
	}
	if input.CategoryID != nil {
		article.CategoryID = input.CategoryID
	}
	if input.Visibility != nil {
		// If transitioning to public and not yet published, set PublishedAt
		if *input.Visibility == model.VisibilityPublic && article.PublishedAt == nil {
			now := time.Now()
			article.PublishedAt = &now
		}
		article.Visibility = *input.Visibility
	}

	if err := s.articleRepo.Update(article); err != nil {
		return nil, err
	}

	// Update tags if provided
	if input.TagIDs != nil {
		tags, err := s.tagRepo.GetByIDs(*input.TagIDs)
		if err != nil {
			return nil, err
		}
		if err := s.articleRepo.UpdateTags(article, tags); err != nil {
			return nil, err
		}
	}

	return s.articleRepo.GetByID(id)
}

// Delete removes an article.
func (s *ArticleService) Delete(id uint) error {
	return s.articleRepo.Delete(id)
}

// GetPublicArchives returns archives for public articles.
func (s *ArticleService) GetPublicArchives() ([]repository.ArchiveItem, error) {
	vis := model.VisibilityPublic
	return s.articleRepo.GetArchives(&vis)
}

// generateSlug creates a URL-friendly slug from a title.
func generateSlug(title string) string {
	slug := strings.ToLower(strings.TrimSpace(title))

	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9\p{Han}-]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	// If slug is empty (e.g. all Chinese), use a timestamp-based slug
	if slug == "" || allChinese(title) {
		slug = fmt.Sprintf("post-%d", time.Now().UnixMilli())
	}

	// Limit length
	if len(slug) > 200 {
		slug = slug[:200]
	}

	return slug
}

func allChinese(s string) bool {
	for _, r := range s {
		if !unicode.Is(unicode.Han, r) && !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}
