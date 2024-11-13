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
	FindByTelegramID(TelegramID uint64) (models.User, error)
	FindAll() []models.User
	Delete(ID uint) error
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

func (ur UserRepositoryImpl) Find(ID uint) (models.User, error) {
	var user models.User

	result := ur.dbConnection.Find(&user, ID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// TODO
	}
	return user, result.Error
}
func (ur UserRepositoryImpl) FindByTelegramID(TelegramID uint64) (models.User, error) {
	var user models.User
	result := ur.dbConnection.First(&user, "telegram_id = ?", TelegramID) // Ensure you query by TelegramID

	if result.RowsAffected == 0 {
		return models.User{}, nil // User was not found, no error
	}

	if result.Error != nil {
		fmt.Println("error occurred")
		return models.User{}, result.Error // Return the error if any
	}

	return user, nil
}

func (ur UserRepositoryImpl) FindAll() []models.User {
	var users []models.User
	result := ur.dbConnection.Find(users)
	fmt.Println(result.RowsAffected)
	return users
}

func (ur UserRepositoryImpl) Delete(id uint) error {
	user, err := ur.Find(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		} else {
			return err
		}
	} else if user.Role == models.SUPER_ADMIN {
		return errors.New("super-admin user not allowed to delete")
	}
	err = ur.dbConnection.Delete(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}
