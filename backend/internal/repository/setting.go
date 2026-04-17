package repository

import (
	"blog-backend/internal/model"

	"gorm.io/gorm"
)

// SettingRepo handles setting database operations.
type SettingRepo struct {
	db *gorm.DB
}

// NewSettingRepo creates a new SettingRepo.
func NewSettingRepo(db *gorm.DB) *SettingRepo {
	return &SettingRepo{db: db}
}

// ListPublic returns all settings marked as public.
func (r *SettingRepo) ListPublic() ([]model.Setting, error) {
	var settings []model.Setting
	err := r.db.Where("is_public = ?", true).Find(&settings).Error
	return settings, err
}

// ListAll returns all settings.
func (r *SettingRepo) ListAll() ([]model.Setting, error) {
	var settings []model.Setting
	err := r.db.Find(&settings).Error
	return settings, err
}

// GetByKey returns a setting by its key.
func (r *SettingRepo) GetByKey(key string) (*model.Setting, error) {
	var s model.Setting
	if err := r.db.Where("`key` = ?", key).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

// Upsert creates or updates a setting by key.
func (r *SettingRepo) Upsert(key, value string) error {
	var s model.Setting
	result := r.db.Where("`key` = ?", key).First(&s)
	if result.Error == gorm.ErrRecordNotFound {
		s = model.Setting{Key: key, Value: value}
		return r.db.Create(&s).Error
	}
	if result.Error != nil {
		return result.Error
	}
	s.Value = value
	return r.db.Save(&s).Error
}
