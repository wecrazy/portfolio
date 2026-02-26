package model

import "gorm.io/gorm"

// ContactMessage represents a message submitted through the public contact form.
type ContactMessage struct {
	gorm.Model
	Name    string `gorm:"column:name;size:100;not null" json:"name" form:"name"`
	Email   string `gorm:"column:email;size:100;not null" json:"email" form:"email"`
	Subject string `gorm:"column:subject;size:200" json:"subject" form:"subject"`
	Message string `gorm:"column:message;type:text;not null" json:"message" form:"message"`
	IsRead  bool   `gorm:"column:is_read;default:false" json:"is_read"`
}

// TableName returns the database table name for ContactMessage.
func (ContactMessage) TableName() string { return tables.ContactMessages }
