package db

import (
	"fmt"
	"github.com/MagicalCrawler/RealEstateApp/types"
	"strconv"

	"github.com/MagicalCrawler/RealEstateApp/models"
	"github.com/MagicalCrawler/RealEstateApp/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection() *gorm.DB {
	host := utils.GetConfig("POSTGRES_HOST")
	user := utils.GetConfig("POSTGRES_USER")
	password := utils.GetConfig("POSTGRES_PASSWORD")
	dbname := utils.GetConfig("POSTGRES_DB_NAME")
	port := utils.GetConfig("POSTGRES_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic("Error connecting to database")
	}

	db.AutoMigrate(&models.User{})
	seedSuperAdminUser(db)

	db.AutoMigrate(&models.Post{}, &models.PostHistory{}, &models.Bookmark{})
	samplePost(db)

	return db
}

func seedSuperAdminUser(db *gorm.DB) {
	superAdminTelegramId, _ := strconv.ParseUint(utils.GetConfig("SUPER_ADMIN"), 10, 64)
	superAdminUser := models.User{
		TelegramID: superAdminTelegramId,
		Role:       models.SUPER_ADMIN,
	}
	if err := db.FirstOrCreate(&superAdminUser, models.User{TelegramID: superAdminTelegramId}).Error; err != nil {
		fmt.Printf("Could not seed super-admin user (%v): %v", superAdminTelegramId, err)
		panic("Could not seed super-admin user")
	}
}

func samplePost(db *gorm.DB) {
	post := models.Post{
		Title: "MyHouse",
	}

	postHistory := models.PostHistory{
		HistoryID:    post.DetailID,
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

	// Save the post along with its details to the database

	if err := db.FirstOrCreate(&post).Error; err != nil {
		fmt.Printf("PostTableCreationError")
		panic("PostTableCreationError")
	}
	if err := db.FirstOrCreate(&postHistory).Error; err != nil {
		fmt.Printf("PostHistoryTableCreationError")
		panic("PostHistoryTableCreationError")
	}

	admin := models.User{
		TelegramID: 111,
		Role:       models.ADMIN,
	}
	db.FirstOrCreate(&admin)

	bookmark := models.Bookmark{UserID: admin.ID, PostID: post.ID}
	if err := db.Create(&bookmark).Error; err != nil {
		fmt.Printf("BookmarkTableCreationError")
		panic("BookmarkTableCreationError")
	}
}
