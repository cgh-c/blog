package repository

import (
	"blog-backend/internal/model"

	"gorm.io/gorm"
)

// ArticleRepo handles article database operations.
type ArticleRepo struct {
	db *gorm.DB
}

// NewArticleRepo creates a new ArticleRepo.
func NewArticleRepo(db *gorm.DB) *ArticleRepo {
	return &ArticleRepo{db: db}
}

// ArticleQuery holds optional query parameters for listing articles.
type ArticleQuery struct {
	CategoryID *uint
	TagID      *uint
	Visibility *model.Visibility
	Keyword    string
	Offset     int
	Limit      int
}

// List returns articles matching the query with preloaded relations.
func (r *ArticleRepo) List(q ArticleQuery) ([]model.Article, int64, error) {
	tx := r.db.Model(&model.Article{})

	if q.Visibility != nil {
		tx = tx.Where("visibility = ?", *q.Visibility)
	}
	if q.CategoryID != nil {
		tx = tx.Where("category_id = ?", *q.CategoryID)
	}
	if q.TagID != nil {
		tx = tx.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Where("article_tags.tag_id = ?", *q.TagID)
	}
	if q.Keyword != "" {
		like := "%" + q.Keyword + "%"
		tx = tx.Where("title LIKE ? OR summary LIKE ?", like, like)
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var articles []model.Article
	err := tx.Preload("Category").Preload("Tags").
		Order("created_at DESC").
		Offset(q.Offset).Limit(q.Limit).
		Find(&articles).Error

	return articles, total, err
}

// GetByID returns an article by its ID with preloaded relations.
func (r *ArticleRepo) GetByID(id uint) (*model.Article, error) {
	var article model.Article
	err := r.db.Preload("Category").Preload("Tags").
		First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// GetBySlug returns a article by slug with preloaded relations.
func (r *ArticleRepo) GetBySlug(slug string) (*model.Article, error) {
	var article model.Article
	err := r.db.Preload("Category").Preload("Tags").
		Where("slug = ?", slug).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// Create inserts a new article.
func (r *ArticleRepo) Create(article *model.Article) error {
	return r.db.Create(article).Error
}

// Update saves changes to an existing article.
func (r *ArticleRepo) Update(article *model.Article) error {
	return r.db.Save(article).Error
}

// UpdateTags replaces the tags association for an article.
func (r *ArticleRepo) UpdateTags(article *model.Article, tags []model.Tag) error {
	return r.db.Model(article).Association("Tags").Replace(tags)
}

// Delete removes an article by ID.
func (r *ArticleRepo) Delete(id uint) error {
	// Clear tag associations first
	r.db.Exec("DELETE FROM article_tags WHERE article_id = ?", id)
	return r.db.Delete(&model.Article{}, id).Error
}

// IncrementViewCount atomically increments the view count.
func (r *ArticleRepo) IncrementViewCount(id uint) error {
	return r.db.Model(&model.Article{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// ArchiveItem represents a month in the archive with article count.
type ArchiveItem struct {
	Year  int   `json:"year"`
	Month int   `json:"month"`
	Count int64 `json:"count"`
}

// GetArchives returns article counts grouped by year and month.
func (r *ArticleRepo) GetArchives(visibility *model.Visibility) ([]ArchiveItem, error) {
	tx := r.db.Model(&model.Article{})
	if visibility != nil {
		tx = tx.Where("visibility = ?", *visibility)
	}

	var archives []ArchiveItem
	err := tx.Select("strftime('%Y', created_at) as year, strftime('%m', created_at) as month, COUNT(*) as count").
		Group("year, month").
		Order("year DESC, month DESC").
		Scan(&archives).Error

	return archives, err
}
