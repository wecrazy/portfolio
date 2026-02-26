package model

import "gorm.io/gorm"

// Comment represents a comment left by an OAuth-authenticated visitor or an owner reply.
type Comment struct {
	gorm.Model
	OAuthUserID  uint      `gorm:"column:oauth_user_id;not null;index" json:"oauth_user_id"`
	ParentID     *uint     `gorm:"column:parent_id;index" json:"parent_id"`
	Body         string    `gorm:"column:body;type:text;not null" json:"body" form:"body"`
	IsOwnerReply bool      `gorm:"column:is_owner_reply;default:false" json:"is_owner_reply"`
	IsApproved   bool      `gorm:"column:is_approved;default:true" json:"is_approved"`
	OAuthUser    OAuthUser `gorm:"foreignKey:OAuthUserID" json:"oauth_user"`
	Replies      []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

// TableName returns the database table name for Comment.
func (Comment) TableName() string { return tables.Comments }
