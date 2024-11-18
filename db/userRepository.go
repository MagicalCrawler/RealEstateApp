package db

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/MagicalCrawler/RealEstateApp/models"
	"github.com/MagicalCrawler/RealEstateApp/utils"
	"gorm.io/gorm"
)

type UserRepository interface {
	Save(models.User) (models.User, error)
	Find(ID uint) (models.User, error)
	FindByTelegramID(TelegramID uint64) (models.User, error)
	FindAll() ([]models.User, error)
	FindAllUsersByRole(models.Role) ([]models.User, error)
	Delete(ID uint) error
	UpdateUserType(ID uint, Type models.UserType) (models.User, error)
	UpdateUserRole(ID uint, Role models.Role) (models.User, error)
	UpdateUser(ID uint, updatedData map[string]interface{}) (models.User, error)
}

type UserRepositoryImpl struct {
	dbConnection *gorm.DB
	logger       *slog.Logger
}

func CreateNewUserRepository(dbConnection *gorm.DB) UserRepository {
	return UserRepositoryImpl{dbConnection: dbConnection, logger: utils.NewLogger("database")}
}

func (ur UserRepositoryImpl) Save(user models.User) (models.User, error) {
	err := ur.dbConnection.Create(&user).Error
	if err != nil {
		ur.logger.Error("Insert User Failed: %v", slog.Attr{Key: "err", Value: slog.AnyValue(err)})
	}
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

func (ur UserRepositoryImpl) FindAll() ([]models.User, error) {
	var users []models.User
	result := ur.dbConnection.Find(&users)
	if result.Error != nil {
		fmt.Println("error occurred")
		return users, result.Error // Return the error if any
	}
	// fmt.Println(result.RowsAffected)
	return users, nil
}
func (ur UserRepositoryImpl) FindAllUsersByRole(role models.Role) ([]models.User, error) {
	var users []models.User
	result := ur.dbConnection.Where("role = ?", role).Find(&users)
	if result.Error != nil {
		fmt.Println("error occurred")
		return users, result.Error // Return the error if any
	}
	// fmt.Println(result.RowsAffected)
	return users, nil
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
func (ur UserRepositoryImpl) UpdateUserType(ID uint, Type models.UserType) (models.User, error) {
	var user models.User
	result := ur.dbConnection.Model(&user).Where("id = ?", ID).Where("role = ?", models.USER).Update("type", Type)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (ur UserRepositoryImpl) UpdateUserRole(ID uint, Role models.Role) (models.User, error) {
	var user models.User
	result := ur.dbConnection.Model(&user).Where("id = ?", ID).Where("role = ?", models.USER).Update("role", Role)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (ur UserRepositoryImpl) UpdateUser(ID uint, updatedData map[string]interface{}) (models.User, error) {
	var user models.User

	// Find the user to update
	if err := ur.dbConnection.First(&user, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user not found")
		}
		ur.logger.Error("Find User Failed", slog.Any("error", err))
		return models.User{}, err
	}

	// Update the fields specified in updatedData
	if err := ur.dbConnection.Model(&user).Updates(updatedData).Error; err != nil {
		ur.logger.Error("Update User Failed", slog.Any("error", err))
		return models.User{}, err
	}

	return user, nil
}
