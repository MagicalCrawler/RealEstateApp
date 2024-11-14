package db

import (
	"github.com/MagicalCrawler/RealEstateApp/models"
	"gorm.io/gorm"
)

type WatchListRepository interface {
	Create(watchList models.WatchList) (models.WatchList, error)
	FindByID(id uint) (models.WatchList, error)
	FindAll() ([]models.WatchList, error)
	Update(id uint, updatedData models.WatchList) (models.WatchList, error)
	Delete(id uint) error
}

type WatchListRepositoryImpl struct {
	dbConnection *gorm.DB
}

func NewWatchListRepository(dbConnection *gorm.DB) WatchListRepository {
	return &WatchListRepositoryImpl{dbConnection: dbConnection}
}

// Create adds a new WatchList to the database
func (repo *WatchListRepositoryImpl) Create(watchList models.WatchList) (models.WatchList, error) {
	err := repo.dbConnection.Create(&watchList).Error
	return watchList, err
}

// FindByID retrieves a WatchList by ID
func (repo *WatchListRepositoryImpl) FindByID(id uint) (models.WatchList, error) {
	var watchList models.WatchList
	err := repo.dbConnection.First(&watchList, id).Error
	return watchList, err
}

// FindAll retrieves all WatchLists
func (repo *WatchListRepositoryImpl) FindAll() ([]models.WatchList, error) {
	var watchLists []models.WatchList
	err := repo.dbConnection.Find(&watchLists).Error
	return watchLists, err
}

// Update updates a WatchList's data
func (repo *WatchListRepositoryImpl) Update(id uint, updatedData models.WatchList) (models.WatchList, error) {
	var watchList models.WatchList
	if err := repo.dbConnection.First(&watchList, id).Error; err != nil {
		return watchList, err
	}

	// Update fields based on updatedData
	err := repo.dbConnection.Model(&watchList).Updates(updatedData).Error
	return watchList, err
}

// Delete removes a WatchList by ID
func (repo *WatchListRepositoryImpl) Delete(id uint) error {
	if err := repo.dbConnection.Delete(&models.WatchList{}, id).Error; err != nil {
		return err
	}
	return nil
}
