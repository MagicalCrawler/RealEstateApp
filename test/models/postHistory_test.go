package models

import (
	"os"
	"testing"

	"github.com/MagicalCrawler/RealEstateApp/db"
	"github.com/MagicalCrawler/RealEstateApp/models"
	"github.com/MagicalCrawler/RealEstateApp/types"
	"gorm.io/gorm"
)

var dbConnection *gorm.DB

func TestMain(m *testing.M) {
	env := map[string]string{
		"POSTGRES_HOST":     "localhost",
		"POSTGRES_USER":     "admin",
		"POSTGRES_PASSWORD": "123456",
		"POSTGRES_DB_NAME":  "MagicCrawler",
		"POSTGRES_PORT":     "5432",
		"SUPER_ADMIN":       "123456789",
	}
	for key, val := range env {
		os.Setenv(key, val)
	}
	dbConnection = db.NewConnection()
	m.Run()
}

func clearData() {
	dbConnection.Exec("DELETE FROM post_histories")
	dbConnection.Exec("DELETE FROM posts")
}

func TestNotExistingPost(t *testing.T) {
	clearData()
	defer clearData()
	postHistory := models.PostHistory{
		PostID:       1,
		PostURL:      "https://divar.ir/v/%DB%B1%DB%B1%DB%B0%D9%85%D8%AA%D8%B1-%DA%A9%D9%84%DB%8C%D8%AF%D9%86%D8%AE%D9%88%D8%B1%D8%AF%D9%87-%D8%A8%D8%B1%D8%AC-%D8%B3%D8%AA%D8%A7%D8%B1%D9%87/gZ8rtsHM",
		Price:        127000000,
		City:         "تهران",
		Neighbor:     "دانشگاه",
		Area:         80,
		Age:          12,
		FloorsNum:    8,
		BuyMode:      types.Shopping,
		Building:     types.Apartment,
		HasWareHouse: false,
		HsaParking:   false,
		HasElevator:  true,
		ImageURL:     "https://s100.divarcdn.com/static/photo/neda/post/SMnTpk3hEBI4iH397OO0oQ/b55a998a-57d8-4b6c-b9a2-149c96908c77.jpg",
	}

	err := dbConnection.Create(&postHistory).Error
	if err == nil {
		t.Fatalf(`Insert PostHistory with wrong postId: %v`, err)
	}
}

func TestInsertPostSimple(t *testing.T) {
	clearData()
	defer clearData()
	post := models.Post{
		Title: "test-title",
	}
	if err := dbConnection.Create(&post).Error; err != nil {
		t.Fatalf(`Insert Post Failed: %v`, err)
	}
	postHistory := models.PostHistory{
		Post:         post,
		PostURL:      "https://divar.ir/v/%DB%B1%DB%B1%DB%B0%D9%85%D8%AA%D8%B1-%DA%A9%D9%84%DB%8C%D8%AF%D9%86%D8%AE%D9%88%D8%B1%D8%AF%D9%87-%D8%A8%D8%B1%D8%AC-%D8%B3%D8%AA%D8%A7%D8%B1%D9%87/gZ8rtsHM",
		Price:        127000000,
		City:         "تهران",
		Neighbor:     "دانشگاه",
		Area:         80,
		Age:          12,
		FloorsNum:    8,
		BuyMode:      types.Shopping,
		Building:     types.Apartment,
		HasWareHouse: false,
		HsaParking:   false,
		HasElevator:  true,
		ImageURL:     "https://s100.divarcdn.com/static/photo/neda/post/SMnTpk3hEBI4iH397OO0oQ/b55a998a-57d8-4b6c-b9a2-149c96908c77.jpg",
	}

	err := dbConnection.Create(&postHistory).Error
	if err != nil {
		t.Fatalf(`Insert PostHistory Failed: %v`, err)
	}
}
