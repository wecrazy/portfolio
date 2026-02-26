// Package model defines all GORM database models.
package model

import (
	"time"

	"gorm.io/gorm"
)

// Admin represents an administrator who manages the portfolio.
type Admin struct {
	gorm.Model
	Username     string     `gorm:"column:username;size:50;uniqueIndex;not null" json:"username"`
	Email        string     `gorm:"column:email;size:100;uniqueIndex;not null" json:"email"`
	PasswordHash string     `gorm:"column:password_hash;size:255;not null" json:"-"`
	SessionToken string     `gorm:"column:session_token;size:255" json:"-"`
	LastLoginAt  *time.Time `gorm:"column:last_login_at" json:"last_login_at"`
}

// TableName returns the database table name for Admin.
func (Admin) TableName() string { return tables.Admins }
