package models

import (
	"time"

	"gorm.io/gorm"
)

func CreateWatchList(db *gorm.DB, userID uint, filterItem SearchItem, interval int) (*WatchList, error) {
	// Save the FilterItem first
	if err := db.Create(&filterItem).Error; err != nil {
		return nil, err
	}

	// Create the WatchList with the associated SearchItemID
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
