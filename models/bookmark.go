package models

import "gorm.io/gorm"

type Bookmark struct {
	User   User
	Post   Post
	UserID uint `gorm:"not null;foreignKey:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PostID uint `gorm:"foreignKey:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	gorm.Model
}
