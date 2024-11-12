package models

import (
	"errors"
	"testing"

	"github.com/MagicalCrawler/RealEstateApp/db"
	"github.com/MagicalCrawler/RealEstateApp/models"
	"gorm.io/gorm"
)

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
	_, err := userRepository.Find(user.ID)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf(`Find User Failed: %v`, err)
	}

	user, err = userRepository.Save(user)
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
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf(`Find User Failed: %v`, err)
	}
}
