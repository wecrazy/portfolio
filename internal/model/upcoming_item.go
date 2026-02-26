package model

import "gorm.io/gorm"

// UpcomingItem represents a future project or announcement in the "What's Next" section.
type UpcomingItem struct {
	gorm.Model
	Title       string `gorm:"column:title;size:200;not null" json:"title" form:"title"`
	Description string `gorm:"column:description;type:text" json:"description" form:"description"`
	Type        string `gorm:"column:type;size:30" json:"type" form:"type"`                         // "project" | "announcement"
	Status      string `gorm:"column:status;size:30;default:'planned'" json:"status" form:"status"` // "planned"|"in-progress"|"coming-soon"
	IconClass   string `gorm:"column:icon_class;size:100" json:"icon_class" form:"icon_class"`
	SortOrder   int    `gorm:"column:sort_order;default:0" json:"sort_order" form:"sort_order"`
	IsVisible   bool   `gorm:"column:is_visible;default:true" json:"is_visible"`
}

// TableName returns the database table name for UpcomingItem.
func (UpcomingItem) TableName() string { return tables.UpcomingItems }
