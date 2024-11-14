package models

import (
	"testing"

	"github.com/MagicalCrawler/RealEstateApp/models"
	"github.com/MagicalCrawler/RealEstateApp/types"
)

func TestNotExistingPost(t *testing.T) {
	clearData()
	defer clearData()
	crawlInfo := models.CrawlHistory{}
	postHistory := models.PostHistory{
		PostID:         1,
		Title:          "آگهی",
		PostURL:        "https://divar.ir/v/%DB%B1%DB%B1%DB%B0%D9%85%D8%AA%D8%B1-%DA%A9%D9%84%DB%8C%D8%AF%D9%86%D8%AE%D9%88%D8%B1%D8%AF%D9%87-%D8%A8%D8%B1%D8%AC-%D8%B3%D8%AA%D8%A7%D8%B1%D9%87/gZ8rtsHM",
		Price:          127000000,
		City:           "تهران",
		Neighborhood:   "دانشگاه",
		Area:           80,
		Age:            12,
		FloorsNum:      8,
		BuyMode:        types.Shopping,
		Building:       types.Apartment,
		HasStorage:     false,
		HsaParking:     false,
		HasElevator:    true,
		ImageURL:       "https://s100.divarcdn.com/static/photo/neda/post/SMnTpk3hEBI4iH397OO0oQ/b55a998a-57d8-4b6c-b9a2-149c96908c77.jpg",
		CrawlHistory:   crawlInfo,
		CrawlHistoryID: crawlInfo.ID,
	}

	err := dbConnection.Create(&postHistory).Error
	if err == nil {
		t.Fatalf(`Insert PostHistory with wrong postId: %v`, err)
	}
}

func TestCrawlHistory(t *testing.T) {
	crawlInfo := models.CrawlHistory{PostNum: 123}
	if err := dbConnection.Create(&crawlInfo).Error; err != nil {
		t.Fatalf(`Insert Crawl History error: %v`, err)
	}
}

func TestInsertPostSimple(t *testing.T) {
	clearData()
	defer clearData()
	crawlInfo := models.CrawlHistory{}
	if err := dbConnection.Create(&crawlInfo).Error; err != nil {
		t.Fatalf(`Insert Crawl History error: %v`, err)
	}
	post := models.Post{
		UniqueCode: "qweqwe",
	}
	if err := dbConnection.Create(&post).Error; err != nil {
		t.Fatalf(`Insert Post Failed: %v`, err)
	}
	postHistory := models.PostHistory{
		Post:           post,
		Title:          "test-title",
		PostURL:        "https://divar.ir/v/%DB%B1%DB%B1%DB%B0%D9%85%D8%AA%D8%B1-%DA%A9%D9%84%DB%8C%D8%AF%D9%86%D8%AE%D9%88%D8%B1%D8%AF%D9%87-%D8%A8%D8%B1%D8%AC-%D8%B3%D8%AA%D8%A7%D8%B1%D9%87/gZ8rtsHM",
		Price:          127000000,
		City:           "تهران",
		Neighborhood:   "دانشگاه",
		Area:           80,
		Age:            12,
		FloorsNum:      8,
		BuyMode:        types.Shopping,
		Building:       types.Apartment,
		HasStorage:     false,
		HsaParking:     false,
		HasElevator:    true,
		ImageURL:       "https://s100.divarcdn.com/static/photo/neda/post/SMnTpk3hEBI4iH397OO0oQ/b55a998a-57d8-4b6c-b9a2-149c96908c77.jpg",
		CrawlHistory:   crawlInfo,
		CrawlHistoryID: crawlInfo.ID,
	}

	err := dbConnection.Create(&postHistory).Error
	if err != nil {
		t.Fatalf(`Insert PostHistory Failed: %v`, err)
	}
}

func TestInsertExistPost(t *testing.T) {
	clearData()
	defer clearData()
	post := models.Post{
		UniqueCode: "FrfT",
	}
	if err := dbConnection.Create(&post).Error; err != nil {
		t.Fatalf(`Insert Post Failed: %v`, err)
	}

	post1 := models.Post{
		UniqueCode: "FrfT",
	}
	var existPost bool
	err := dbConnection.Table("posts").Select("count(*) > 0").Where("unique_code = ?", post1.UniqueCode).Find(&existPost).Error

	if err == nil {
		if existPost == false {
			if err := dbConnection.Create(&post1).Error; err != nil {
				t.Fatalf(`Insert Post Failed: %v`, err)
			}
		} else {
			println("Post already exist")
		}
	}

}
