package db

import (
	"fmt"
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
