package model

import "gorm.io/gorm"

// Skill represents a single skill with category and proficiency level.
type Skill struct {
	gorm.Model
	Name        string `gorm:"column:name;size:100;not null" json:"name" form:"name"`
	Category    string `gorm:"column:category;size:50" json:"category" form:"category"`
	IconClass   string `gorm:"column:icon_class;size:100" json:"icon_class" form:"icon_class"`
	Proficiency int    `gorm:"column:proficiency;default:0" json:"proficiency" form:"proficiency"`
	SortOrder   int    `gorm:"column:sort_order;default:0" json:"sort_order" form:"sort_order"`
}

// TableName returns the database table name for Skill.
func (Skill) TableName() string { return tables.Skills }
