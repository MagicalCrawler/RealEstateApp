package models

import (
	"os"
	"testing"

	"github.com/MagicalCrawler/RealEstateApp/db"
	"github.com/MagicalCrawler/RealEstateApp/models"
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
	superAdmin := models.User{}
	err := dbConnection.Where("Role = ?", models.SUPER_ADMIN).First(&superAdmin).Error
	if err != nil {
		panic("super admin does not exists")
	}
	m.Run()
}

func clearData() {
	dbConnection.Unscoped().Where("1 = 1").Delete(&models.Bookmark{})
	dbConnection.Unscoped().Where("1 = 1").Delete(&models.PostHistory{})
	dbConnection.Unscoped().Where("1 = 1").Delete(&models.Post{})
	dbConnection.Unscoped().Where("Role <> ?", models.SUPER_ADMIN).Delete(&models.User{})

}
