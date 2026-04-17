package repository

import (
	"blog-backend/internal/model"

	"gorm.io/gorm"
)

// TagRepo handles tag database operations.
type TagRepo struct {
	db *gorm.DB
}

// NewTagRepo creates a new TagRepo.
func NewTagRepo(db *gorm.DB) *TagRepo {
	return &TagRepo{db: db}
}

// List returns all tags.
func (r *TagRepo) List() ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.Order("name ASC").Find(&tags).Error
	return tags, err
}

// GetByID returns a tag by ID.
func (r *TagRepo) GetByID(id uint) (*model.Tag, error) {
	var tag model.Tag
	if err := r.db.First(&tag, id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetByIDs returns tags matching the given IDs.
func (r *TagRepo) GetByIDs(ids []uint) ([]model.Tag, error) {
	var tags []model.Tag
	if len(ids) == 0 {
		return tags, nil
	}
	err := r.db.Where("id IN ?", ids).Find(&tags).Error
	return tags, err
}

// Create inserts a new tag.
func (r *TagRepo) Create(tag *model.Tag) error {
	return r.db.Create(tag).Error
}

// Delete removes a tag by ID.
func (r *TagRepo) Delete(id uint) error {
	// Clean up article_tags associations
	r.db.Exec("DELETE FROM article_tags WHERE tag_id = ?", id)
	return r.db.Delete(&model.Tag{}, id).Error
}
