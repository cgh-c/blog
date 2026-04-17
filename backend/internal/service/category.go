package service

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"blog-backend/internal/model"
	"blog-backend/internal/repository"
)

// CategoryService handles category business logic.
type CategoryService struct {
	repo *repository.CategoryRepo
}

// NewCategoryService creates a new CategoryService.
func NewCategoryService(repo *repository.CategoryRepo) *CategoryService {
	return &CategoryService{repo: repo}
}

// CreateCategoryInput holds data for creating a category.
type CreateCategoryInput struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"max=500"`
	SortOrder   int    `json:"sort_order"`
}

// UpdateCategoryInput holds data for updating a category.
type UpdateCategoryInput struct {
	Name        *string `json:"name" binding:"omitempty,min=1,max=100"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	SortOrder   *int    `json:"sort_order"`
}

// ListWithPublicCount returns categories with public article counts.
func (s *CategoryService) ListWithPublicCount() ([]repository.CategoryWithCount, error) {
	return s.repo.ListWithPublicCount()
}

// ListAll returns all categories.
func (s *CategoryService) ListAll() ([]model.Category, error) {
	return s.repo.List()
}

// Create creates a new category.
func (s *CategoryService) Create(input CreateCategoryInput) (*model.Category, error) {
	cat := model.Category{
		Name:        input.Name,
		Slug:        generateCategorySlug(input.Name),
		Description: input.Description,
		SortOrder:   input.SortOrder,
	}
	if err := s.repo.Create(&cat); err != nil {
		return nil, err
	}
	return &cat, nil
}

// Update updates an existing category.
func (s *CategoryService) Update(id uint, input UpdateCategoryInput) (*model.Category, error) {
	cat, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if input.Name != nil {
		cat.Name = *input.Name
		cat.Slug = generateCategorySlug(*input.Name)
	}
	if input.Description != nil {
		cat.Description = *input.Description
	}
	if input.SortOrder != nil {
		cat.SortOrder = *input.SortOrder
	}
	if err := s.repo.Update(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

// Delete removes a category.
func (s *CategoryService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func generateCategorySlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	reg := regexp.MustCompile(`[^a-z0-9-]+`)
	slug = reg.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = fmt.Sprintf("cat-%d", time.Now().UnixMilli())
	}
	return slug
}
