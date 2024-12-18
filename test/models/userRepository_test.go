package models

import (

	// "os"

	"errors"

	"gorm.io/gorm"

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

func TestInsertAndDeleteUserWithRepository(t *testing.T) {
	clearData()
	defer clearData()

	userRepository := db.CreateNewUserRepository(dbConnection)

	user := models.User{
		TelegramID: 1,
		Role:       models.ADMIN,
	}

	//_, err := userRepository.Find(user.ID)
	//if !errors.Is(err, gorm.ErrRecordNotFound) {
	//	t.Fatalf(`Find User Failed: %v`, err)
	//}

	user, err := userRepository.Save(user)
	if err != nil {
		t.Fatalf(`Insert User Failed: %v`, err)
	}

	user, err = userRepository.Find(user.ID)
	if err != nil {
		t.Fatalf(`Find User Failed: %v`, err)
	}

	err = dbConnection.Delete(&user, user.ID).Error
	if err != nil {
		t.Fatalf(`Delete User Failed: %v`, err)
	}
	user, err = userRepository.Find(user.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf(`Find User Failed: %v`, err)
	}
}
func TestInsertAndFindByTelegramIDUserWithRepository(t *testing.T) {
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

	_, err = userRepository.FindByTelegramID(user.TelegramID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf(`Find User Failed: %v`, err)
	}
	//u, err := userRepository.FindByTelegramID(1)
	//empty_user := models.User{}
	//if !errors.Is(err, nil) && u != empty_user {
	//	t.Fatalf(`Find User Failed: %v`, err)
	//}

	user, err = userRepository.Find(user.ID)
	if err != nil {
		t.Fatalf(`Find User Failed: %v`, err)
	}

	err = dbConnection.Delete(&user, user.ID).Error
	if err != nil {
		t.Fatalf(`Delete User Failed: %v`, err)
	}
	user, err = userRepository.Find(user.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf(`Find User Failed: %v`, err)
	}
}
