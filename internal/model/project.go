package model

import "gorm.io/gorm"

// Project represents a portfolio project entry.
type Project struct {
	gorm.Model
	Title        string `gorm:"column:title;size:200;not null" json:"title" form:"title"`
	Slug         string `gorm:"column:slug;size:200;uniqueIndex" json:"slug"`
	Description  string `gorm:"column:description;type:text" json:"description" form:"description"`
	LongDesc     string `gorm:"column:long_desc;type:text" json:"long_desc" form:"long_desc"`
	ThumbnailURL string `gorm:"column:thumbnail_url;size:500" json:"thumbnail_url" form:"thumbnail_url"`
	LiveURL      string `gorm:"column:live_url;size:500" json:"live_url" form:"live_url"`
	RepoURL      string `gorm:"column:repo_url;size:500" json:"repo_url" form:"repo_url"`
	Tags         string `gorm:"column:tags;size:500" json:"tags" form:"tags"`
	Status       string `gorm:"column:status;size:20;default:'draft'" json:"status" form:"status"`
	SortOrder    int    `gorm:"column:sort_order;default:0" json:"sort_order" form:"sort_order"`
	Featured     bool   `gorm:"column:featured;default:false" json:"featured" form:"featured"`
}

// TableName returns the database table name for Project.
func (Project) TableName() string { return tables.Projects }
