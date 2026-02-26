package model

import "gorm.io/gorm"

// UploadedFile tracks files uploaded through the admin dashboard.
type UploadedFile struct {
	gorm.Model
	OriginalName string `gorm:"column:original_name;size:255;not null" json:"original_name"`
	StoredName   string `gorm:"column:stored_name;size:255;not null" json:"stored_name"`
	FilePath     string `gorm:"column:file_path;size:500;not null" json:"file_path"`
	MimeType     string `gorm:"column:mime_type;size:100" json:"mime_type"`
	FileSize     int64  `gorm:"column:file_size" json:"file_size"`
	Category     string `gorm:"column:category;size:50" json:"category"`
}

// TableName returns the database table name for UploadedFile.
func (UploadedFile) TableName() string { return tables.UploadedFiles }
