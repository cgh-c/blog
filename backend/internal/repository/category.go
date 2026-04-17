package repository

import (
	"blog-backend/internal/model"

	"gorm.io/gorm"
)

// CategoryRepo handles category database operations.
type CategoryRepo struct {
	db *gorm.DB
}

// NewCategoryRepo creates a new CategoryRepo.
func NewCategoryRepo(db *gorm.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

// CategoryWithCount includes the number of public articles in a category.
type CategoryWithCount struct {
	model.Category
	ArticleCount int64 `json:"article_count"`
}

// List returns all categories ordered by sort_order.
func (r *CategoryRepo) List() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Order("sort_order ASC, id ASC").Find(&categories).Error
	return categories, err
}

// ListWithPublicCount returns categories with their public article counts.
func (r *CategoryRepo) ListWithPublicCount() ([]CategoryWithCount, error) {
	var results []CategoryWithCount
	err := r.db.Model(&model.Category{}).
		Select("categories.*, (SELECT COUNT(*) FROM articles WHERE articles.category_id = categories.id AND articles.visibility = 'public') as article_count").
		Order("sort_order ASC, id ASC").
		Scan(&results).Error
	return results, err
}

// GetByID returns a category by ID.
func (r *CategoryRepo) GetByID(id uint) (*model.Category, error) {
	var cat model.Category
	if err := r.db.First(&cat, id).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

// GetBySlug returns a category by slug.
func (r *CategoryRepo) GetBySlug(slug string) (*model.Category, error) {
	var cat model.Category
	if err := r.db.Where("slug = ?", slug).First(&cat).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

// Create inserts a new category.
func (r *CategoryRepo) Create(cat *model.Category) error {
	return r.db.Create(cat).Error
}

// Update saves changes to an existing category.
func (r *CategoryRepo) Update(cat *model.Category) error {
	return r.db.Save(cat).Error
}

// Delete removes a category by ID.
func (r *CategoryRepo) Delete(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}
