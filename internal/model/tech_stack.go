package model

import "gorm.io/gorm"

// TechStack represents a technology or tool in the portfolio owner's stack.
type TechStack struct {
	gorm.Model
	Name        string `gorm:"column:name;size:100;not null" json:"name" form:"name"`
	Category    string `gorm:"column:category;size:100;not null" json:"category" form:"category"`
	IconClass   string `gorm:"column:icon_class;size:100" json:"icon_class" form:"icon_class"`
	IconURL     string `gorm:"column:icon_url;size:500" json:"icon_url" form:"icon_url"`
	IconURLDark string `gorm:"column:icon_url_dark;size:500" json:"icon_url_dark" form:"icon_url_dark"`
	Desc        string `gorm:"column:description;size:500" json:"description" form:"description"`
	URL         string `gorm:"column:url;size:500" json:"url" form:"url"`
	SortOrder   int    `gorm:"column:sort_order;default:0" json:"sort_order" form:"sort_order"`
}

// TableName returns the database table name for TechStack.
func (TechStack) TableName() string { return tables.TechStacks }
