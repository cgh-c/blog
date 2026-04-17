package model

import "time"

// Visibility defines the access level of an article.
type Visibility string

const (
	VisibilityPublic  Visibility = "public"
	VisibilityPrivate Visibility = "private"
	VisibilityDraft   Visibility = "draft"
)

// Article represents a blog post.
type Article struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Slug        string     `gorm:"uniqueIndex;size:200;not null" json:"slug"`
	Title       string     `gorm:"size:200;not null" json:"title"`
	Content     string     `gorm:"type:text;not null" json:"content"`
	Summary     string     `gorm:"size:500" json:"summary"`
	CoverImage  string     `gorm:"size:500" json:"cover_image"`
	CategoryID  *uint      `gorm:"index" json:"category_id"`
	Category    *Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags        []Tag      `gorm:"many2many:article_tags" json:"tags,omitempty"`
	Visibility  Visibility `gorm:"size:20;default:draft;not null;index" json:"visibility"`
	ViewCount   uint       `gorm:"default:0" json:"view_count"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at"`
}
