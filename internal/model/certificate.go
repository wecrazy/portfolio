package model

import (
	"time"

	"gorm.io/gorm"
)

// Certificate represents a professional certificate or achievement.
type Certificate struct {
	gorm.Model
	Title       string        `gorm:"column:title;size:200;not null" json:"title" form:"title"`
	Description string        `gorm:"column:description;type:text" json:"description" form:"description"`
	Issuer      string        `gorm:"column:issuer;size:100" json:"issuer" form:"issuer"`
	IssueDate   time.Time     `gorm:"column:issue_date" json:"issue_date" form:"issue_date"`
	CertURL     string        `gorm:"column:cert_url;size:500" json:"cert_url" form:"cert_url"`
	FileID      *uint         `gorm:"column:file_id" json:"file_id"`
	File        *UploadedFile `gorm:"foreignKey:FileID" json:"file,omitempty"`
	IsVisible   bool          `gorm:"column:is_visible;default:true" json:"is_visible" form:"is_visible"`
	SortOrder   int           `gorm:"column:sort_order;default:0" json:"sort_order" form:"sort_order"`
}

// TableName overrides the default GORM table name for certificates.  This
// is necessary because the table name is configurable via the tables struct.
func (Certificate) TableName() string { return tables.Certificates }
