package models

import (
	"github.com/MagicalCrawler/RealEstateApp/types"
	"gorm.io/gorm"
)

type Post struct {
	UniqueCode string              `gorm:"not null;unique"`      // each ads has a unique code in divar
	Website    types.WebsiteSource `gorm:"not null;type:string"` // for search between some sources
	WatchedNum uint
	gorm.Model
}
