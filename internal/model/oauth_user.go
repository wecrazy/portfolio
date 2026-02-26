package model

import "gorm.io/gorm"

// OAuthUser represents a visitor who authenticated via Google or GitHub to leave comments.
type OAuthUser struct {
	gorm.Model
	Provider    string    `gorm:"column:provider;size:20;not null" json:"provider"`
	ProviderID  string    `gorm:"column:provider_id;size:255;not null" json:"provider_id"`
	Email       string    `gorm:"column:email;size:100" json:"email"`
	DisplayName string    `gorm:"column:display_name;size:100" json:"display_name"`
	AvatarURL   string    `gorm:"column:avatar_url;size:500" json:"avatar_url"`
	Comments    []Comment `gorm:"foreignKey:OAuthUserID" json:"comments,omitempty"`
}

// TableName returns the database table name for OAuthUser.
func (OAuthUser) TableName() string { return tables.OAuthUsers }
