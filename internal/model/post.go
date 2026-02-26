package model

import (
	"time"

	"gorm.io/gorm"
)

// Post represents a blog/feed post entry.
type Post struct {
	gorm.Model
	Title           string        `gorm:"column:title;size:300;not null" json:"title" form:"title"`
	Slug            string        `gorm:"column:slug;size:300;uniqueIndex" json:"slug"`
	Excerpt         string        `gorm:"column:excerpt;type:text" json:"excerpt" form:"excerpt"`
	Content         string        `gorm:"column:content;type:text" json:"content" form:"content"`
	ThumbnailFileID *uint         `gorm:"column:thumbnail_file_id" json:"thumbnail_file_id" form:"thumbnail_file_id"`
	ThumbnailFile   *UploadedFile `gorm:"foreignKey:ThumbnailFileID" json:"thumbnail_file,omitempty"`
	Tags            string        `gorm:"column:tags;size:500" json:"tags" form:"tags"`
	Status          string        `gorm:"column:status;size:20;default:'draft'" json:"status" form:"status"`
	SortOrder       int           `gorm:"column:sort_order;default:0" json:"sort_order" form:"sort_order"`
	PublishedAt     *time.Time    `gorm:"column:published_at" json:"published_at"`
}

// TableName returns the database table name for Post.
func (Post) TableName() string { return tables.Posts }
