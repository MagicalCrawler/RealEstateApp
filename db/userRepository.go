package db

import (
	"errors"
	"fmt"

	"github.com/MagicalCrawler/RealEstateApp/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Save(models.User) (models.User, error)
	Find(ID uint) (models.User, error)
	FindAll() []models.User
}

type UserRepositoryImpl struct {
	dbConnection *gorm.DB
}

func CreateNewUserRepository(dbConnection *gorm.DB) UserRepository {
	return UserRepositoryImpl{dbConnection: dbConnection}
}

func (ur UserRepositoryImpl) Save(user models.User) (models.User, error) {
	err := ur.dbConnection.Create(&user).Error
	return user, err
}

func (ur UserRepositoryImpl) Find(id uint) (models.User, error) {
	var user models.User
	result := ur.dbConnection.Find(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// TODO
	}
	return user, result.Error
}

func (ur UserRepositoryImpl) FindAll() []models.User {
	var users []models.User
	result := ur.dbConnection.Find(users)
	fmt.Println(result.RowsAffected)
	return users
}