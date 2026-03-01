package model

import "gorm.io/gorm"

// Owner represents the portfolio owner's profile and about information.
type Owner struct {
	gorm.Model
	FullName       string        `gorm:"column:full_name;size:100;not null" json:"full_name" form:"full_name"`
	Title          string        `gorm:"column:title;size:150" json:"title" form:"title"`
	Tagline        string        `gorm:"column:tagline;size:300" json:"tagline" form:"tagline"`
	Bio            string        `gorm:"column:bio;type:text" json:"bio" form:"bio"`
	Email          string        `gorm:"column:email;size:100" json:"email" form:"email"`
	Phone          string        `gorm:"column:phone;size:30" json:"phone" form:"phone"`
	Location       string        `gorm:"column:location;size:100" json:"location" form:"location"`
	ProfileImageID *uint         `gorm:"column:profile_image_id" json:"profile_image_id"`
	ResumeFileID   *uint         `gorm:"column:resume_file_id" json:"resume_file_id"`
	ProfileImage   *UploadedFile `gorm:"foreignKey:ProfileImageID" json:"profile_image,omitempty"`
	ResumeFile     *UploadedFile `gorm:"foreignKey:ResumeFileID" json:"resume_file,omitempty"`
}

// TableName returns the database table name for Owner.
func (Owner) TableName() string { return tables.Owners }
