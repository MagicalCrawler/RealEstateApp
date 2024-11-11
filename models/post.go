package models

import (
	"gorm.io/gorm"
)

type Post struct {
	Title string `gorm:"type:text"`
	gorm.Model
}
