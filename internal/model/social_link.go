package model

import "gorm.io/gorm"

// SocialLink represents a linked social media or external account.
type SocialLink struct {
	gorm.Model
	Platform  string `gorm:"column:platform;size:50;not null" json:"platform" form:"platform"`
	URL       string `gorm:"column:url;size:500;not null" json:"url" form:"url"`
	IconClass string `gorm:"column:icon_class;size:100" json:"icon_class" form:"icon_class"`
	Label     string `gorm:"column:label;size:100" json:"label" form:"label"`
	SortOrder int    `gorm:"column:sort_order;default:0" json:"sort_order" form:"sort_order"`
}

// TableName returns the database table name for SocialLink.
func (SocialLink) TableName() string { return tables.SocialLinks }
