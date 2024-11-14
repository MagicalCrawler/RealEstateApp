package models

import (
	"github.com/MagicalCrawler/RealEstateApp/types"
	"gorm.io/gorm"
)

type Post struct {
	UniqueCode string              `gorm:"not null"`
	Website    types.WebsiteSource `gorm:"not nll;type:string"`
	WatchedNum uint
	gorm.Model
}
