package db

import (
	"time"

	"github.com/MagicalCrawler/RealEstateApp/models"
	"gorm.io/gorm"
)

type WatchListRepository interface {
	Create(userID uint, filterItem models.FilterItem, interval int) (models.WatchList, error)
	FindByID(id uint) (models.WatchList, error)
	FindAll() ([]models.WatchList, error)
	Update(id uint, updatedData models.WatchList, updatedFilterItem models.FilterItem) (models.WatchList, error)
	Delete(id uint) error
	GetAllWatchListsForUser(userID uint) ([]models.WatchList, error)
}

type WatchListRepositoryImpl struct {
	dbConnection *gorm.DB
}

func NewWatchListRepository(dbConnection *gorm.DB) WatchListRepository {
	return &WatchListRepositoryImpl{dbConnection: dbConnection}
}

func (repo WatchListRepositoryImpl) Create(userID uint, filterItem models.FilterItem, interval int) (models.WatchList, error) {
	// Save the FilterItem first
	if err := repo.dbConnection.Create(&filterItem).Error; err != nil {
		return models.WatchList{}, err
	}

	// Create the WatchList with the associated FilterItemID
	watchList := models.WatchList{
		UserID:          userID,
		FilterItemID:    filterItem.ID,
		RefreshInterval: interval,
		LastChecked:     time.Now(),
	}

	if err := repo.dbConnection.Create(&watchList).Error; err != nil {
		return models.WatchList{}, err
	}
	return watchList, nil
}

func (repo WatchListRepositoryImpl) FindByID(id uint) (models.WatchList, error) {
	var watchList models.WatchList
	err := repo.dbConnection.Preload("FilterItem").First(&watchList, id).Error
	return watchList, err
}

func (repo WatchListRepositoryImpl) FindAll() ([]models.WatchList, error) {
	var watchLists []models.WatchList
	err := repo.dbConnection.Preload("FilterItem").Find(&watchLists).Error
	return watchLists, err
}

func (repo WatchListRepositoryImpl) Update(id uint, updatedData models.WatchList, updatedFilterItem models.FilterItem) (models.WatchList, error) {
	var watchList models.WatchList
	if err := repo.dbConnection.First(&watchList, id).Error; err != nil {
		return models.WatchList{}, err
	}

	// Update FilterItem fields if needed
	updatedFilterItem.ID = watchList.FilterItemID // Ensure we're updating the correct FilterItem
	if err := repo.dbConnection.Model(&models.FilterItem{}).Where("id = ?", watchList.FilterItemID).Updates(updatedFilterItem).Error; err != nil {
		return models.WatchList{}, err
	}

	// Update WatchList fields
	if err := repo.dbConnection.Model(&watchList).Updates(updatedData).Error; err != nil {
		return models.WatchList{}, err
	}

	return watchList, nil
}

func (repo WatchListRepositoryImpl) Delete(id uint) error {
	var watchList models.WatchList
	if err := repo.dbConnection.First(&watchList, id).Error; err != nil {
		return err
	}

	// Delete the associated FilterItem
	if err := repo.dbConnection.Delete(&models.FilterItem{}, watchList.FilterItemID).Error; err != nil {
		return err
	}

	// Delete the WatchList entry
	if err := repo.dbConnection.Delete(&models.WatchList{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo WatchListRepositoryImpl) GetAllWatchListsForUser(userID uint) ([]models.WatchList, error) {
	var watchLists []models.WatchList
	err := repo.dbConnection.Preload("FilterItem").Where("user_id = ?", userID).Find(&watchLists).Error
	return watchLists, err
}
