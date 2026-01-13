package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Slug        string         `gorm:"unique;index" json:"slug"`
	ContentMD   string         `gorm:"type:text" json:"content_md"`
	ContentHTML string         `gorm:"type:text" json:"content_html"`
	Published   bool           `gorm:"default:false" json:"published"`
	AuthorID    uint           `gorm:"not null" json:"author_id"`
	Author      User           `json:"author"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
