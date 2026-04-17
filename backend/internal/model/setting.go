package model

import "time"

// Setting represents a site configuration key-value pair.
type Setting struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Key         string    `gorm:"size:100;uniqueIndex;not null" json:"key"`
	Value       string    `gorm:"type:text" json:"value"`
	Type        string    `gorm:"size:20;default:string" json:"type"` // string, text, json
	IsPublic    bool      `gorm:"default:true;index" json:"is_public"`
	Description string    `gorm:"size:200" json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}
