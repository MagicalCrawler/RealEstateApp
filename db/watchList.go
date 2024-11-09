package models

import (
	"time"

	"gorm.io/gorm"
)

func CreateWatchList(db *gorm.DB, userID uint, filterItem FilterItem, interval int) (*WatchList, error) {
	// Save the FilterItem first
	if err := db.Create(&filterItem).Error; err != nil {
		return nil, err
	}

	// Create the WatchList with the associated FilterItemID
	watchList := WatchList{
		UserID:          userID,
		FilterItemID:    filterItem.ID,
		RefreshInterval: interval,
		LastChecked:     time.Now(),
	}

	if err := db.Create(&watchList).Error; err != nil {
		return nil, err
	}
	return &watchList, nil
}

func GetAllWatchListsForUser(db *gorm.DB, userID uint) ([]WatchList, error) {
	var watchLists []WatchList
	if err := db.Preload("FilterItem").Where("user_id = ?", userID).Find(&watchLists).Error; err != nil {
		return nil, err
	}
	return watchLists, nil
}

func DeleteWatchList(db *gorm.DB, watchListID uint) error {
	var watchList WatchList
	if err := db.First(&watchList, watchListID).Error; err != nil {
		return err
	}

	// Delete the associated FilterItem
	if err := db.Delete(&FilterItem{}, watchList.FilterItemID).Error; err != nil {
		return err
	}

	// Delete the WatchList entry
	if err := db.Delete(&WatchList{}, watchListID).Error; err != nil {
		return err
	}

	return nil
}

func UpdateWatchList(db *gorm.DB, watchListID uint, updatedData WatchList, updatedFilterItem FilterItem) (*WatchList, error) {
	var watchList WatchList
	if err := db.First(&watchList, watchListID).Error; err != nil {
		return nil, err
	}

	// Update FilterItem fields if needed
	updatedFilterItem.ID = watchList.FilterItemID // Ensure we're updating the correct FilterItem
	if err := db.Model(&FilterItem{}).Where("id = ?", watchList.FilterItemID).Updates(updatedFilterItem).Error; err != nil {
		return nil, err
	}

	// Update WatchList fields
	if err := db.Model(&watchList).Updates(updatedData).Error; err != nil {
		return nil, err
	}

	return &watchList, nil
}
