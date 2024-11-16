package models

import (
	"github.com/MagicalCrawler/RealEstateApp/types"
	"gorm.io/gorm"
)

type Post struct {
	UniqueCode string              `gorm:"not null;unique"`
	Website    types.WebsiteSource `gorm:"not null;type:string"`
	WatchedNum uint
	gorm.Model
}
