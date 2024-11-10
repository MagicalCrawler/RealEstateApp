package models

import "gorm.io/gorm"

type Bookmark struct {
	gorm.Model
	UserID uint `gorm:"not null;foreignKey:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PostID uint `gorm:"foreignKey:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
