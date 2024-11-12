package models

import (
	// "os"
	"testing"

	"github.com/MagicalCrawler/RealEstateApp/db"
	"github.com/MagicalCrawler/RealEstateApp/models"
	// 	"gorm.io/gorm"
)

// var dbConnection *gorm.DB

// func TestMain(m *testing.M) {
// 	env := map[string]string{
// 		"POSTGRES_HOST":     "localhost",
// 		"POSTGRES_USER":     "admin",
// 		"POSTGRES_PASSWORD": "123456",
// 		"POSTGRES_DB_NAME":  "MagicCrawler",
// 		"POSTGRES_PORT":     "5432",
// 		"SUPER_ADMIN":       "123456789",
// 	}
// 	for key, val := range env {
// 		os.Setenv(key, val)
// 	}
// 	dbConnection = db.NewConnection()
// 	m.Run()
// }

// func clearData() {
// 	dbConnection.Exec("DELETE FROM users")
// }

func TestInsertUserDirectly(t *testing.T) {
	clearData()
	defer clearData()
	user := models.User{
		TelegramID: 1,
		Role:       models.ADMIN,
	}

	err := dbConnection.Create(&user).Error
	if err != nil {
		t.Fatalf(`Insert User Failed: %v`, err)
	}
}

func TestInsertUserWithRepository(t *testing.T) {
	clearData()
	defer clearData()

	userRepository := db.CreateNewUserRepository(dbConnection)

	user := models.User{
		TelegramID: 1,
		Role:       models.ADMIN,
	}
	user, err := userRepository.Save(user)
	if err != nil {
		t.Fatalf(`Insert User Failed: %v`, err)
	}

	_, err = userRepository.Find(user.ID)
	if err != nil {
		t.Fatalf(`Find User Failed: %v`, err)
	}
}
