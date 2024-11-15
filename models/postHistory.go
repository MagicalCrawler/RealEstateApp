package models

import (
	"github.com/MagicalCrawler/RealEstateApp/types"
)

type PostHistory struct {
	ID             uint `gorm:"primary_key;auto_increment"`
	PostID         uint
	Post           Post
	Title          string `gorm:"type:text"`
	PostURL        string `gorm:"type:text"`
	Price          int64
	Deposit        int64
	Rent           int64
	City           string `gorm:"type:varchar(63)"`
	Neighborhood   string `gorm:"type:varchar(63)"`
	Area           int
	BedroomNum     int
	BuyMode        types.BuyMode  `gorm:"type:string"`
	Building       types.Building `gorm:"type:string"`
	Age            uint8
	FloorsNum      uint8
	HasStorage     bool
	HasParking     bool
	HasElevator    bool
	ImageURL       string `gorm:"type:text"`
	Description    string `gorm:"type:text"`
	CrawlHistory   CrawlHistory
	CrawlHistoryID uint
	Capacity       string
	NormalDays     string
	Weekend        string
	Holidays       string
	CostPerPerson  string
}
