package models

import (
	"github.com/MagicalCrawler/RealEstateApp/types"
	"gorm.io/gorm"
)

type PostHistory struct {
	PostID       uint
	Post         Post
	PostURL      string `gorm:"type:text"`
	Price        int64
	Deposit      int64
	Rent         int64
	City         string `gorm:"type:varchar(63)"`
	Neighborhood string `gorm:"type:varchar(63)"`
	Area         int
	BedroomNum   int
	BuyMode      types.BuyMode  `gorm:"type:string"`
	Building     types.Building `gorm:"type:string"`
	Age          uint8
	FloorsNum    uint8
	HasStorage   bool
	HsaParking   bool
	HasElevator  bool
	ImageURL     string `gorm:"type:text"`
	Description  string `gorm:"type:text"`
	gorm.Model
}
