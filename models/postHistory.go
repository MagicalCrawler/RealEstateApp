package models

import (
	"github.com/MagicalCrawler/RealEstateApp/types"
	"gorm.io/gorm"
)

type PostHistory struct {
	gorm.Model
	PostID       uint
	Post         Post
	PostURL      string `gorm:"type:text"`
	Price        int64
	City         string `gorm:"type:varchar(63)"`
	Neighbor     string `gorm:"type:varchar(63)"`
	Area         int
	BedroomNum   int
	BuyMode      types.BuyMode  `gorm:"type:string"`
	Building     types.Building `gorm:"type:string"`
	Age          uint8
	FloorsNum    uint8
	HasWareHouse bool
	HsaParking   bool
	HasElevator  bool
	ImageURL     string `gorm:"type:text"`
	Description  string `gorm:"type:text"`
}
