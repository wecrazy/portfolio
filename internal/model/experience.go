package model

import (
	"time"

	"gorm.io/gorm"
)

// Experience represents a work, education, or certification entry on the timeline.
type Experience struct {
	gorm.Model
	Type        string        `gorm:"column:type;size:20;not null" json:"type" form:"type"`
	Title       string        `gorm:"column:title;size:200;not null" json:"title" form:"title"`
	Org         string        `gorm:"column:org;size:200" json:"org" form:"org"`
	Location    string        `gorm:"column:location;size:100" json:"location" form:"location"`
	StartDate   time.Time     `gorm:"column:start_date;not null" json:"start_date" form:"start_date"`
	EndDate     *time.Time    `gorm:"column:end_date" json:"end_date" form:"end_date"`
	IsCurrent   bool          `gorm:"column:is_current;default:false" json:"is_current" form:"is_current"`
	Description string        `gorm:"column:description;type:text" json:"description" form:"description"`
	SortOrder   int           `gorm:"column:sort_order;default:0" json:"sort_order" form:"sort_order"`
	ImageURL    string        `gorm:"column:image_url;size:500" json:"image_url" form:"image_url"`
	ImageID     *uint         `gorm:"column:image_id" json:"image_id"`
	Image       *UploadedFile `gorm:"foreignKey:ImageID" json:"image,omitempty"`
}

// TableName returns the database table name for Experience.
func (Experience) TableName() string { return tables.Experiences }
